package fast

import (
	"math/rand"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func BenchmarkInsert(b *testing.B) {
	var node *Node
	perm := rand.Perm(1000)
	for i := 0; i < 1000; i++ {
		node = Insert(node, BasicInt64Symbol(int64(perm[i])), XXX)
	}
}
