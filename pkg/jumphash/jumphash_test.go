package jumphash_test

import (
	"strconv"
	"testing"

	. "github.com/andy2046/gopie/pkg/jumphash"
)

var jumphashTest = []struct {
	key      uint64
	buckets  int
	expected int
}{
	{1, 1, 0},
	{42, 57, 43},
	{0xDEAD10CC, 1, 0},
	{0xDEAD10CC, 666, 361},
	{256, 1024, 520},
	{0, -10, 0},
}

func TestHash(t *testing.T) {
	for _, v := range jumphashTest {
		h := Hash(v.key, v.buckets)
		if h != v.expected {
			t.Errorf("expected bucket for key=%d to be %d, got %d",
				v.key, v.expected, h)
		}
	}
}

var jumphashStringTest = []struct {
	key      string
	buckets  int
	expected int
}{
	{"localhost", 10, 6},
	{"中国", 10, 1},
}

func TestHashString(t *testing.T) {
	for _, v := range jumphashStringTest {
		h := HashString(v.key, v.buckets)
		if h != v.expected {
			t.Errorf("invalid bucket for key=%s, expected %d, got %d",
				strconv.Quote(v.key), v.expected, h)
		}
	}
}

func TestHasher(t *testing.T) {
	for _, v := range jumphashStringTest {
		hasher := New(v.buckets)
		h := hasher.Hash(v.key)
		if h != v.expected {
			t.Errorf("invalid bucket for key=%s, expected %d, got %d",
				strconv.Quote(v.key), v.expected, h)
		}
	}
}

func BenchmarkHashN100(b *testing.B)   { benchmarkHash(b, 100) }
func BenchmarkHashN1000(b *testing.B)  { benchmarkHash(b, 1000) }
func BenchmarkHashN10000(b *testing.B) { benchmarkHash(b, 10000) }

func benchmarkHash(b *testing.B, n int) {
	h := New(n)
	for i := 0; i < b.N; i++ {
		if x := h.Hash(strconv.Itoa(i)); x > n {
			b.Fatal("invalid hash:", x)
		}
	}
}
