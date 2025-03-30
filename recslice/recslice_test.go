package recslice_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/elecbug/go-dspkg/recslice"
)

func TestRecslice(t *testing.T) {
	size := 100000
	check := 100

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	recs := recslice.New[int](100)
	recs.AutoCompact(ctx, 100*time.Millisecond)
	for i := 0; i < size; i++ {
		recs.Insert(i, i)
	}
	// t.Logf("Recslice: %v", recs.ToSlice())
	recs.Set(50, -1)
	t.Logf("Recslice[50]: %d", recs.Get(50))

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
}
