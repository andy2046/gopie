// Package batch implements batch utility.
package batch

import (
	"log"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Batch process batch of work by time duration or batch count.
type Batch struct {
	timeout time.Duration
	lock    sync.Mutex
	once    sync.Once
	fn      func([]interface{})
	count   uint32
	temp    uint32
	input   chan interface{}
	notify  chan struct{}
	closer  chan struct{}
}

// New returns a Batch with time duration, batch count and batch function provided.
func New(timeout time.Duration, count int, fn func([]interface{})) *Batch {
	if count < 1 || count > math.MaxUint32 {
		panic("invalid batch count")
	}

	b := &Batch{
		timeout: timeout,
		fn:      fn,
		count:   uint32(count),
		notify:  make(chan struct{}),
		closer:  make(chan struct{}),
		input:   make(chan interface{}, count),
	}

	go b.batch()
	return b
}

// Batch batches the given item.
func (b *Batch) Batch(item interface{}) {
	b.input <- item

	if atomic.AddUint32(&b.temp, 1) == b.count {
		atomic.StoreUint32(&b.temp, 0)
		b.notify <- struct{}{}
	}
}

// Close the batch.
func (b *Batch) Close() {
	b.once.Do(func() {
		close(b.closer)
	})
}

func (b *Batch) batch() {
	var (
		batch []interface{}
		ch    = make(chan struct{})
		h     = func() {
			d := append(batch[:0:0], batch...)
			batch = batch[:0]
			go func() {
				defer func() {
					if r := recover(); r != nil {
						var buf [4096]byte
						n := runtime.Stack(buf[:], false)
						log.Println(r)
						log.Println(string(buf[:n]))
					}
				}()

				b.fn(d)
			}()

			go b.counter(ch)
		}
	)

	go b.counter(ch)

	for {
		select {
		case <-b.closer:
			return
		case <-ch:
			h()
			break
		default:
		}

		select {
		case <-b.closer:
			return
		case <-ch:
			h()
			break
		case item := <-b.input:
			batch = append(batch, item)
			continue
		}
	}
}

func (b *Batch) counter(ch chan struct{}) {
	timer := time.NewTimer(b.timeout)

	select {
	case <-b.closer:
		timer.Stop()
		return
	case <-timer.C:
		break
	case <-b.notify:
		break
	}

	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}

	select {
	case <-b.notify:
	default:
	}

	ch <- struct{}{}
}
