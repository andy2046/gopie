// Package semaphore provides a semaphore implementation.
package semaphore

import (
	"context"
)

type (
	// Semaphore is the semaphore implementation.
	Semaphore struct {
		cur chan struct{}
	}

	// ISemaphore is the semaphore interface.
	ISemaphore interface {
		Acquire(context.Context) error
		Release(context.Context) error
	}
)

// Acquire acquires the semaphore.
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.cur <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Release releases the semaphore.
func (s *Semaphore) Release(ctx context.Context) error {
	select {
	case _ = <-s.cur:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// New creates a new semaphore with given maximum concurrent access.
func New(n int) ISemaphore {
	if n <= 0 {
		panic("the number of max concurrent access must be positive int")
	}
	return &Semaphore{
		cur: make(chan struct{}, n),
	}
}
