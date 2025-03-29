package slice

import "github.com/elecbug/go-type/slice/compare_type"

// `Bsearch` performs a binary search on a sorted slice.
// `compare` should return true if the `target` is less than or equal to the current element.
// The function returns the first index i such that `compare(slice[i]) == compare_type.EQUAL`.
// If no such index exists, it returns `-1`.
func Bsearch[T any](slice []T, compare func(target T) int) int {
	left, right := 0, len(slice)-1

	for left <= right {
		mid := (left + right) / 2

		if compare(slice[mid]) == compare_type.EQUAL {
			return mid
		} else if compare(slice[mid]) == compare_type.TARGET_SMALL {
			left = mid + 1
		} else if compare(slice[mid]) == compare_type.TARGET_BIG {
			right = mid - 1
		} else {
			panic("The `compare` function must returns compare_type.XXX")
		}
	}

	return -1
}
