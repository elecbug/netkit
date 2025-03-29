package slice_test

import (
	"testing"

	"github.com/elecbug/go-type/slice"
	"github.com/elecbug/go-type/slice/compare_type"
)

func TestBsearch(t *testing.T) {
	arr := make([]int, 1000000)

	for i := 0; i < len(arr); i++ {
		arr[i] = i * 2
	}

	idx := slice.Bsearch(arr, func(target int) int {
		if target == 123456 {
			return compare_type.EQUAL
		} else if target < 123456 {
			return compare_type.TARGET_SMALL
		} else {
			return compare_type.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) int {
		if target == 123457 {
			return compare_type.EQUAL
		} else if target < 123457 {
			return compare_type.TARGET_SMALL
		} else {
			return compare_type.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) int {
		if target == -1 {
			return compare_type.EQUAL
		} else if target < -1 {
			return compare_type.TARGET_SMALL
		} else {
			return compare_type.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) int {
		if target == 10000000 {
			return compare_type.EQUAL
		} else if target < 10000000 {
			return compare_type.TARGET_SMALL
		} else {
			return compare_type.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}
}
