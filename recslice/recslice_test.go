package recslice_test

import (
	"fmt"
	"testing"

	"github.com/elecbug/go-dspkg/recslice"
)

func TestRecslice(t *testing.T) {
	size := 100

	a := recslice.New[int](10)

	for i := 0; i < size; i++ {
		a.Insert(a.Length()/2, i)
	}

	t.Logf("Slice: %v", a.Print(0, func(value int) string { return fmt.Sprintf("%d", value) }))

	t.Log("After deleting 3 items:")
	a.Delete(size / 2)
	a.Delete(size / 3)
	a.Delete(size / 4)

	t.Logf("Slice: %v", a.Print(0, func(value int) string { return fmt.Sprintf("%d", value) }))
	t.Logf("Get middle: %v", a.Get(size/2))
}
