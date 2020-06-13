package backoff

import (
	"testing"
	"time"
)

var result time.Duration

func TestExponential(t *testing.T) {
	baseDelay, capDelay := 1*time.Second, 2*time.Second
	bk := New(baseDelay, capDelay, Exponential)
	v := bk.Backoff(0)
	if v != baseDelay {
		t.Errorf("want %v got %v", baseDelay, v)
	}

	for i := 1; i < 10; i++ {
		v = bk.Backoff(i)
		if v != capDelay {
			t.Errorf("want %v got %v", capDelay, v)
		}
	}

	v = bk.Backoff(0)
	if v != baseDelay {
		t.Errorf("want %v got %v", baseDelay, v)
	}
}

func TestFullJitter(t *testing.T) {
	baseDelay, capDelay := 1*time.Second, 2*time.Second
	bk := New(baseDelay, capDelay, FullJitter)
	result = bk.Backoff(0)

	for i := 1; i < 10; i++ {
		result = bk.Backoff(i)
	}
}

func TestEqualJitter(t *testing.T) {
	baseDelay, capDelay := 1*time.Second, 2*time.Second
	bk := New(baseDelay, capDelay, EqualJitter)
	result = bk.Backoff(0)

	for i := 1; i < 10; i++ {
		result = bk.Backoff(i)
	}
}

func TestDecorrelatedJitter(t *testing.T) {
	baseDelay, capDelay := 1*time.Second, 2*time.Second
	bk := New(baseDelay, capDelay, DecorrelatedJitter)
	result = bk.Backoff(0)

	for i := 1; i < 10; i++ {
		result = bk.Backoff(i)
	}
}

func TestMinDuration(t *testing.T) {
	testcases := []struct {
		x    time.Duration
		y    time.Duration
		want time.Duration
	}{
		{1 * time.Second, 2 * time.Second, 1 * time.Second},
		{1 * time.Second, 1 * time.Second, 1 * time.Second},
		{1 * time.Millisecond, 2 * time.Millisecond, 1 * time.Millisecond},
		{1 * time.Millisecond, 1 * time.Millisecond, 1 * time.Millisecond},
		{1 * time.Minute, 1 * time.Hour, 1 * time.Minute},
	}

	for _, c := range testcases {
		if r := minDuration(c.x, c.y); r != c.want {
			t.Errorf("want %v got %v", c.want, r)
		}
	}
}
