package spsc

import (
	"unsafe"
)

type iface struct {
	t, d unsafe.Pointer
}

func extractptr(i interface{}) unsafe.Pointer {
	return (*iface)(unsafe.Pointer(&i)).d
}

func inject(i interface{}, ptr unsafe.Pointer) {
	var v = (*unsafe.Pointer)((*iface)(unsafe.Pointer(&i)).d)
	*v = ptr
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
