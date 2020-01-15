// Package multilane implements a concurrent blocking multiset.
package multilane

import (
	"github.com/andy2046/gopie/pkg/spsc"
	"runtime"
)

const (
	defaultQSize uint32 = 1024
	defaultWidth uint32 = 8
)

var defaultQueFunc = func(size uint32) Queue {
	return spsc.New(size)
}

type (
	// Queue represents the queue with Get / Put / Close methods.
	Queue interface {
		Get(interface{}) bool
		Put(interface{})
		Close()
	}

	// Config for MultiLane.
	Config struct {
		LaneWidth uint32 // number of queue for MultiLane
		QueueSize uint32 // size of each queue in MultiLane
		New       func(uint32) Queue
	}

	// MultiLane represents the concurrent multiset.
	MultiLane struct {
		_        [spsc.CacheLinePadSize]byte
		lanes    []Queue
		laneMask int64
		putCh    chan int64
		getCh    chan int64
		queSize  uint32
		_        [4]byte
	}
)

// New create a new MultiLane.
func New(conf Config) *MultiLane {
	m := MultiLane{}
	qs, w, queFunc := defaultQSize, defaultWidth, defaultQueFunc
	if conf.New != nil {
		queFunc = conf.New
	}
	if conf.QueueSize > 7 {
		qs = conf.QueueSize
	}
	if conf.LaneWidth > 0 {
		w = conf.LaneWidth
	}
	m.queSize = nextPowerOf2(qs)
	m.lanes = make([]Queue, nextPowerOf2(w))
	m.laneMask = int64(len(m.lanes) - 1)
	m.putCh = make(chan int64, len(m.lanes))
	m.getCh = make(chan int64, len(m.lanes))
	for i := range m.lanes {
		m.lanes[i] = queFunc(m.queSize)
		m.putCh <- int64(i)
		m.getCh <- int64(i)
	}
	return &m
}

// GetLane get the value at the provided lane to given variable,
// blocking.
func (m *MultiLane) GetLane(lane uint32, i interface{}) bool {
	return m.lanes[int64(lane)&m.laneMask].Get(i)
}

// PutLane put given variable at the provided lane,
// blocking.
func (m *MultiLane) PutLane(lane uint32, i interface{}) {
	m.lanes[int64(lane)&m.laneMask].Put(i)
}

// Get the value at one of the lanes pointed by the cursor to given variable,
// blocking.
func (m *MultiLane) Get(i interface{}) bool {
	for {
		select {
		case curs := <-m.getCh:
			r := m.lanes[curs&m.laneMask].Get(i)
			m.getCh <- curs
			return r
		default:
			runtime.Gosched()
		}
	}
}

// Put given variable at one of the lanes pointed by the cursor,
// blocking.
func (m *MultiLane) Put(i interface{}) {
	for {
		select {
		case curs := <-m.putCh:
			m.lanes[curs&m.laneMask].Put(i)
			m.putCh <- curs
			return
		default:
			runtime.Gosched()
		}
	}
}

// Close the MultiLane, it shall NOT be called before `Put()`.
func (m *MultiLane) Close() {
	for _, l := range m.lanes {
		l.Close()
	}
}

// LaneWidth is the number of lanes.
func (m *MultiLane) LaneWidth() uint32 {
	return uint32(m.laneMask + 1)
}

// QueueSize is the size of the underlying queue.
func (m *MultiLane) QueueSize() uint32 {
	return m.queSize
}

func nextPowerOf2(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	return v + 1
}
