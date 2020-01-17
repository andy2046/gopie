package randomsequence_test

import (
	bt "github.com/andy2046/bitmap"
	. "github.com/andy2046/gopie/pkg/randomsequence"
	"testing"
	"time"
)

func TestUnique(t *testing.T) {
	size := uint64(^uint32(0) >> 10) // by right the size is 4294967295 but it is too slow
	t.Log("sample size is", size)
	btmap := bt.New(size + 1)
	seed := uint32(time.Now().UTC().UnixNano())
	rnd := New(seed, seed+1)
	for i := uint64(0); i <= size; i++ {
		n := rnd.Next()
		if btmap.GetBit(uint64(n)) {
			t.Fatal("dup")
		}
		btmap.SetBit(uint64(n), true)
	}
}
