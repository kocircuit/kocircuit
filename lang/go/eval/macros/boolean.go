package macros

import (
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGateAt("", "And", new(EvalAnd))
	RegisterEvalGateAt("", "Or", new(EvalOr))
	RegisterEvalGateAt("", "Xor", new(EvalXor))
	RegisterEvalGateAt("", "Not", new(EvalNot))
}

type EvalAnd struct {
	Series []bool `ko:"name=series,monadic"`
}

func (and *EvalAnd) Play(ctx *runtime.Context) bool {
	r := true
	for _, x := range and.Series {
		r = r && x
	}
	return r
}

type EvalOr struct {
	Series []bool `ko:"name=series,monadic"`
}

func (or *EvalOr) Play(ctx *runtime.Context) bool {
	r := false
	for _, x := range or.Series {
		r = r || x
	}
	return r
}

type EvalXor struct {
	Series []bool `ko:"name=series,monadic"`
}

func (xor *EvalXor) Play(ctx *runtime.Context) bool {
	var q uint
	for _, x := range xor.Series {
		q ^= boolUint(x)
	}
	return q != 0
}

func boolUint(b bool) uint {
	if b {
		return 1
	} else {
		return 0
	}
}

type EvalNot struct {
	Bool bool `ko:"name=bool,monadic"`
}

func (not *EvalNot) Play(ctx *runtime.Context) bool {
	return !not.Bool
}
