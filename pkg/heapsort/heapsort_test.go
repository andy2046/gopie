package heapsort_test

import (
	"testing"
	"time"

	"math/rand"

	"github.com/andy2046/gopie/pkg/heapsort"
)

func TestHeapsort(t *testing.T) {
	size := 20
	max := 999
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Intn(max))
	}
	t.Logf("before -> %v", s)
	heapsort.Sort(s)
	t.Logf("after -> %v", s)
}

func BenchmarkHeapsort(b *testing.B) {
	size := 1000000
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Int())
	}
	b.ResetTimer()
	heapsort.Sort(s)
}
