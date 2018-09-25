// Package breaker implements Circuit Breaker pattern.
package breaker

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// State is the type representing a state of CircuitBreaker.
type State int

const (
	// StateClosed represents Closed State.
	StateClosed State = iota
	// StateHalfOpen represents HalfOpen State.
	StateHalfOpen
	// StateOpen represents Open State.
	StateOpen
)

var (
	// ErrTooManyRequests is returned when the state is half open
	// and the requests count is more the maxRequests.
	ErrTooManyRequests = errors.New("too many requests, requests count is more the maxRequests in half open state")
	// ErrOpenState is returned when the state is open.
	ErrOpenState = errors.New("circuit breaker is open")
)

// String implements stringer interface.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return fmt.Sprintf("unknown state: %d", s)
	}
}

// Counts holds the numbers of requests and their successes/failures.
// CircuitBreaker clears the internal Counts either
// on the change of the state or at the closed-state intervals.
// Counts ignores the results of the requests sent before clearing.
type Counts struct {
	Requests             uint64
	TotalSuccesses       uint64
	TotalFailures        uint64
	ConsecutiveSuccesses uint64
	ConsecutiveFailures  uint64
}

func (c *Counts) onRequest() {
	c.Requests++
}

func (c *Counts) onSuccess() {
	c.TotalSuccesses++
	c.ConsecutiveSuccesses++
	c.ConsecutiveFailures = 0
}

func (c *Counts) onFailure() {
	c.TotalFailures++
	c.ConsecutiveFailures++
	c.ConsecutiveSuccesses = 0
}

func (c *Counts) clear() {
	c.Requests = 0
	c.TotalSuccesses = 0
	c.TotalFailures = 0
	c.ConsecutiveSuccesses = 0
	c.ConsecutiveFailures = 0
}

type (
	// Settings represents settings for CircuitBreaker.
	Settings struct {
		// Name is the name of the CircuitBreaker.
		Name string
		// MaxRequests is the maximum number of requests allowed to pass through
		// when the CircuitBreaker is half-open.
		// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
		MaxRequests uint64
		// Interval is the cyclic period of the closed state
		// for the CircuitBreaker to clear the internal Counts.
		// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
		Interval time.Duration
		// Timeout is the period of the open state,
		// after which the state of the CircuitBreaker becomes half-open.
		// If Timeout is 0, the timeout for the CircuitBreaker is 60 seconds.
		Timeout time.Duration
		// ShouldTrip is called with a copy of Counts whenever a request fails in the closed state.
		// If ShouldTrip returns true, the CircuitBreaker will be placed into the open state.
		// If ShouldTrip is nil, default ShouldTrip is used.
		// Default ShouldTrip returns true when the number of consecutive failures is more than 5.
		ShouldTrip func(counts Counts) bool
		// OnStateChange is called whenever the state of the CircuitBreaker changes.
		OnStateChange func(name string, from, to State)
	}

	// CircuitBreaker prevent an application repeatedly trying to execute an operation that is likely to fail.
	CircuitBreaker struct {
		name          string
		maxRequests   uint64
		interval      time.Duration
		timeout       time.Duration
		shouldTrip    func(counts Counts) bool
		onStateChange func(name string, from, to State)

		mutex      sync.Mutex
		state      State
		generation uint64
		counts     Counts
		expiry     time.Time
	}

	// Option applies settings to CircuitBreaker Settings.
	Option = func(*Settings) error
)

// DefaultSettings is the default CircuitBreaker Settings.
var DefaultSettings = Settings{
	Name:        "CircuitBreaker",
	MaxRequests: 1,
	Interval:    0,
	Timeout:     60 * time.Second,
	ShouldTrip: func(counts Counts) bool {
		return counts.ConsecutiveFailures > 5
	},
	OnStateChange: nil,
}

func setOption(s *Settings, options ...func(*Settings) error) error {
	for _, opt := range options {
		if err := opt(s); err != nil {
			return err
		}
	}
	return nil
}

// New returns a new CircuitBreaker with options applied.
func New(options ...Option) *CircuitBreaker {
	st := DefaultSettings
	err := setOption(&st, options...)
	if err != nil {
		log.Panicf("fail to apply Settings -> %v\n", err)
	}

	cb := &CircuitBreaker{
		name:          st.Name,
		maxRequests:   st.MaxRequests,
		interval:      st.Interval,
		timeout:       st.Timeout,
		shouldTrip:    st.ShouldTrip,
		onStateChange: st.OnStateChange,
	}

	cb.toNewGeneration(time.Now())
	return cb
}

// Name returns the name of the CircuitBreaker.
func (cb *CircuitBreaker) Name() string {
	return cb.name
}

// State returns the current state of the CircuitBreaker.
func (cb *CircuitBreaker) State() State {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, _ := cb.currentState(now)
	return state
}

// Execute runs the given request if the CircuitBreaker accepts it.
// Execute returns an error instantly if the CircuitBreaker rejects the request.
// Otherwise, Execute returns the result of the request.
// If a panic occurs in the request, the CircuitBreaker handles it as an error
// and causes the same panic again.
func (cb *CircuitBreaker) Execute(request func() (interface{}, error)) (interface{}, error) {
	generation, err := cb.beforeRequest()
	if err != nil {
		return nil, err
	}

	defer func() {
		e := recover()
		if e != nil {
			cb.afterRequest(generation, false)
			panic(e)
		}
	}()

	result, err := request()
	cb.afterRequest(generation, err == nil)
	return result, err
}

func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)

	if state == StateOpen {
		return generation, ErrOpenState
	} else if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequests {
		return generation, ErrTooManyRequests
	}

	cb.counts.onRequest()
	return generation, nil
}

func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	now := time.Now()
	state, generation := cb.currentState(now)
	if generation != before {
		return
	}

	if success {
		cb.onSuccess(state, now)
	} else {
		cb.onFailure(state, now)
	}
}

func (cb *CircuitBreaker) onSuccess(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onSuccess()
	case StateHalfOpen:
		cb.counts.onSuccess()
		if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
			cb.setState(StateClosed, now)
		}
	}
}

func (cb *CircuitBreaker) onFailure(state State, now time.Time) {
	switch state {
	case StateClosed:
		cb.counts.onFailure()
		if cb.shouldTrip(cb.counts) {
			cb.setState(StateOpen, now)
		}
	case StateHalfOpen:
		cb.setState(StateOpen, now)
	}
}

func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
	switch cb.state {
	case StateClosed:
		if !cb.expiry.IsZero() && cb.expiry.Before(now) {
			cb.toNewGeneration(now)
		}
	case StateOpen:
		if cb.expiry.Before(now) {
			cb.setState(StateHalfOpen, now)
		}
	}
	return cb.state, cb.generation
}

func (cb *CircuitBreaker) setState(state State, now time.Time) {
	if cb.state == state {
		return
	}

	prev := cb.state
	cb.state = state

	cb.toNewGeneration(now)

	if cb.onStateChange != nil {
		cb.onStateChange(cb.name, prev, state)
	}
}

func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
	cb.generation++
	cb.counts.clear()

	var zero time.Time
	switch cb.state {
	case StateClosed:
		if cb.interval == 0 {
			cb.expiry = zero
		} else {
			cb.expiry = now.Add(cb.interval)
		}
	case StateOpen:
		cb.expiry = now.Add(cb.timeout)
	default: // for StateHalfOpen
		cb.expiry = zero
	}
}
