// Package subset implements deterministic subsetting.
package subset

/*
   https://landing.google.com/sre/book/chapters/load-balancing-datacenter.html
*/

import (
	"math/rand"
)

// Subset returns a subset of backends with size subsetSize.
func Subset(backends []string, clientID, subsetSize int) []string {

	subsetCount := len(backends) / subsetSize

	// Group clients into rounds; each round uses the same shuffled list:
	round := clientID / subsetCount

	r := rand.New(rand.NewSource(int64(round)))
	r.Shuffle(len(backends), func(i, j int) { backends[i], backends[j] = backends[j], backends[i] })

	// The subset id corresponding to the current client:
	subsetID := clientID % subsetCount

	start := subsetID * subsetSize
	return backends[start : start+subsetSize]
}
