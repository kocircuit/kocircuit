package tree

import (
	"testing"
)

func TestMultiSet(t *testing.T) {
	mset := MultiSet{}
	if mset.Count(1) != 0 {
		t.Errorf("count 0")
	}
	mset = mset.Add(1)
	if mset.Count(1) != 1 {
		t.Errorf("count 1")
	}
	mset = mset.Add(1)
	if mset.Count(1) != 2 {
		t.Errorf("count 2")
	}
}
