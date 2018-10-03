// Package ratelimit implements a rate limiter.
package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// Limit defines the maximum number of requests per second.
type Limit float64

// Every converts a time interval between requests to a Limit.
func Every(interval time.Duration) Limit {
	if interval <= 0 {
		panic("ratelimit: invalid time interval for Every")
	}
	return 1 / Limit(interval.Seconds())
}

// Limiter implements a token bucket limiter at rate `r` tokens per second with burst size of `b` tokens.
type Limiter struct {
	limit  Limit
	burst  int
	mu     sync.Mutex
	tokens float64
	// last is the last time the limiter's tokens got updated
	last time.Time
	// lastRequest is the latest time of a request
	lastRequest time.Time
}

// Limit returns the Limiter's rate.
func (l *Limiter) Limit() Limit {
	return l.limit
}

// Burst returns the Limiter's burst size.
func (l *Limiter) Burst() int {
	return l.burst
}

// New returns a new Limiter at rate `r` tokens per second with burst of `b` tokens.
func New(r Limit, b int) *Limiter {
	return &Limiter{
		limit: r,
		burst: b,
	}
}

// Allow is the shortcut for AllowN(time.Now(), 1).
func (l *Limiter) Allow() bool {
	return l.AllowN(time.Now(), 1)
}

// AllowN checks whether `n` requests may happen at time `now`.
func (l *Limiter) AllowN(now time.Time, n int) bool {
	return l.reserveN(now, n, 0).ok
}

// Wait is the shortcut for WaitN(time.Now(), 1).
func (l *Limiter) Wait() (time.Duration, error) {
	return l.WaitN(time.Now(), 1)
}

// WaitN calculates the time duration to wait before `n` requests may happen at time `now`.
func (l *Limiter) WaitN(now time.Time, n int) (time.Duration, error) {
	return l.waitN(now, n)
}

func (l *Limiter) waitN(now time.Time, n int) (time.Duration, error) {
	if n > l.burst {
		return 0, fmt.Errorf("ratelimit: WaitN %d exceeds limiter's burst %d", n, l.burst)
	}

	_, _, tokens := l.adjust(now)

	// calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n)

	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = l.limit.durationFromTokens(-tokens)
	}

	maxWaitDuration := l.limit.durationFromTokens(float64(n))
	ok := n <= l.burst && waitDuration <= maxWaitDuration

	if ok {
		return waitDuration, nil
	}
	return 0, fmt.Errorf("ratelimit: WaitN %d exceeds maximum wait duration", n)
}

type reservation struct {
	ok        bool
	tokens    int
	timeToAct time.Time
}

func (l *Limiter) reserveN(now time.Time, n int, maxWaitDuration time.Duration) reservation {
	l.mu.Lock()

	noow, last, tokens := l.adjust(now)

	// calculate the remaining number of tokens resulting from the request.
	tokens -= float64(n)

	var waitDuration time.Duration
	if tokens < 0 {
		waitDuration = l.limit.durationFromTokens(-tokens)
	}

	ok := n <= l.burst && waitDuration <= maxWaitDuration

	r := reservation{
		ok: ok,
	}

	if ok {
		r.tokens = n
		r.timeToAct = noow.Add(waitDuration)
		l.last = noow
		l.tokens = tokens
		l.lastRequest = r.timeToAct
	} else {
		l.last = last
	}

	l.mu.Unlock()
	return r
}

// adjust calculates the updated state for Limiter resulting from the passage of time.
func (l *Limiter) adjust(now time.Time) (newNow, newLast time.Time, newTokens float64) {
	last := l.last
	if now.Before(last) {
		last = now
	}

	maxElapsed := l.limit.durationFromTokens(float64(l.burst) - l.tokens)
	elapsed := now.Sub(last)
	if elapsed > maxElapsed {
		elapsed = maxElapsed
	}

	delta := l.limit.tokensFromDuration(elapsed)
	tokens := l.tokens + delta
	if burst := float64(l.burst); tokens > burst {
		tokens = burst
	}

	return now, last, tokens
}

func (lmt Limit) durationFromTokens(tokens float64) time.Duration {
	seconds := tokens / float64(lmt)
	return time.Nanosecond * time.Duration(1e9*seconds)
}

func (lmt Limit) tokensFromDuration(d time.Duration) float64 {
	return d.Seconds() * float64(lmt)
}
