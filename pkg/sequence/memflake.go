package sequence

import (
	"errors"
	"sync"
	"time"
)

// MemFlake is an implementation of in-memory Iceflake.
type MemFlake struct {
	sync.Mutex
	startTime      time.Time
	bitLenSequence uint8
	elapsedTime    uint64
	startEpoch     int64
	machineID      uint64
}

var (
	errLessThanOne          = errors.New("`n` should not be less than 1")
	errTimeDrift            = errors.New("time drifts too much")
	_              Iceflake = &MemFlake{}
)

// NewMemFlake creates a MemFlake.
func NewMemFlake(startTime time.Time, bitLenSequence uint8, machineID uint64) *MemFlake {
	if startTime.After(time.Now()) {
		// no future time
		return nil
	}
	if startTime.IsZero() {
		startTime = time.Date(2019, 10, 9, 0, 0, 0, 0, time.UTC)
	}
	startEpoch := toIceflakeTime(startTime)
	if startEpoch < 0 {
		// should not be too far away from now
		// should be after 1970/1/1
		return nil
	}
	if bitLenSequence < 10 || bitLenSequence > 21 {
		// keep control
		bitLenSequence = 18
	}

	return &MemFlake{
		startTime:      startTime,
		bitLenSequence: bitLenSequence,
		startEpoch:     startEpoch,
		machineID:      machineID,
	}
}

// Next ...
func (m *MemFlake) Next() (uint64, error) {
	return m.NextN(1)
}

// NextN ...
func (m *MemFlake) NextN(n int) (uint64, error) {
	if n < 1 {
		return 0, errLessThanOne
	}

	current := currentElapsedTime(m.startEpoch)
	if current < 0 {
		return 0, errTimeDrift
	}

	nextTime, un := uint64(m.toID(current)), uint64(n)
	m.Lock()
	defer m.Unlock()
	// [elapsedTime - n + 1, elapsedTime] will be returned
	if m.elapsedTime < nextTime {
		m.elapsedTime = nextTime + un - 1
	} else {
		m.elapsedTime += un
	}
	// return the first available timestamp
	return m.elapsedTime - un + 1, nil
}

// MachineID ...
func (m *MemFlake) MachineID() uint64 {
	return m.machineID
}

// StartTime ...
func (m *MemFlake) StartTime() time.Time {
	return m.startTime
}

// BitLenSequence ...
func (m *MemFlake) BitLenSequence() uint8 {
	return m.bitLenSequence
}

func (m *MemFlake) toID(t int64) int64 {
	return t << m.bitLenSequence
}

func toIceflakeTime(t time.Time) int64 {
	return t.UTC().Unix() * 1000 // in milliseconds
}

func currentElapsedTime(startEpoch int64) int64 {
	return toIceflakeTime(time.Now()) - startEpoch
}
