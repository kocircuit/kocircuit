package boolean

import (
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterGoGateAt("boolean", "And", &GoAnd{})
	RegisterGoGateAt("boolean", "Or", &GoOr{})
	RegisterGoGateAt("boolean", "Xor", &GoXor{})
}

type GoAnd struct {
	Series []bool `ko:"name=series,monadic"`
}

func (and *GoAnd) Play(ctx *runtime.Context) bool {
	r := true
	for _, x := range and.Series {
		r = r && x
	}
	return r
}

type GoOr struct {
	Series []bool `ko:"name=series,monadic"`
}

func (or *GoOr) Play(ctx *runtime.Context) bool {
	r := false
	for _, x := range or.Series {
		r = r || x
	}
	return r
}

type GoXor struct {
	Series []bool `ko:"name=series,monadic"`
}

func (xor *GoXor) Play(ctx *runtime.Context) bool {
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
