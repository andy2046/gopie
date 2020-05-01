package dll

import (
	"sync/atomic"
	"unsafe"
)

type (
	atomicMarkableReference struct {
		pair *pair
	}

	pair struct {
		reference *Element
		mark      bool
	}
)

func newAtomicMarkableReference(initialRef *Element, initialMark bool) *atomicMarkableReference {
	return &atomicMarkableReference{
		&pair{
			reference: initialRef,
			mark:      initialMark,
		},
	}
}

func (amr *atomicMarkableReference) getPair() *pair {
	return (*pair)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&amr.pair))))
}

func (amr *atomicMarkableReference) getReference() *Element {
	p := amr.getPair()
	return p.reference
}

func (amr *atomicMarkableReference) isMarked() bool {
	p := amr.getPair()
	return p.mark
}

func (amr *atomicMarkableReference) get() (bool, *Element) {
	p := amr.getPair()
	return p.mark, p.reference
}

func (amr *atomicMarkableReference) compareAndSet(expectedReference *Element,
	newReference *Element,
	expectedMark bool,
	newMark bool) bool {
	current := amr.getPair()
	val := &pair{newReference, newMark}

	return expectedReference == current.reference &&
		expectedMark == current.mark &&
		((newReference == current.reference &&
			newMark == current.mark) ||
			amr.casPair(current, val))
}

func (amr *atomicMarkableReference) tryMark(expectedReference *Element, newMark bool) bool {
	current := amr.getPair()
	val := &pair{expectedReference, newMark}
	return expectedReference == current.reference &&
		(newMark == current.mark ||
			amr.casPair(current, val))
}

func (amr *atomicMarkableReference) casPair(cmp *pair, val *pair) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&amr.pair)),
		unsafe.Pointer(cmp),
		unsafe.Pointer(val))
}
