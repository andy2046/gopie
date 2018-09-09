// Package mergesort implements Mergesort.
package mergesort

// Sort sorts the slice.
func Sort(data []int) {
	aux := make([]int, len(data))
	sort(data, aux)
}

func sort(a, aux []int) {
	n := len(a)
	for size := 1; size < n; size = 2 * size {
		for lo := 0; lo < n-size; lo += 2 * size {
			merge(a, aux, lo, lo+size-1, min(lo+size+size-1, n-1))
		}
	}
}

func merge(a, aux []int, lo, mid, hi int) {
	for k := lo; k <= hi; k++ {
		aux[k] = a[k]
	}

	i, j := lo, mid+1

	for k := lo; k <= hi; k++ {
		if i > mid {
			a[k] = aux[j]
			j++
		} else if j > hi {
			a[k] = aux[i]
			i++
		} else if aux[j] <= aux[i] {
			a[k] = aux[j]
			j++
		} else {
			a[k] = aux[i]
			i++
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
