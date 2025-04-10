package slice

import "sync"

// `Sort` performs a merge sort on the given slice using the provided comparison function.
// The `compare` function should return true if `a` should come before `b`.
func Sort[T any](slice []T, compare func(a, b T) bool) {
	if len(slice) <= 1 {
		return
	}

	var mergeSort func(arr []T, compare func(a, b T) bool) []T
	var merge func(left, right []T, compare func(a, b T) bool) []T

	mergeSort = func(arr []T, compare func(a, b T) bool) []T {
		if len(arr) <= 1 {
			return arr
		}

		mid := len(arr) / 2
		left := mergeSort(arr[:mid], compare)
		right := mergeSort(arr[mid:], compare)
		return merge(left, right, compare)
	}

	merge = func(left, right []T, compare func(a, b T) bool) []T {
		result := make([]T, 0, len(left)+len(right))
		i, j := 0, 0

		for i < len(left) && j < len(right) {
			if compare(left[i], right[j]) {
				result = append(result, left[i])
				i++
			} else {
				result = append(result, right[j])
				j++
			}
		}

		// Append remaining elements
		result = append(result, left[i:]...)
		result = append(result, right[j:]...)
		return result
	}

	// Sort and copy back the result
	sorted := mergeSort(slice, compare)
	copy(slice, sorted)
}

// `ParallelSort` performs a merge sort with parallel on the given slice using the provided comparison function.
// The `compare` function should return true if `a` should come before `b`.
// The `level` is depth of splitting through thread.
func ParallelSort[T any](slice []T, compare func(a, b T) bool, level int) {
	if len(slice) <= 1 {
		return
	}

	var mergeSort func(arr []T, compare func(a, b T) bool, depth int, wg *sync.WaitGroup) []T
	var merge func(left, right []T, compare func(a, b T) bool) []T

	mergeSort = func(arr []T, compare func(a, b T) bool, depth int, parentWG *sync.WaitGroup) []T {
		defer func() {
			if parentWG != nil {
				parentWG.Done()
			}
		}()

		if len(arr) <= 1 {
			return arr
		}

		mid := len(arr) / 2

		left := arr[:mid]
		right := arr[mid:]

		if depth < level {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				left = mergeSort(left, compare, depth+1, &wg)
			}()
			right = mergeSort(right, compare, depth+1, &wg)

			wg.Wait()
		} else {
			left = mergeSort(left, compare, depth+1, nil)
			right = mergeSort(right, compare, depth+1, nil)
		}

		return merge(left, right, compare)
	}

	merge = func(left, right []T, compare func(a, b T) bool) []T {
		result := make([]T, 0, len(left)+len(right))
		i, j := 0, 0

		for i < len(left) && j < len(right) {
			if compare(left[i], right[j]) {
				result = append(result, left[i])
				i++
			} else {
				result = append(result, right[j])
				j++
			}
		}

		// Append remaining elements
		result = append(result, left[i:]...)
		result = append(result, right[j:]...)
		return result
	}

	// Sort and copy back the result

	var wg sync.WaitGroup
	depth := 0

	wg.Add(1)
	sorted := mergeSort(slice, compare, depth+1, &wg)
	wg.Wait()

	copy(slice, sorted)
}

// Verify that the `slice` is sorted.
// The `compare` is a validation function, This return `true` only when every
// `compare(a, b) == true` for two adjacent `a` and `b` elements in `slice`.
func IsSorted[T any](slice []T, compare func(a, b T) bool) bool {
	for i := 0; i < len(slice)-1; i++ {
		if !compare(slice[i], slice[i+1]) {
			return false
		}
	}

	return true
}
