// Package bloom implements a Bloom filter.
package bloom

import (
	"math"
)

type (
	// Bloom is the standard bloom filter.
	Bloom interface {
		Add([]byte)
		AddString(string)
		Exist([]byte) bool
		ExistString(string) bool
		FalsePositive() float64
		GuessFalsePositive(uint64) float64
		M() uint64
		K() uint64
		N() uint64
		Clear()
	}

	// CountingBloom is the bloom filter which allows deletion of entries.
	// Take note that an 16-bit counter is maintained for each entry.
	CountingBloom interface {
		Bloom
		Remove([]byte)
		RemoveString(string)
	}

	bloomFilter struct {
		bitmap []uint16 // bloom filter counter
		k      uint64   // number of hash functions
		n      uint64   // number of elements in the bloom filter
		m      uint64   // size of the bloom filter bits
		shift  uint8    // the shift to get high/low bit fragments
	}
)

const (
	ln2                  float64 = 0.6931471805599453 // math.Log(2)
	maxCountingBloomSize uint64  = 1 << 37            // to avoid panic: makeslice: len out of range
	maxCounter           uint16  = 65535
)

// New creates counting bloom filter based on the provided m/k.
// m is the size of bloom filter bits.
// k is the number of hash functions.
func New(m, k uint64) CountingBloom {
	mm, exponent := adjustM(m)
	return &bloomFilter{
		bitmap: make([]uint16, mm),
		m:      mm - 1, // x % 2^i = x & (2^i - 1)
		k:      k,
		shift:  64 - exponent,
	}
}

// NewGuess estimates m/k based on the provided n/p then creates counting bloom filter.
// n is the estimated number of elements in the bloom filter.
// p is the false positive probability.
func NewGuess(n uint64, p float64) CountingBloom {
	m, k := Guess(n, p)
	return New(m, k)
}

// Guess estimates m/k based on the provided n/p.
func Guess(n uint64, p float64) (m, k uint64) {
	mm := math.Ceil(-1 * float64(n) * math.Log(p) / math.Pow(ln2, 2))
	kk := math.Ceil(ln2 * mm / float64(n))
	m, k = uint64(mm), uint64(kk)
	return
}

func (bf *bloomFilter) Add(entry []byte) {
	hash := sipHash(entry)
	h := hash >> bf.shift
	l := hash << bf.shift >> bf.shift
	var idx uint64
	for i := uint64(0); i < bf.k; i++ {
		idx = (h + i*l) & bf.m
		// avoid overflow
		if bf.bitmap[idx] < maxCounter {
			bf.bitmap[idx]++
		}
	}
	bf.n++
}

func (bf *bloomFilter) AddString(entry string) {
	bf.Add([]byte(entry))
}

func (bf *bloomFilter) Remove(entry []byte) {
	hash := sipHash(entry)
	h := hash >> bf.shift
	l := hash << bf.shift >> bf.shift
	var idx uint64
	for i := uint64(0); i < bf.k; i++ {
		idx = (h + i*l) & bf.m
		// avoid overflow
		if bf.bitmap[idx] > 0 {
			bf.bitmap[idx]--
		}
	}
	bf.n--
}

func (bf *bloomFilter) RemoveString(entry string) {
	bf.Remove([]byte(entry))
}

func (bf *bloomFilter) Exist(entry []byte) bool {
	hash := sipHash(entry)
	h := hash >> bf.shift
	l := hash << bf.shift >> bf.shift
	var idx uint64
	for i := uint64(0); i < bf.k; i++ {
		idx = (h + i*l) & bf.m
		if bf.bitmap[idx] == 0 {
			return false
		}
	}

	return true
}

func (bf *bloomFilter) ExistString(entry string) bool {
	return bf.Exist([]byte(entry))
}

func (bf *bloomFilter) FalsePositive() float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/float64(bf.m))),
		float64(bf.k))
}

func (bf *bloomFilter) GuessFalsePositive(n uint64) float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*n)/float64(bf.m))),
		float64(bf.k))
}

func (bf *bloomFilter) M() uint64 {
	return bf.m + 1
}

func (bf *bloomFilter) K() uint64 {
	return bf.k
}

func (bf *bloomFilter) N() uint64 {
	return bf.n
}

func (bf *bloomFilter) Clear() {
	for i := range bf.bitmap {
		bf.bitmap[i] = 0
	}
	bf.n = 0
}

func adjustM(x uint64) (m uint64, exponent uint8) {
	if x < 512 {
		x = 512
	}
	m = uint64(1)
	for m < x && m < maxCountingBloomSize {
		m <<= 1
		exponent++
	}
	return m, exponent
}
