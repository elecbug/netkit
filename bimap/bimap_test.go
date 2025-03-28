package bimap_test

import (
	"testing"

	"github.com/elecbug/go-type/bimap"
)

func TestBimap(t *testing.T) {
	bimap := bimap.New[string, int]()

	bimap.Set("Apple", 5)
	bimap.Set("Banana", 4)
	bimap.Set("Orange", 3)

	if apple, ok := bimap.GetByKey("Apple"); ok {
		t.Logf(`bimap.GetByKey("Apple") is `+"%d", apple)
	} else {
		t.Logf("Do not find {\"Apple\", ?}")
	}

	if _4, ok := bimap.GetByValue(4); ok {
		t.Logf(`bimap.GetByValue(4) is `+"%s", _4)
	} else {
		t.Logf("Do not find {?, 4}")
	}

	bimap.Set("Cake", 4)

	if banana, ok := bimap.GetByKey("Banana"); ok {
		t.Logf(`bimap.GetByKey("Banana") is `+"%d", banana)
	} else {
		t.Logf("Do not find {\"Banana\", ?}")
	}

	if _4, ok := bimap.GetByValue(4); ok {
		t.Logf(`bimap.GetByValue(4) is `+"%s", _4)
	} else {
		t.Logf("Do not find {?, 4}")
	}

	if orange, ok := bimap.GetByKey("Orange"); ok {
		t.Logf(`bimap.GetByKey("Orange") is `+"%d", orange)
	} else {
		t.Logf("Do not find {\"Orange\", ?}")
	}

	bimap.DeleteByKey("Orange")

	if orange, ok := bimap.GetByKey("Orange"); ok {
		t.Logf(`bimap.GetByKey("Orange") is `+"%d", orange)
	} else {
		t.Logf("Do not find {\"Orange\", ?}")
	}

	t.Logf("bimap to list: %v", bimap.ToList())
}
