package kahnsort

import (
	"fmt"
	"math"
	"testing"
)

// XXX: TestKahnSortFibonacci is flaky due to sorting non-determinism introduced by ranging over maps.
func TestKahnSortFibonacci(t *testing.T) {
	f := fibonacciNode{5, 1, 1}
	sorted := Sort([]Node{f})
	for i, g := range sorted {
		fmt.Printf("%v\n", g)
		if f != g {
			t.Errorf("node %d: expecting %v, got %v", i, f, g)
		}
		f, _ = f.next()
	}
}

type fibonacciNode struct {
	LogLimit float64
	M, N     float64
}

func (f fibonacciNode) next() (fibonacciNode, bool) {
	if math.Log(f.M) > f.LogLimit {
		return fibonacciNode{}, false
	}
	return fibonacciNode{LogLimit: f.LogLimit, M: f.N, N: f.M + f.N}, true
}

func (f fibonacciNode) Down() []Node {
	g, ok := f.next()
	if !ok {
		return nil
	}
	h, ok := g.next()
	if !ok {
		return nil
	}
	return []Node{g, h}
}
