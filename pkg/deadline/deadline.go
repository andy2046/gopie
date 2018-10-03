// Package deadline implements Deadline pattern.
package deadline

import (
	"errors"
	"time"
)

// ErrTimeout is the error for deadline timeout.
var ErrTimeout = errors.New("time out executing function")

// Deadline represents the deadline.
type Deadline struct {
	timeout time.Duration
}

// New returns a new Deadline with the provided timeout.
func New(timeout time.Duration) *Deadline {
	return &Deadline{
		timeout: timeout,
	}
}

// Go executes the provided function with a done channel as parameter to signal the timeout.
func (d *Deadline) Go(fn func(<-chan struct{}) error) error {
	result, done := make(chan error), make(chan struct{})

	go func() {
		result <- fn(done)
	}()

	select {
	case ret := <-result:
		return ret
	case <-time.After(d.timeout):
		close(done)
		return ErrTimeout
	}
}
