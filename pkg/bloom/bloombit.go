package bloom

import (
	"math"

	"github.com/andy2046/bitmap"
)

type (
	bloomFilterBit struct {
		bitmap *bitmap.Bitmap // bloom filter bitmap
		k      uint64         // number of hash functions
		n      uint64         // number of elements in the bloom filter
		m      uint64         // size of the bloom filter bits
		shift  uint8          // the shift to get high/low bit fragments
	}
)

// NewB creates standard bloom filter based on the provided m/k.
// m is the size of bloom filter bits.
// k is the number of hash functions.
func NewB(m, k uint64) Bloom {
	mm, exponent := adjustM(m)
	return &bloomFilterBit{
		bitmap: bitmap.New(mm),
		m:      mm - 1, // x % 2^i = x & (2^i - 1)
		k:      k,
		shift:  64 - exponent,
	}
}

// NewBGuess estimates m/k based on the provided n/p then creates standard bloom filter.
// n is the estimated number of elements in the bloom filter.
// p is the false positive probability.
func NewBGuess(n uint64, p float64) Bloom {
	m, k := Guess(n, p)
	return NewB(m, k)
}

func (bf *bloomFilterBit) Add(entry []byte) {
	hash := sipHash(entry)
	h := hash >> bf.shift
	l := hash << bf.shift >> bf.shift
	for i := uint64(0); i < bf.k; i++ {
		bf.bitmap.SetBit((h+i*l)&bf.m, true)
	}
	bf.n++
}

func (bf *bloomFilterBit) AddString(entry string) {
	bf.Add([]byte(entry))
}

func (bf *bloomFilterBit) Exist(entry []byte) bool {
	hash := sipHash(entry)
	h := hash >> bf.shift
	l := hash << bf.shift >> bf.shift

	for i := uint64(0); i < bf.k; i++ {
		if !bf.bitmap.GetBit((h + i*l) & bf.m) {
			return false
		}
	}

	return true
}

func (bf *bloomFilterBit) ExistString(entry string) bool {
	return bf.Exist([]byte(entry))
}

func (bf *bloomFilterBit) FalsePositive() float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*bf.n)/float64(bf.m))),
		float64(bf.k))
}

func (bf *bloomFilterBit) GuessFalsePositive(n uint64) float64 {
	return math.Pow((1 - math.Exp(-float64(bf.k*n)/float64(bf.m))),
		float64(bf.k))
}

func (bf *bloomFilterBit) M() uint64 {
	return bf.m + 1
}

func (bf *bloomFilterBit) K() uint64 {
	return bf.k
}

func (bf *bloomFilterBit) N() uint64 {
	return bf.n
}

func (bf *bloomFilterBit) Clear() {
	s := bf.bitmap.Size()
	for i := uint64(0); i < s; i++ {
		bf.bitmap.SetBit(i, false)
	}
	bf.n = 0
}

func (bf *bloomFilterBit) estimatedFillRatio() float64 {
	return 1 - math.Exp(-float64(bf.n)/math.Ceil(float64(bf.m)/float64(bf.k)))
}
