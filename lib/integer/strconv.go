package integer

import (
	"strconv"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	// weave
	RegisterGoGate(new(GoFormatInt64))
	// eval
	RegisterEvalGate(new(GoFormatInt64))
}

type GoFormatInt64 struct {
	Int64 int64  `ko:"name=int64,monadic"`
	Base  *int64 `ko:"name=base"`
}

func (g *GoFormatInt64) Play(ctx *runtime.Context) string {
	return strconv.FormatInt(g.Int64, int(OptInt64(g.Base, 10)))
}
