package recslice_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/elecbug/go-dspkg/recslice"
)

func TestRecslice(t *testing.T) {
	size := 100000
	check := 100

	recs := recslice.New[int](100)
	for i := 0; i < size; i++ {
		recs.Insert(recs.Length(), i)
	}
	t.Logf("Recslice: %v", recs.ToSlice())

	arr := make([]int, size)
	for i := 0; i < size; i++ {
		arr[i] = i
	}

	start := time.Now().UnixNano()
	for i := 0; i < check; i++ {
		recs.Insert(rand.Intn(recs.Length()), 0)
		recs.Delete(rand.Intn(recs.Length()))
	}
	end := float64(time.Now().UnixNano()-start) / 10e6
	t.Logf("Recslice test: %f", end)

	start = time.Now().UnixNano()
	for i := 0; i < check; i++ {
		idx := rand.Intn(len(arr))
		arr = append(arr, 0)
		copy(arr[idx+1:], arr[idx:])
		arr[idx] = 0

		idx = rand.Intn(len(arr))
		copy(arr[idx:], arr[idx+1:])
		arr = arr[:len(arr)-1]
	}
	end = float64(time.Now().UnixNano()-start) / 10e6
	t.Logf("General array test: %f", end)

	// // t.Logf("Recslice: %v", a.Print(0, func(value int) string { return fmt.Sprintf("%d", value) }))
	// t.Logf("Recslice: %v", recs.ToSlice())
	// t.Logf("General array: %v", arr)

	// t.Log("After deleting 3 items:")
	// recs.Delete(size / 2)
	// recs.Delete(size / 3)
	// recs.Delete(size / 4)

	// // t.Logf("Recslice: %v", a.Print(0, func(value int) string { return fmt.Sprintf("%d", value) }))
	// t.Logf("Recslice: %v", recs.ToSlice())

	// t.Logf("Get middle: %v", recs.Get(size/2))
}
