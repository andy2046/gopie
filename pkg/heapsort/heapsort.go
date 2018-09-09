// Package heapsort implements Heapsort.
package heapsort

// Sort sorts the slice.
func Sort(data []int) {
	sort(data)
}

func sort(a []int) {
	n := len(a)
	buildMaxHeap(a, n)
	sortDown(a, n)
}

func buildMaxHeap(a []int, n int) {
	for k := n / 2; k >= 1; k-- {
		sink(a, k, n)
	}
}

func sortDown(a []int, n int) {
	for n > 1 {
		swap(a, 1, n)
		n--
		sink(a, 1, n)
	}
}

func sink(a []int, k, n int) {
	for 2*k <= n {
		j := 2 * k

		// is right key greater than left key
		if j < n && less(a, j, j+1) {
			j++
		}

		// when both right and left child are not greater than parent
		if !less(a, k, j) {
			break
		}

		// moves the greater key up
		swap(a, k, j)

		k = j
	}
}

func less(a []int, i, j int) bool {
	return a[i-1] < a[j-1]
}

func swap(a []int, i, j int) {
	a[i-1], a[j-1] = a[j-1], a[i-1]
}
