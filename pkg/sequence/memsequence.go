package sequence

import (
	"errors"
	"sync"
)

// MemSeq is an implementation of in-memory Sequencer.
type MemSeq struct {
	sync.Mutex
	current uint64
	machine uint64
}

var (
	errLessThanToo           = errors.New("`n` should not be less than 1")
	_              Sequencer = &MemSeq{}
)

// NewMemSeq creates a MemSeq.
func NewMemSeq(machineID uint64) *MemSeq {
	return &MemSeq{
		machine: machineID,
	}
}

// Next ...
func (m *MemSeq) Next() (uint64, error) {
	return m.NextN(1)
}

// NextN ...
func (m *MemSeq) NextN(n int) (uint64, error) {
	if n < 1 {
		return 0, errLessThanToo
	}
	m.Lock()
	r := m.current
	m.current += uint64(n)
	m.Unlock()
	return r, nil
}

// MachineID ...
func (m *MemSeq) MachineID() uint64 {
	return m.machine
}
