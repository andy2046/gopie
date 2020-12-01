package pushsum

import "sync"

type store struct {
	index map[uint64]uint64 // map key to ring index
	ring  []valueWeight
	sync.RWMutex
}

func newStore(len uint32) *store {
	return &store{
		index: make(map[uint64]uint64, len),
		ring:  make([]valueWeight, len, len),
	}
}

func (s *store) get(key uint64) (uint64, valueWeight, bool) {
	s.RLock()
	pos, ok := s.index[key]
	if !ok {
		// not exist
		s.RUnlock()
		return 0, valueWeight{}, false
	}

	v := s.ring[pos]
	if v.key != key {
		// k deleted from index?
		s.RUnlock()
		return 0, valueWeight{}, false
	}

	s.RUnlock()
	return pos, v, true
}

func (s *store) compareAndSet(idx uint64, vw valueWeight) bool {
	s.Lock()
	if v := s.ring[idx]; v.key != vw.key {
		s.Unlock()
		return false
	}

	s.ring[idx] = vw
	s.Unlock()
	return true
}

func (s *store) update(key uint64, value, weight float64, inTransition bool) bool {
	s.Lock()
	defer s.Unlock()

	i, ok := s.index[key]
	if !ok {
		return false
	}

	vw := s.ring[i]
	if vw.key == key {
		vw.value += value
		vw.weight += weight
		vw.inTransition = inTransition
		s.ring[i] = vw
	}

	return true
}

func (s *store) insert(idx uint64, vw valueWeight) {
	s.Lock()
	// remove index for existing vw
	e := s.ring[idx]
	delete(s.index, e.key)

	s.ring[idx] = vw
	s.index[vw.key] = idx
	s.Unlock()
}
