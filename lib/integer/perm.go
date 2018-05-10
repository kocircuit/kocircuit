package integer

import (
	"math/rand"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoPermuteInt64))
}

type GoPermuteInt64 struct {
	Len  int64 `ko:"name=len"`
	Seed int64 `ko:"name=seed"`
}

func (g *GoPermuteInt64) Play(ctx *runtime.Context) []int64 {
	r := make([]int64, g.Len)
	for i, index := range rand.New(rand.NewSource(g.Seed)).Perm(int(g.Len)) {
		r[i] = int64(index)
	}
	return r
}
