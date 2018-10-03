package ratelimit

import (
	"testing"
	"time"
)

var (
	t0 = time.Date(2017, time.August, 8, 8, 8, 8, 0, time.UTC)
)

func TestEvery(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "ratelimit: invalid time interval for Every" {
				t.Error("Every panic messgae error")
			}
		}
	}()

	cases := []struct {
		interval time.Duration
		l        Limit
	}{
		{1 * time.Second, Limit(1)},
		{2 * time.Second, Limit(0.5)},
		{4 * time.Second, Limit(0.25)},
		{10 * time.Second, Limit(0.1)},
		{1 * time.Millisecond, Limit(1e3)},
		{10 * time.Millisecond, Limit(100)},
		{-1 * time.Millisecond, Limit(-1)},
	}
	for _, c := range cases {
		l := Every(c.interval)
		if l-c.l > 0.0001 {
			t.Errorf("Every(%v) = %v want %v", c.interval, l, c.l)
		}
	}
}

func TestAllow(t *testing.T) {
	t1 := t0.Add(100 * time.Millisecond)
	t2 := t0.Add(200 * time.Millisecond)
	cases := []struct {
		t  time.Time
		n  int
		ok bool
	}{
		{t0, 1, true},
		{t0, 1, false},
		{t0, 1, false},
		{t1, 1, true},
		{t1, 1, false},
		{t1, 1, false},
		{t2, 2, false}, // exceeds burst
		{t2, 1, true},
		{t2, 1, false},
	}

	l := New(10, 1)
	for _, c := range cases {
		ok := l.AllowN(c.t, c.n)
		if ok != c.ok {
			t.Errorf("AllowN(%v, %v) = %v want %v",
				c.t, c.n, ok, c.ok)
		}
	}
}

func TestWait(t *testing.T) {
	cases := []struct {
		t      time.Time
		n      int
		nilErr bool
	}{
		{t0, 1, true},
		{t0, 2, false}, // exceeds burst
		{t0, 1, true},
		{t0, 1, true},
	}

	l := New(10, 1)
	for _, c := range cases {
		w, err := l.WaitN(c.t, c.n)
		delay := l.limit.durationFromTokens(float64(c.n))
		if (c.nilErr && err != nil) || (!c.nilErr && err == nil) || w > delay {
			t.Errorf("WaitN(%v, %v) = %v, %v want %v, nilErr? %v",
				c.t, c.n, w, err, delay, c.nilErr)
		}
	}
}
