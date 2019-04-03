// Package hyperloglog implements HyperLogLog cardinality estimation.
package hyperloglog

import (
	"errors"
	"hash"
	"hash/fnv"
	"math"
)

// HyperLogLog probabilistic data struct for cardinality estimation.
type HyperLogLog struct {
	registers []uint8     // registers bucket
	m         uint        // number of registers
	b         uint32      // number of bits to find registers bucket number
	alpha     float64     // bias-correction constant
	hash      hash.Hash32 // hash function
}

const (
	exp32    float64 = 4294967296
	negexp32 float64 = -4294967296
	alpha16  float64 = 0.673
	alpha32  float64 = 0.697
	alpha64  float64 = 0.709
)

// New creates a new HyperLogLog with `m` registers bucket.
// `m` should be a power of two.
func New(m uint) (*HyperLogLog, error) {
	if (m & (m - 1)) != 0 {
		m = adjustM(m)
	}

	return &HyperLogLog{
		registers: make([]uint8, m),
		m:         m,
		b:         uint32(math.Ceil(math.Log2(float64(m)))),
		alpha:     calculateAlpha(m),
		hash:      fnv.New32(),
	}, nil
}

// NewGuess creates a new HyperLogLog within the given standard error.
func NewGuess(stdErr float64) (*HyperLogLog, error) {
	m := math.Pow(1.04/stdErr, 2)
	return New(uint(math.Pow(2, math.Ceil(math.Log2(m)))))
}

// Add adds the data to the set.
func (h *HyperLogLog) Add(data []byte) {
	var (
		hash = h.calculateHash(data)
		k    = 32 - h.b
		r    = calculateConsecutiveZeros(hash, k)
		j    = hash >> uint(k)
	)

	if r > h.registers[j] {
		h.registers[j] = r
	}
}

// Count returns the estimated cardinality of the set.
func (h *HyperLogLog) Count() uint64 {
	sum, m := 0.0, float64(h.m)
	for _, rv := range h.registers {
		sum += 1.0 / math.Pow(2.0, float64(rv))
	}
	estimate := h.alpha * m * m / sum
	if estimate <= 5.0/2.0*m {
		// Small range correction
		v := 0
		for _, r := range h.registers {
			if r == 0 {
				v++
			}
		}
		if v > 0 {
			estimate = m * math.Log(m/float64(v))
		}
	} else if estimate > 1.0/30.0*exp32 {
		// Large range correction
		estimate = negexp32 * math.Log(1-estimate/exp32)
	}
	return uint64(estimate)
}

// Merge combines the HyperLogLog with the other.
func (h *HyperLogLog) Merge(other *HyperLogLog) error {
	if h.m != other.m {
		return errors.New("registers bucket number must match")
	}

	for j, r := range other.registers {
		if r > h.registers[j] {
			h.registers[j] = r
		}
	}

	return nil
}

// Reset restores the HyperLogLog to its original state.
func (h *HyperLogLog) Reset() {
	h.registers = make([]uint8, h.m)
}

// SetHash sets the hashing function.
func (h *HyperLogLog) SetHash(hasher hash.Hash32) {
	h.hash = hasher
}

func (h *HyperLogLog) calculateHash(data []byte) uint32 {
	h.hash.Reset()
	h.hash.Write(data)
	sum := h.hash.Sum32()
	return sum
}

func calculateAlpha(m uint) float64 {
	var a float64
	switch m {
	case 16:
		a = alpha16
	case 32:
		a = alpha32
	case 64:
		a = alpha64
	default:
		a = 0.7213 / (1.0 + 1.079/float64(m))
	}
	return a
}

// calculateConsecutiveZeros calculates the position of the rightmost 1-bit.
func calculateConsecutiveZeros(val, max uint32) uint8 {
	r := uint32(1)
	for val&1 == 0 && r <= max {
		r++
		val >>= 1
	}
	return uint8(r)
}

func adjustM(x uint) uint {
	m := uint(1)
	for m < x {
		m <<= 1
	}
	return m
}
