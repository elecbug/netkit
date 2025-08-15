// Package slice provides generic slice utilities such as binary search and sorting.
package slice

// CompareType indicates the relation of the probed element to the target.
type CompareType int

const (
	// TARGET_SMALL indicates the target is larger than the probed element; search right.
	TARGET_SMALL CompareType = -1
	// EQUAL indicates the target matches the probed element.
	EQUAL CompareType = 0
	// TARGET_BIG indicates the target is smaller than the probed element; search left.
	TARGET_BIG CompareType = 1
)

// Bsearch performs a binary search on a sorted slice.
// The compare function should return one of compare_type.TARGET_SMALL,
// compare_type.EQUAL, or compare_type.TARGET_BIG for the probed element.
// It returns the index i such that compare(slice[i]) == CompareType.EQUAL,
// or -1 if no such element exists.
func Bsearch[T any](slice []T, compare func(target T) CompareType) int {
	left, right := 0, len(slice)-1

	for left <= right {
		mid := (left + right) / 2

		if compare(slice[mid]) == EQUAL {
			return mid
		} else if compare(slice[mid]) == TARGET_SMALL {
			left = mid + 1
		} else if compare(slice[mid]) == TARGET_BIG {
			right = mid - 1
		} else {
			panic("compare must return a value from compare_type")
		}
	}

	return -1
}
