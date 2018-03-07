package tree

import (
	"fmt"
	"testing"
)

func TestPrintStruct(t *testing.T) {
	fmt.Println(Sprint(_two{_one{}, map[string]int{"a": 1, "b": 2}}))
}

type _one struct{}
type _two struct {
	X interface{} `ko:"name=x"`
	Y interface{} `ko:"name=y"`
}

func TestPrintCycle(t *testing.T) {
	r := &node{S: &node{}}
	r.T = r.S
	r.S.S = r
	r.S.T = r
	fmt.Println(Sprint(r))
}

type node struct {
	S *node
	T interface{}
}
