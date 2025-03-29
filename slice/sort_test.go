package slice_test

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/elecbug/go-type/slice"
)

func isAscended(slice []int) bool {
	for i := 0; i < len(slice)-1; i++ {
		if slice[i] > slice[i+1] {
			return false
		}
	}

	return true
}

func TestSort(t *testing.T) {
	result := []float64{0, 0, 0}
	iter := 10

	for i := 0; i < iter; i++ {
		data := make([]int, 1000000)

		for i := range data {
			data[i] = rand.Int()
		}

		copy1 := make([]int, len(data))
		copy2 := make([]int, len(data))
		copy3 := make([]int, len(data))
		copy(copy1, data)
		copy(copy2, data)
		copy(copy3, data)

		start := time.Now().UnixNano()
		sort.Slice(copy1, func(i, j int) bool { return copy1[i] < copy1[j] })
		end := float64(time.Now().UnixNano()-start) / 10e6
		t.Logf("General sort: %v, time: %fms", isAscended(copy1), end)
		result[0] += end

		start = time.Now().UnixNano()
		slice.Sort(copy2, func(i, j int) bool { return i < j })
		end = float64(time.Now().UnixNano()-start) / 10e6
		t.Logf("go_type sort: %v, time: %fms", isAscended(copy2), end)
		result[1] += end

		start = time.Now().UnixNano()
		slice.ParallelSort(copy3, func(i, j int) bool { return i < j }, 7)
		end = float64(time.Now().UnixNano()-start) / 10e6
		t.Logf("go_type parallel sort: %v, time: %fms", isAscended(copy3), end)
		result[2] += end
	}

	t.Logf("--Benchmark--")
	t.Logf("General sort: %fms", result[0]/float64(iter))
	t.Logf("go_type sort: %fms", result[1]/float64(iter))
	t.Logf("go_type parallel sort: %fms", result[2]/float64(iter))
}
