package bloom

import (
	"math"

	"github.com/andy2046/bitmap"
)

type (
	scalableBloomFilter struct {
		filterz   []*bloomFilterBit // bloom filters list
		count     uint64            // number of elements in the bloom filter
		n         uint64            // estimated number of elements
		p         float64           // target False Positive rate
		r         float64           // optimal tightening ratio
		fillRatio float64           // fill ratio
	}
)

const (
	rDefault  float64 = 0.8
	fillRatio float64 = 0.5
)

// NewS creates scalable bloom filter based on the provided fpRate.
// fpRate is the target False Positive probability.
func NewS(fpRate float64) Bloom {
	return NewSGuess(10000, fpRate, rDefault)
}

// NewSGuess estimates m/k based on the provided n/p then creates scalable bloom filter.
// n is the estimated number of elements in the bloom filter.
// p is the false positive probability.
// r is the optimal tightening ratio.
func NewSGuess(n uint64, p, r float64) Bloom {
	m, k := Guess(n, p)
	mm, exponent := adjustM(m)

	sBF := scalableBloomFilter{
		filterz:   make([]*bloomFilterBit, 0, 1),
		r:         r,
		fillRatio: fillRatio,
		p:         p,
		n:         n,
	}

	sBF.filterz = append(sBF.filterz, &bloomFilterBit{
		bitmap: bitmap.New(mm),
		m:      mm - 1, // x % 2^i = x & (2^i - 1)
		k:      k,
		shift:  64 - exponent,
	})
	return &sBF
}

func (bf *scalableBloomFilter) Add(entry []byte) {
	idx := len(bf.filterz) - 1
	if bf.filterz[idx].estimatedFillRatio() >= bf.fillRatio {
		fp := bf.p * math.Pow(bf.r, float64(len(bf.filterz)))
		m, k := Guess(bf.n, fp)
		mm, exponent := adjustM(m)
		bf.filterz = append(bf.filterz, &bloomFilterBit{
			bitmap: bitmap.New(mm),
			m:      mm - 1, // x % 2^i = x & (2^i - 1)
			k:      k,
			shift:  64 - exponent,
		})
		idx++
	}
	bf.filterz[idx].Add(entry)
	bf.count++
}

func (bf *scalableBloomFilter) AddString(entry string) {
	bf.Add([]byte(entry))
}

func (bf *scalableBloomFilter) Exist(entry []byte) bool {
	for _, f := range bf.filterz {
		if f.Exist(entry) {
			return true
		}
	}
	return false
}

func (bf *scalableBloomFilter) ExistString(entry string) bool {
	return bf.Exist([]byte(entry))
}

func (bf *scalableBloomFilter) FalsePositive() float64 {
	rez := 1.0
	for _, f := range bf.filterz {
		rez *= (1.0 - f.FalsePositive())
	}
	return 1.0 - rez
}

func (bf *scalableBloomFilter) GuessFalsePositive(n uint64) float64 {
	rez := 1.0
	for _, f := range bf.filterz {
		rez *= (1.0 - f.GuessFalsePositive(n))
	}
	return 1.0 - rez
}

func (bf *scalableBloomFilter) M() uint64 {
	m := uint64(0)
	for _, f := range bf.filterz {
		m += f.M()
	}
	return m
}

func (bf *scalableBloomFilter) K() uint64 {
	return bf.filterz[0].K()
}

func (bf *scalableBloomFilter) N() uint64 {
	return bf.count
}

func (bf *scalableBloomFilter) Clear() {
	for i := range bf.filterz {
		bf.filterz[i] = nil
	}
	bf.filterz = make([]*bloomFilterBit, 0, 1)
	m, k := Guess(bf.n, bf.p)
	mm, exponent := adjustM(m)
	bf.filterz = append(bf.filterz, &bloomFilterBit{
		bitmap: bitmap.New(mm),
		m:      mm - 1, // x % 2^i = x & (2^i - 1)
		k:      k,
		shift:  64 - exponent,
	})
	bf.count = 0
}
