package pubsub

import (
	"context"
	"fmt"
	"testing"
)

func TestNewTopic(t *testing.T) {
	psName := "pubsub-01"
	psNameExpected := "projects/" + psName
	ps := New(psName)
	if ps.Name() != psNameExpected {
		t.Error("wrong PubSub name, expected", psNameExpected, ", got", ps.Name())
	}

	topic1, err := ps.NewTopic("topic-01", 10, 2)
	if err != nil {
		t.Error("NewTopic error", err)
	}

	sub11, _ := topic1.NewSubscription(2)
	sub11.Receive(func(m *Message) {
		t.Logf("sub11 Message -> %#v", *m)
	})
	sub12, _ := topic1.NewSubscription(2)
	sub12.Receive(func(m *Message) {
		t.Logf("sub12 Message -> %#v", *m)
	})

	err1 := topic1.Publish(context.Background(), &Message{
		ID:   "001",
		Data: []byte("data"),
	})
	if err1 != nil {
		t.Logf("topic1.Publish error -> %#v", err1)
	}

	topic1.Delete()

	topic2, err := ps.NewTopic("topic-02", 10, 2)
	if err != nil {
		t.Error("NewTopic error", err)
	}

	sub2, _ := topic2.NewSubscription(2)
	sub2.Receive(func(m *Message) {
		t.Logf("sub2 Message -> %#v", *m)
	})

	err2 := topic2.Publish(context.Background(), &Message{
		ID:   "002",
		Data: []byte("data"),
	})
	if err2 != nil {
		t.Logf("topic2.Publish error -> %#v", err2)
	}

	t.Logf("Topics -> %#v", ps.Topics())
	if fmt.Sprintf("%v", ps.Topics()) != fmt.Sprintf("%v", []string{"topic-02"}) {
		t.Error("Topics should return topic name")
	}

	topic2.Delete()
}

func TestNewSubscription(t *testing.T) {
	psName := "pubsub-02"
	psNameExpected := "projects/" + psName
	ps := New(psName)
	if ps.Name() != psNameExpected {
		t.Error("wrong PubSub name, expected", psNameExpected, ", got", ps.Name())
	}

	topic1, err := ps.NewTopic("topic-01", 10, 2)
	if err != nil {
		t.Error("NewTopic error", err)
	}

	sub11, _ := topic1.NewSubscription(2)
	sub11.Receive(func(m *Message) {
		t.Logf("sub11 Message -> %#v", *m)
	})
	sub12, _ := topic1.NewSubscription(2)
	sub12.Receive(func(m *Message) {
		t.Logf("sub12 Message -> %#v", *m)
	})
	sub13, _ := topic1.NewSubscription(2)
	sub13.Receive(func(m *Message) {
		t.Logf("sub13 Message -> %#v", *m)
	})

	err1 := topic1.Publish(context.Background(), &Message{
		ID:   "001",
		Data: []byte("data"),
	})
	if err1 != nil {
		t.Logf("topic1.Publish error -> %#v", err1)
	}

	sub13.Delete()
	topic1.Delete()
}
