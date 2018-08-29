package subset

import (
	"strconv"
	"testing"
)

func TestSubset(t *testing.T) {
	subsetSize := 10
	clientSize := 300
	backendSize := 300
	loopC := make([]struct{}, clientSize)
	loopB := make([]struct{}, backendSize)
	results := make(map[string]int)
	min, max := clientSize, 1

	clients := make([]int, 0, clientSize)
	for i := range loopC {
		clients = append(clients, i)
	}

	for i := range loopC {
		backends := make([]string, 0, backendSize)
		for i := range loopB {
			backends = append(backends, strconv.Itoa(i))
		}
		sets := Subset(backends, clients[i], subsetSize)
		for _, b := range sets {
			if _, ok := results[b]; !ok {
				results[b] = 1
			} else {
				results[b]++
			}
		}
	}

	t.Log("results:")
	for k, v := range results {
		t.Logf("backend %s -> client count %d", k, v)
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	t.Logf("min -> %d max -> %d (max-min)/min -> %d%%", min, max, 100*(max-min)/min)
}
