// Package quicksort implements Quicksort.
package quicksort

// Sort sorts the slice.
func Sort(data []int) {
	sort(data, 0, len(data)-1)
}

func sort(c []int, start, end int) {
	if end <= start {
		return
	}
	i, j := start, end+1
	// ensure: c[start] <= c[start+1] <= c[end]
	if c[start] > c[end] {
		c[start], c[end] = c[end], c[start]
	}
	if c[start+1] > c[end] {
		c[start+1], c[end] = c[end], c[start+1]
	}
	if c[start] > c[start+1] {
		c[start], c[start+1] = c[start+1], c[start]
	}
	comp := c[start]
	for {
		for ok := true; ok; ok = c[i] < comp {
			i++
		}
		for ok := true; ok; ok = c[j] > comp {
			j--
		}
		if j <= i {
			break
		}
		c[i], c[j] = c[j], c[i]
	}
	c[start], c[j] = c[j], c[start]
	sort(c, start, j-1)
	sort(c, j+1, end)
}
