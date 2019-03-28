// Package countminsketch implements Count-Min Sketch.
package countminsketch

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"math"
)

// CountMinSketch struct.
type CountMinSketch struct {
	matrix [][]uint64  // count matrix
	width  uint        // matrix width
	depth  uint        // matrix depth
	count  uint64      // total number of items added
	hash   hash.Hash64 // hash function
}

// For a sketch matrix w x d with total sum of all counts N,
// the estimate has error at most 2N/w, with probability at least 1-(1/2)^d.

// New returns new Count-Min Sketch with the given `width` and `depth`.
func New(width, depth uint) (*CountMinSketch, error) {
	if width < 1 || depth < 1 {
		return nil, errors.New("Dimensions must be positive")
	}

	matrix := make([][]uint64, depth)
	for i := uint(0); i < depth; i++ {
		matrix[i] = make([]uint64, width)
	}

	return &CountMinSketch{
		matrix: matrix,
		width:  width,
		depth:  depth,
		hash:   fnv.New64(),
	}, nil
}

// NewGuess returns new Count-Min Sketch with the given error rate `epsilon` and confidence `delta`.
func NewGuess(epsilon, delta float64) (*CountMinSketch, error) {
	if epsilon <= 0 || epsilon >= 1 {
		return nil, errors.New("epsilon must be in range (0, 1)")
	}
	if delta <= 0 || delta >= 1 {
		return nil, errors.New("delta must be in range (0, 1)")
	}

	width, depth := uint(math.Ceil(math.E/epsilon)),
		uint(math.Ceil(math.Log(1-delta)/math.Log(0.5)))

	return New(width, depth)
}

// Count returns the number of items added to the sketch.
func (c *CountMinSketch) Count() uint64 {
	return c.count
}

// Add add the `data` to the sketch. `count` default to 1.
func (c *CountMinSketch) Add(data []byte, count ...uint64) {
	cnt := uint64(1)
	if len(count) > 0 {
		cnt = count[0]
	}

	lower, upper := hashn(data, c.hash)

	for i := uint(0); i < c.depth; i++ {
		c.matrix[i][(uint(lower)+uint(upper)*i)%c.width] += cnt
	}

	c.count += cnt
}

// AddString add the `data` string to the sketch. `count` default to 1.
func (c *CountMinSketch) AddString(data string, count ...uint64) {
	c.Add([]byte(data), count...)
}

// Estimate estimate the frequency of the `data`.
func (c *CountMinSketch) Estimate(data []byte) uint64 {
	var (
		lower, upper = hashn(data, c.hash)
		count        uint64
	)

	for i := uint(0); i < c.depth; i++ {
		j := (uint(lower) + uint(upper)*i) % c.width
		if i == 0 || c.matrix[i][j] < count {
			count = c.matrix[i][j]
		}
	}

	return count
}

// EstimateString estimate the frequency of the `data` string.
func (c *CountMinSketch) EstimateString(data string) uint64 {
	return c.Estimate([]byte(data))
}

// Reset reset the sketch to its original state.
func (c *CountMinSketch) Reset() {
	matrix := make([][]uint64, c.depth)
	for i := uint(0); i < c.depth; i++ {
		matrix[i] = make([]uint64, c.width)
	}

	c.matrix = matrix
	c.count = 0
}

// Merge combines the sketch with another.
func (c *CountMinSketch) Merge(other *CountMinSketch) error {
	if c.depth != other.depth {
		return errors.New("matrix depth must match")
	}

	if c.width != other.width {
		return errors.New("matrix width must match")
	}

	for i := uint(0); i < c.depth; i++ {
		for j := uint(0); j < c.width; j++ {
			c.matrix[i][j] += other.matrix[i][j]
		}
	}

	c.count += other.count
	return nil
}

// Depth returns the matrix depth.
func (c *CountMinSketch) Depth() uint {
	return c.depth
}

// Width returns the matrix width.
func (c *CountMinSketch) Width() uint {
	return c.width
}

func hashn(data []byte, hasher hash.Hash64) (uint32, uint32) {
	hasher.Reset()
	hasher.Write(data)
	sum := hasher.Sum(nil)
	return binary.BigEndian.Uint32(sum[4:8]), binary.BigEndian.Uint32(sum[0:4])
}
