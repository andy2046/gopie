// Package rwspinlock implements RWSpinLock.
package rwspinlock

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// RWLocker is the RWSpinLock implementation.
type RWLocker struct {
	noCopy noCopy
	// the writer bit is placed on the MSB, allowing 2^31 readers.
	// using recursion for the readers lock should be done with great caution,
	// if the lock is acquired for a thread and another thread acquire the writer lock,
	// reader recursion will lead to a deadlock.
	lock uint32
}

// New creates a RWLocker.
func New() *RWLocker {
	return &RWLocker{}
}

// RLock locks rw for reading.
func (rw *RWLocker) RLock() {
	for {
		// wait until there is no active writer.
		for atomic.LoadUint32(&rw.lock)&0x80000000 != 0 {
			runtime.Gosched()
		}

		o := atomic.LoadUint32(&rw.lock) & 0x7fffffff
		if atomic.CompareAndSwapUint32(&rw.lock, o, o+1) {
			return
		}
	}
}

// RUnlock undoes a single RLock call.
func (rw *RWLocker) RUnlock() {
	atomic.AddUint32(&rw.lock, ^uint32(0))
}

// TryRLock try to locks rw for reading,
// it returns true if succeed, false otherwise.
func (rw *RWLocker) TryRLock() bool {
	if atomic.LoadUint32(&rw.lock)&0x80000000 != 0 {
		return false
	}

	o := atomic.LoadUint32(&rw.lock) & 0x7fffffff
	return atomic.CompareAndSwapUint32(&rw.lock, o, o+1)
}

// Lock locks rw for writing.
func (rw *RWLocker) Lock() {
	for {
		// wait until there is no active writer.
		for atomic.LoadUint32(&rw.lock)&0x80000000 != 0 {
			runtime.Gosched()
		}

		o := atomic.LoadUint32(&rw.lock) & 0x7fffffff
		n := o | 0x80000000

		if atomic.CompareAndSwapUint32(&rw.lock, o, n) {
			// wait for active readers to release locks.
			for atomic.LoadUint32(&rw.lock)&0x7fffffff != 0 {
				runtime.Gosched()
			}
			return
		}
	}
}

// Unlock unlocks rw for writing.
func (rw *RWLocker) Unlock() {
	if atomic.LoadUint32(&rw.lock) != 0x80000000 {
		panic("Unlock")
	}
	atomic.StoreUint32(&rw.lock, 0)
}

// TryLock try to locks rw for writing,
// it returns true if succeed, false otherwise.
func (rw *RWLocker) TryLock() bool {
	if atomic.LoadUint32(&rw.lock)&0x80000000 != 0 {
		return false
	}

	o := atomic.LoadUint32(&rw.lock) & 0x7fffffff
	n := o | 0x80000000

	if !atomic.CompareAndSwapUint32(&rw.lock, o, n) {
		return false
	}

	// wait for active readers to release locks.
	for atomic.LoadUint32(&rw.lock)&0x7fffffff != 0 {
		runtime.Gosched()
	}
	return true
}

// IsLocked returns true if there is active writer, false otherwise.
func (rw *RWLocker) IsLocked() bool {
	return atomic.LoadUint32(&rw.lock)&0x80000000 != 0
}

// IsRLocked returns true if there is active reader, false otherwise.
func (rw *RWLocker) IsRLocked() bool {
	return atomic.LoadUint32(&rw.lock)&0x7fffffff != 0
}

// RLocker returns a sync.Locker interface implementation by calling rw.RLock and rw.RUnlock.
func (rw *RWLocker) RLocker() sync.Locker {
	return (*rlocker)(rw)
}

type rlocker RWLocker

func (r *rlocker) Lock()   { (*RWLocker)(r).RLock() }
func (r *rlocker) Unlock() { (*RWLocker)(r).RUnlock() }

// noCopy may be embedded into structs which must not be copied
// after the first use.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
