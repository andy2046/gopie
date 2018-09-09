package quickselect_test

import (
	"testing"
	"time"

	"math/rand"

	"github.com/andy2046/gopie/pkg/quickselect"
	"github.com/andy2046/gopie/pkg/quicksort"
)

func TestQuickselect(t *testing.T) {
	size := 20
	selected := 10
	max := 999
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Intn(max))
	}
	t.Logf("before -> %v", s)
	quickselect.Select(s, selected)
	t.Logf("after -> %v", s)
	se := s[selected-1]
	t.Logf("%dth selected -> %v", selected, se)
	quicksort.Sort(s)
	so := s[selected-1]
	t.Logf("sorted -> %v", s)
	if se != so {
		t.Fatalf("expected %d, got %d", so, se)
	}
}

func BenchmarkQuickselect(b *testing.B) {
	size := 1000000
	selected := size - 10
	s := make([]int, 0, size)
	rand.Seed(time.Now().UTC().UnixNano())
	for range make([]struct{}, size) {
		s = append(s, rand.Int())
	}
	b.ResetTimer()
	quickselect.Select(s, selected)
}
