package slice_test

import (
	"testing"

	"github.com/elecbug/netkit/slice"
)

func TestBsearch(t *testing.T) {
	arr := make([]int, 1000000)

	for i := 0; i < len(arr); i++ {
		arr[i] = i * 2
	}

	idx := slice.Bsearch(arr, func(target int) slice.CompareType {
		if target == 123456 {
			return slice.EQUAL
		} else if target < 123456 {
			return slice.TARGET_SMALL
		} else {
			return slice.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) slice.CompareType {
		if target == 123457 {
			return slice.EQUAL
		} else if target < 123457 {
			return slice.TARGET_SMALL
		} else {
			return slice.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) slice.CompareType {
		if target == -1 {
			return slice.EQUAL
		} else if target < -1 {
			return slice.TARGET_SMALL
		} else {
			return slice.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}

	idx = slice.Bsearch(arr, func(target int) slice.CompareType {
		if target == 10000000 {
			return slice.EQUAL
		} else if target < 10000000 {
			return slice.TARGET_SMALL
		} else {
			return slice.TARGET_BIG
		}
	})

	if idx != -1 {
		t.Logf("Index: %d, value: %d", idx, arr[idx])
	} else {
		t.Logf("Do not find: %d", idx)
	}
}
