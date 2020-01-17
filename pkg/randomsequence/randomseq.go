// Package randomsequence implements quadratic residues based random sequence.
package randomsequence

// https://preshing.com/20121224/how-to-generate-a-sequence-of-unique-random-integers

// Random represents the random sequence.
type Random struct {
	index  uint32
	offset uint32
}

const (
	prime    uint32 = 4294967291
	primeBy2 uint32 = prime / 2
)

// New creates a random sequence with the seed provided.
func New(seedBase, seedOffset uint32) *Random {
	return &Random{
		index:  permute(permute(seedBase) + 0x682f0161),
		offset: permute(permute(seedOffset) + 0x46790905),
	}
}

// Next returns next random number.
func (r *Random) Next() uint32 {
	i := r.index
	r.index++
	return permute((permute(i) + r.offset) ^ 0x5bf03635)
}

func permute(x uint32) uint32 {
	if x >= prime {
		// The 5 integers in the range [4294967291, 2^32] are mapped to themselves.
		return x
	}

	residue := uint32(uint64(x) * uint64(x) % uint64(prime))
	if x > primeBy2 {
		return prime - residue
	}
	return residue
}
