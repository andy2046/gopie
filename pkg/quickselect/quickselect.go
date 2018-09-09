// Package quickselect implements Quickselect.
package quickselect

// Select selects the kth element.
func Select(data []int, k int) int {
	return do(data, 0, len(data)-1, k)
}

func do(a []int, left, right, k int) int {
	for {
		if right <= left+1 {
			if right == left+1 && a[right] < a[left] {
				swap(a, left, right)
			}
			return a[k]
		}

		middle := (left + right) >> 1
		swap(a, middle, left+1)

		if a[left] > a[right] {
			swap(a, left, right)
		}

		if a[left+1] > a[right] {
			swap(a, left+1, right)
		}

		if a[left] > a[left+1] {
			swap(a, left, left+1)
		}

		i, j := left+1, right
		pivot := a[left+1]

		for {
			for ok := true; ok; ok = a[i] < pivot {
				i++
			}
			for ok := true; ok; ok = a[j] > pivot {
				j--
			}

			if j < i {
				break
			}

			swap(a, i, j)
		}

		a[left+1] = a[j]
		a[j] = pivot

		if j >= k {
			right = j - 1
		}

		if j <= k {
			left = i
		}
	}
}

func swap(a []int, i, j int) {
	a[i], a[j] = a[j], a[i]
}
