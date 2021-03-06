// Package spinlock implements Spinlock.
package spinlock

import (
	"runtime"
	"sync"
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

// TryLock try to acquire the spinlock,
// it returns true if succeed, false otherwise.
func (l *Locker) TryLock() bool {
	return atomic.CompareAndSwapUintptr(&l.lock, 0, 1)
}

// IsLocked returns true if locked, false otherwise.
func (l *Locker) IsLocked() bool {
	return atomic.LoadUintptr(&l.lock) == 1
}

// Locker returns a sync.Locker interface implementation.
func (l *Locker) Locker() sync.Locker {
	return l
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
