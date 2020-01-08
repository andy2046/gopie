// Package spsc implements a Single-Producer / Single-Consumer queue.
package spsc

import (
	"golang.org/x/sys/cpu"
	"runtime"
	"sync/atomic"
	"unsafe"
)

const (
	// CacheLinePadSize is the size of OS CacheLine.
	CacheLinePadSize = unsafe.Sizeof(cpu.CacheLinePad{})
	// DefaultMaxBatch is the default max batch size.
	DefaultMaxBatch uint32 = (1 << 8) - 1
)

// SPSC is the SPSC queue structure.
type SPSC struct {
	_          [CacheLinePadSize]byte
	writeCache int64
	_          [CacheLinePadSize - 8%CacheLinePadSize]byte
	writeIndex int64
	_          [CacheLinePadSize - 8%CacheLinePadSize]byte
	readIndex  int64
	_          [CacheLinePadSize - 8%CacheLinePadSize]byte
	readCache  int64
	_          [CacheLinePadSize - 8%CacheLinePadSize]byte
	done       int32
	_          [CacheLinePadSize - 4%CacheLinePadSize]byte
	data       []unsafe.Pointer
	mask       int64
	maxbatch   int64
}

// New create a new SPSC with bounded `size`.
func New(size uint32, batchSize ...uint32) *SPSC {
	sp := SPSC{}
	sp.data = make([]unsafe.Pointer, nextPowerOf2(size))
	sp.mask = int64(len(sp.data) - 1)
	bSize := DefaultMaxBatch
	if len(batchSize) > 0 && batchSize[0] != 0 {
		bSize = batchSize[0]
	}
	sp.maxbatch = int64(nextPowerOf2(bSize) - 1)
	return &sp
}

// Close the SPSC, it shall NOT be called before `Offer()` or `Put()`.
func (sp *SPSC) Close() {
	atomic.AddInt32(&sp.done, 1)
}

// Poll the value at the head of the queue to given variable,
// non-blocking, return false if the queue is empty / closed.
func (sp *SPSC) Poll(i interface{}) bool {
	readCache := sp.readCache
	if writeIndex := atomic.LoadInt64(&sp.writeIndex); readCache >= writeIndex {
		if readCache > sp.readIndex {
			atomic.StoreInt64(&sp.readIndex, readCache)
		}
		if atomic.LoadInt32(&sp.done) > 0 {
			// queue is closed
			return false
		}
		if writeIndex = atomic.LoadInt64(&sp.writeCache); readCache >= writeIndex {
			// queue is empty
			return false
		}
	}

	inject(i, sp.data[readCache&sp.mask])
	readCache++
	sp.readCache = readCache
	if readCache-sp.readIndex > sp.maxbatch {
		atomic.StoreInt64(&sp.readIndex, readCache)
	}
	return true
}

// Get the value at the head of the queue to given variable,
// blocking, return false if the queue is closed.
func (sp *SPSC) Get(i interface{}) bool {
	readCache := sp.readCache
	if writeIndex := atomic.LoadInt64(&sp.writeIndex); readCache >= writeIndex {
		if readCache > sp.readIndex {
			atomic.StoreInt64(&sp.readIndex, readCache)
		}
		for readCache >= writeIndex {
			if atomic.LoadInt32(&sp.done) > 0 {
				// queue is closed
				return false
			}
			sp.spin()
			writeIndex = atomic.LoadInt64(&sp.writeCache)
		}
	}

	inject(i, sp.data[readCache&sp.mask])
	readCache++
	sp.readCache = readCache
	if readCache-sp.readIndex > sp.maxbatch {
		atomic.StoreInt64(&sp.readIndex, readCache)
	}
	return true
}

// Offer given variable at the tail of the queue,
// non-blocking, return false if the queue is full.
func (sp *SPSC) Offer(i interface{}) bool {
	writeCache := sp.writeCache
	if masked, readIndex := writeCache-sp.mask,
		atomic.LoadInt64(&sp.readIndex); masked >= readIndex {
		if writeCache > sp.writeIndex {
			atomic.StoreInt64(&sp.writeIndex, writeCache)
		}
		if readIndex = atomic.LoadInt64(&sp.readIndex); masked >= readIndex {
			// queue is full
			return false
		}
	}

	sp.data[writeCache&sp.mask] = extractptr(i)
	writeCache = atomic.AddInt64(&sp.writeCache, 1)
	if writeCache-sp.writeIndex > sp.maxbatch {
		atomic.StoreInt64(&sp.writeIndex, writeCache)
	}
	return true
}

// Put given variable at the tail of the queue,
// blocking.
func (sp *SPSC) Put(i interface{}) {
	writeCache := sp.writeCache
	if masked, readIndex := writeCache-sp.mask,
		atomic.LoadInt64(&sp.readIndex); masked >= readIndex {
		if writeCache > sp.writeIndex {
			atomic.StoreInt64(&sp.writeIndex, writeCache)
		}
		for masked >= readIndex {
			sp.spin()
			readIndex = atomic.LoadInt64(&sp.readIndex)
		}
	}

	sp.data[writeCache&sp.mask] = extractptr(i)
	writeCache = atomic.AddInt64(&sp.writeCache, 1)
	if writeCache-sp.writeIndex > sp.maxbatch {
		atomic.StoreInt64(&sp.writeIndex, writeCache)
	}
}

func (sp *SPSC) spin() {
	runtime.Gosched()
}
