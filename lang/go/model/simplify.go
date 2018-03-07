package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func Simplify(span *Span, complex GoType) (simple GoType, simplifier Shaper) {
	ctx := NewSimplifyCtx(span, complex)
	for {
		ctx = ctx.Reduce(ReduceNeverNilPtr)   // reduce GoNeverNilPtr{█} to █
		ctx = ctx.Reduce(ReduceAliasAfterPtr) // reduce GoPtr{GoAlias{█}} to GoPtr{█}
		ctx = ctx.Reduce(ReduceAlias)         // reduce GoAlias{█} to █
		ctx = ctx.Reduce(ReducePtrPtr)        // reduce GoPtr{...GoPtr{█}} to GoPtr{█}
		ctx = ctx.Reduce(ReducePtrSlice)      // reduce Go[NeverNil]Ptr{GoSlice{█}} to GoSlice{█}
		ctx = ctx.Reduce(ReducePtrEmpty)      // reduce Go[NeverNil]Ptr{GoEmpty} to GoEmpty
		ctx = ctx.Reduce(ReduceSliceEmpty)    // reduce GoSlice{GoEmpty} to GoEmpty
		ctx = ctx.Reduce(ReduceStructEmpty)   // reduce GoStruct{} to GoEmpty
		ctx = ctx.Reduce(ReduceArrayEmpty)    // reduce GoArray{GoEmpty} to GoEmpty
		ctx = ctx.Reduce(ReduceSingleton)     // reduce GoArray{1, █} to █
		if !ctx.Changed {                     // greedily
			break
		} else {
			ctx = ctx.Clear()
		}
	}
	return ctx.Type, ctx.Simplifier
}

type SimplifyCtx struct {
	Origin *Span `ko:"name=origin"`
	//
	Type       GoType `ko:"name=type"`
	Simplifier Shaper `ko:"name=simplifier"`
	Changed    bool   `ko:"name=changed"`
}

func NewSimplifyCtx(span *Span, typ GoType) *SimplifyCtx {
	return &SimplifyCtx{Origin: span, Type: typ, Simplifier: IdentityShaper(span, typ)}
}

type GoTypeReducer func(*Span, GoType) (GoType, Shaper)

func (ctx *SimplifyCtx) Reduce(reducer GoTypeReducer) *SimplifyCtx {
	rt, d := reducer(ctx.Origin, ctx.Type)
	return &SimplifyCtx{
		Origin:     ctx.Origin,
		Type:       rt,
		Simplifier: CompressShapers(ctx.Origin, ctx.Simplifier, d),
		Changed:    ctx.Changed || !IsIdentityShaper(d),
	}
}

func (ctx *SimplifyCtx) Clear() *SimplifyCtx {
	return &SimplifyCtx{Origin: ctx.Origin, Type: ctx.Type, Simplifier: ctx.Simplifier, Changed: false}
}
