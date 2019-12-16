// Package spinlock implements Spinlock.
package spinlock

import (
	"runtime"
	"sync/atomic"
)

// Locker is the Spinlock implementation.
type Locker struct {
	noCopy noCopy
	lock   uintptr
}

// New creates a Locker.
func New() *Locker {
	return &Locker{}
}

// Lock wait in a loop to acquire the spinlock.
func (l *Locker) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock release the spinlock.
func (l *Locker) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
