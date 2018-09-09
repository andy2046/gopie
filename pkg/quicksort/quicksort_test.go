package quicksort_test

import (
	"testing"
	"time"

	"math/rand"

	"github.com/andy2046/gopie/pkg/quicksort"
)

func TestQuicksort(t *testing.T) {
	size := 20
	max := 999
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Intn(max))
	}
	t.Logf("before -> %v", s)
	quicksort.Sort(s)
	t.Logf("after -> %v", s)
}

func BenchmarkQuicksort(b *testing.B) {
	size := 1000000
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Int())
	}
	b.ResetTimer()
	quicksort.Sort(s)
}
