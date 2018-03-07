package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func RefineWeaveCtx(span *Span, ctx *GoWeaveCtx) *Span {
	return span.Refine(ctx)
}

func NearestWeaveCtx(span *Span) *GoWeaveCtx {
	if span == nil {
		return nil
	} else if ctx, _ := span.Sheath.(*GoWeaveCtx); ctx != nil {
		return ctx
	} else {
		return NearestWeaveCtx(span.Parent)
	}
}

func SpanWeave(span *Span, f *Func, arg GoStructure) (GoType, *GoCombineEffect, error) {
	return NearestWeaveCtx(span).Weave(span, f, arg)
}

func SpanClearWeavingCtx(span *Span) *Span {
	return RefineWeaveCtx(span, NearestWeaveCtx(span).Clear())
}

func AugmentSpanCache(span *Span, cache *AssignCache) *Span {
	ctx := NearestWeaveCtx(span)
	ctx2 := ctx.UseCache(
		AssignCacheUnion(ctx.AssignCache(), cache),
	)
	return RefineWeaveCtx(span, ctx2)
}
