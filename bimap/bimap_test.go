package bimap_test

import (
	"go-type/bimap"
	"testing"
)

func TestBiMap(t *testing.T) {
	bimap := bimap.New[string, int]()

	bimap.Set("Apple", 5)
	bimap.Set("Banana", 4)
	bimap.Set("Orange", 3)

	if apple, ok := bimap.GetByKey("Apple"); ok {
		t.Logf(`bimap.GetByKey("Apple") is `+"%d", apple)
	} else {
		t.Fatal("Do not find {\"Apple\", ?}")
	}

	if _4, ok := bimap.GetByValue(4); ok {
		t.Logf(`bimap.GetByValue(4) is `+"%s", _4)
	} else {
		t.Fatal("Do not find {?, 4}")
	}
}
