// Package pubsub provides a pubsub implementation.
package pubsub

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"sync"
)

type (
	// PubSub is a Pub/Sub instance for a single project.
	PubSub struct {
		projectID string
		topics    map[string]*Topic
		mu        sync.RWMutex
	}

	// Topic represents a PubSub topic.
	Topic struct {
		// The identifier for the topic,
		// in the format "projects/<projid>/topics/<name>".
		name string

		inbox chan Message

		stopped bool

		pubSub *PubSub

		subscriptions map[string]*Subscription

		// Errors is the error output channel back to the user. You MUST read from this
		// channel or the Publish will deadlock when the channel is full.
		Errors chan PublishError

		numGoroutines int

		wgPublish sync.WaitGroup

		once sync.Once

		mu sync.RWMutex

		wg sync.WaitGroup
	}

	// Subscription represents a PubSub subscription.
	Subscription struct {
		// The identifier for the subscription,
		// in the format "projects/<projid>/topics/<name>/subscriptions/<name>".
		name string

		inbox chan Message

		topic *Topic

		stopped bool

		mu sync.RWMutex

		wg sync.WaitGroup

		done chan struct{}

		// numGoroutines is the number of goroutines it will spawn to pull msg concurrently.
		numGoroutines int
	}

	// Message represents a Pub/Sub message.
	Message struct {
		ID   string
		Data []byte
	}

	// PublishError is the error generated when it fails to publish a message.
	PublishError struct {
		Msg *Message
		Err error
	}
)

func (pe PublishError) Error() string {
	return fmt.Sprintf("failed to publish message %s -> %s", pe.Msg.ID, pe.Err)
}

// UUID generates uuid.
func UUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// New creates a new PubSub.
func New(project string) *PubSub {
	return &PubSub{
		projectID: project,
		topics:    make(map[string]*Topic),
	}
}

// Name returns the full name for the PubSub.
func (p *PubSub) Name() string {
	return fmt.Sprintf("projects/%s", p.projectID)
}

// NewTopic creates a new Topic with the given name,
// size is the channel buffer size for topic message chan,
// numGoroutines is the number of goroutines it will spawn to push msg concurrently.
func (p *PubSub) NewTopic(name string, size int, numGoroutines int) (*Topic, error) {
	p.mu.RLock()
	if _, ok := p.topics[name]; ok {
		p.mu.RUnlock()
		return nil, errors.New("duplicated topic name")
	}
	p.mu.RUnlock()
	p.mu.Lock()
	t := &Topic{
		name:          name,
		pubSub:        p,
		subscriptions: make(map[string]*Subscription),
		inbox:         make(chan Message, size),
		Errors:        make(chan PublishError, size),
		numGoroutines: numGoroutines,
	}
	p.topics[name] = t
	p.mu.Unlock()
	return t, nil
}

// Topic returns the topic by name.
func (p *PubSub) Topic(name string) *Topic {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if t, ok := p.topics[name]; ok {
		return t
	}
	return nil
}

// Topics list all the topics in the PubSub.
func (p *PubSub) Topics() []string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	ts := make([]string, 0, len(p.topics))
	for k := range p.topics {
		ts = append(ts, k)
	}
	return ts
}

// Publish publishes msg to the topic asynchronously.
func (t *Topic) Publish(ctx context.Context, msg *Message) error {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.stopped {
		return errors.New("topic stopped")
	}
	t.wgPublish.Add(1)
	go func() {
		select {
		case <-ctx.Done():
			t.Errors <- PublishError{
				msg,
				ctx.Err(),
			}
		case t.inbox <- *msg:
		}
		t.wgPublish.Done()
	}()
	return nil
}

// Delete removes itself from PubSuband stop it.
func (t *Topic) Delete() {
	t.pubSub.mu.Lock()
	for k := range t.pubSub.topics {
		if k == t.name {
			delete(t.pubSub.topics, k)
		}
	}
	t.pubSub.mu.Unlock()
	t.Stop()
}

// Name returns the full name for the topic.
func (t *Topic) Name() string {
	return fmt.Sprintf("projects/%s/topics/%s", t.pubSub.projectID, t.name)
}

// Stop stops the topic.
func (t *Topic) Stop() {
	t.mu.Lock()
	t.stopped = true
	t.wgPublish.Wait()
	close(t.inbox)
	t.mu.Unlock()
	t.wg.Wait()

	for _, v := range t.subscriptions {
		go func(s *Subscription) {
			s.mu.Lock()
			if s.stopped {
				s.mu.Unlock()
				return
			}
			s.stopped = true
			close(s.inbox)
			s.wg.Wait()
			close(s.done)
			s.mu.Unlock()
		}(v)
	}
}

func (t *Topic) start() {
	for {
		m, open := <-t.inbox
		if !open {
			log.Printf("topic %s inbox closed, exit", t.Name())
			t.wg.Done()
			return
		}
		t.mu.RLock()
		subs := make(map[string]*Subscription, len(t.subscriptions))
		for k, v := range t.subscriptions {
			subs[k] = v
		}
		t.mu.RUnlock()

		for _, v := range subs {
			go func(s *Subscription) {
				s.mu.RLock()
				if s.stopped {
					s.mu.RUnlock()
					return
				}
				s.inbox <- m
				s.mu.RUnlock()
			}(v)
		}
	}
}

// Subscriptions list all the subscriptions to this topic.
func (t *Topic) Subscriptions() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	sub := make([]string, 0, len(t.subscriptions))
	for k := range t.subscriptions {
		sub = append(sub, k)
	}
	return sub
}

// NewSubscription creates a new Subscription to this topic,
// numGoroutines is the number of goroutines it will spawn to pull msg concurrently.
func (t *Topic) NewSubscription(numGoroutines int) (*Subscription, error) {
	t.once.Do(func() {
		for range make([]struct{}, t.numGoroutines) {
			t.wg.Add(1)
			go t.start()
		}
	})
	t.mu.Lock()
	n := fmt.Sprintf("%s-sub-%s", t.name, UUID())
	s := &Subscription{
		name:          n,
		inbox:         make(chan Message, 10*numGoroutines),
		topic:         t,
		done:          make(chan struct{}),
		numGoroutines: numGoroutines,
	}
	t.subscriptions[n] = s
	t.mu.Unlock()
	return s, nil
}

// Subscription returns the subscription by name..
func (t *Topic) Subscription(name string) *Subscription {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if s, ok := t.subscriptions[name]; ok {
		return s
	}
	return nil
}

// Receive receives message for this subscription.
func (s *Subscription) Receive(f func(*Message)) {
	for range make([]struct{}, s.numGoroutines) {
		s.wg.Add(1)
		d, i := s.done, s.inbox
		go func() {
			for {
				select {
				case <-d:
					d = nil
					return
				case m, open := <-i:
					if !open {
						i = nil
						s.wg.Done()
						break
					}
					f(&m)
				}
			}
		}()
	}
}

// Delete unsubscribes itself from topic.
func (s *Subscription) Delete() {
	s.topic.mu.Lock()
	for k := range s.topic.subscriptions {
		if k == s.name {
			delete(s.topic.subscriptions, k)
		}
	}
	s.topic.mu.Unlock()
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
	}
	s.stopped = true
	close(s.inbox)
	s.wg.Wait()
	close(s.done)
	s.mu.Unlock()
}
