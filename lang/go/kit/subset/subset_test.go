package subset

import (
	"testing"
)

type testSmall struct {
	String string
	Int    *int
	Float  *float64
}

type testLarge struct {
	String  string
	Int     *int
	Float   *float64
	Complex *complex128
}

var testSubset = []struct {
	Left, Right interface{}
	Expect      bool
}{
	{1, 2, false},
	{"Abc", "Abc", true},
	{[]int{1}, []int{1, 2}, true},
	{[]int{1}, []byte{1, 2}, false},
	{([]int)(nil), []int{}, true},
	{testSmall{"a", nil, nil}, testLarge{"a", nil, newFloat(3.14), nil}, true},
}

func newFloat(f float64) *float64 {
	return &f
}

func TestSubset(t *testing.T) {
	for i, test := range testSubset {
		got := IsSubset(test.Left, test.Right)
		if got != test.Expect {
			t.Errorf("test %d: exepecting %v, got %v", i, test.Expect, got)
		}
	}
}
