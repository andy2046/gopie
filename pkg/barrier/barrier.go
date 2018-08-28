// Package barrier provides a barrier implementation.
package barrier

import (
	"context"
	"errors"
	"sync"
)

// Barrier is a synchronizer that allows members to wait for each other.
type Barrier struct {
	count      int
	n          int
	isBroken   bool
	waitChan   chan struct{}
	fallChan   chan struct{}
	brokenChan chan struct{}
	mu         sync.RWMutex
}

var (
	// ErrBroken when barrier is broken.
	ErrBroken = errors.New("barrier is broken")
)

// New returns a new barrier.
func New(n int) *Barrier {
	if n <= 0 {
		panic("number of members must be positive int")
	}
	b := &Barrier{
		n:          n,
		waitChan:   make(chan struct{}),
		brokenChan: make(chan struct{}),
	}
	return b
}

// Await waits until all members have called await on the barrier.
func (b *Barrier) Await(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	b.mu.Lock()

	if b.isBroken {
		b.mu.Unlock()
		return ErrBroken
	}

	b.count++
	waitChan := b.waitChan
	brokenChan := b.brokenChan
	count := b.count

	if count == b.n {
		b.reset(true)
		b.mu.Unlock()
		return nil
	}

	b.mu.Unlock()

	select {
	case <-waitChan:
		return nil
	case <-brokenChan:
		return ErrBroken
	case <-ctx.Done():
		b.broke(true)
		return ctx.Err()
	}
}

func (b *Barrier) broke(toLock bool) {
	if toLock {
		b.mu.Lock()
		defer b.mu.Unlock()
	}

	if !b.isBroken {
		b.isBroken = true
		close(b.brokenChan)
	}
}

func (b *Barrier) reset(ok bool) {
	if ok {
		close(b.waitChan)
	} else if b.count > 0 {
		b.broke(false)
	}

	b.waitChan = make(chan struct{})
	b.brokenChan = make(chan struct{})
	b.count = 0
	b.isBroken = false
}

// Reset resets the barrier to initial state.
func (b *Barrier) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.reset(false)
}

// N returns the number of members for the barrier.
func (b *Barrier) N() int {
	return b.n
}

// NWaiting returns the number of members currently waiting at the barrier.
func (b *Barrier) NWaiting() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.count
}

// IsBroken returns true if the barrier is broken.
func (b *Barrier) IsBroken() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.isBroken
}
