// Package backoff implements Backoff pattern.
// http://www.awsarchitectureblog.com/2015/03/backoff.html
package backoff

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"math"
	"math/rand"
	"time"
)

type (
	// Strategy type.
	Strategy string

	// Backoff holds the backoff logic.
	Backoff struct {
		strategy iStrategy
	}

	// iStrategy defines the methodology for backing off.
	iStrategy interface {
		// Backoff returns the amount of time to wait before the next retry.
		Backoff(retries int) time.Duration
	}

	// exponential Strategy.
	exponential struct {
		base
	}

	// fullJitter Strategy.
	fullJitter struct {
		base
	}

	// equalJitter Strategy.
	equalJitter struct {
		base
	}

	// decorrelatedJitter Strategy.
	decorrelatedJitter struct {
		base
		sleep time.Duration
	}

	base struct {
		// capDelay is the upper bound of backoff delay.
		capDelay time.Duration
		// baseDelay is the initial backoff delay.
		baseDelay time.Duration
	}
)

// predefined Strategy types.
const (
	Exponential        Strategy = "Exponential"
	FullJitter         Strategy = "FullJitter"
	EqualJitter        Strategy = "EqualJitter"
	DecorrelatedJitter Strategy = "DecorrelatedJitter"
)

func init() {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("fail to seed math/rand pkg with crypto/rand random number generator")
	}
	rand.Seed(int64(binary.LittleEndian.Uint64(b[:])))
}

// New returns a Backoff instance with provided delay and strategy.
func New(baseDelay, capDelay time.Duration, strategy Strategy) *Backoff {
	if capDelay < baseDelay {
		panic("capDelay must be greater than baseDelay")
	}

	var (
		b = base{baseDelay: baseDelay, capDelay: capDelay}
		s iStrategy
	)

	switch strategy {
	case Exponential:
		s = &exponential{b}
	case FullJitter:
		s = &fullJitter{b}
	case EqualJitter:
		s = &equalJitter{b}
	case DecorrelatedJitter:
		s = &decorrelatedJitter{base: b, sleep: baseDelay}
	default:
		panic("unknown strategy " + strategy)
	}

	bk := &Backoff{s}
	return bk
}

// Backoff returns the amount of time to wait before the next retry.
// `retries` starts from zero.
func (bk *Backoff) Backoff(retries int) time.Duration {
	return bk.strategy.Backoff(retries)
}

func (b base) expo(retries int) time.Duration {
	// min(cap, pow(2, n)*base)
	pow2n := math.Pow(2, float64(retries))
	n := time.Duration(pow2n) * b.baseDelay
	return minDuration(b.capDelay, n)
}

func (ex *exponential) Backoff(retries int) time.Duration {
	return ex.expo(retries)
}

func (f *fullJitter) Backoff(retries int) time.Duration {
	n := f.expo(retries)
	return time.Duration(rand.Int63n(int64(n) + 1))
}

func (e *equalJitter) Backoff(retries int) time.Duration {
	n := e.expo(retries)
	return n/2 + time.Duration(rand.Int63n(1+int64(n)/2))
}

func (d *decorrelatedJitter) Backoff(retries int) time.Duration {
	if retries == 0 {
		d.sleep = d.baseDelay
	}

	// sleep = min(cap, random.uniform(base, sleep * 3))
	base, sleep3 := int64(d.baseDelay), int64(d.sleep)*3
	n := base + rand.Int63n(sleep3-base+1)
	d.sleep = minDuration(d.capDelay, time.Duration(n))
	return d.sleep
}

func minDuration(x, y time.Duration) time.Duration {
	if x < y {
		return x
	}
	return y
}
