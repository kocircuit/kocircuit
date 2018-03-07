package model

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type UnwrapSingletonShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *UnwrapSingletonShaper) ShaperID() string {
	return Mix("UnwrapSingletonShaper", sh.ShapingID())
}

func (sh *UnwrapSingletonShaper) String() string { return Sprint(sh) }

func (sh *UnwrapSingletonShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	expr := &GoIndexExpr{Container: ofExpr, Index: ZeroExpr}
	return expr.RenderExpr(fileCtx)
}

func (sh *UnwrapSingletonShaper) Reverse() Shaper {
	return &WrapSingletonShaper{Shaping: sh.Flip()}
}

type WrapSingletonShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *WrapSingletonShaper) ShaperID() string {
	return Mix("WrapSingletonShaper", sh.ShapingID())
}

func (sh *WrapSingletonShaper) String() string { return Sprint(sh) }

func (sh *WrapSingletonShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	expr := &GoMakeSequenceExpr{Type: sh.To, Elem: []GoExpr{ofExpr}}
	return expr.RenderExpr(fileCtx)
}

func (sh *WrapSingletonShaper) Reverse() Shaper {
	return &UnwrapSingletonShaper{Shaping: sh.Flip()}
}
