package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type GoCtx interface {
	AssignCache() *AssignCache
	Valve() *GoValve
	Label() string
}

func NearestGoCtx(span *Span) GoCtx {
	if span == nil {
		return nil
	} else if ctx, _ := span.Sheath.(GoCtx); ctx != nil {
		return ctx
	} else {
		return NearestGoCtx(span.Parent)
	}
}

func NewCacheSpan(cache *AssignCache) *Span {
	return RefineAssignCache(NewSpan(), cache) // AssignCache is a GoCtx and a Sheath.
}

func SpanCache(span *Span) *AssignCache {
	if ctx := NearestGoCtx(span); ctx != nil {
		return ctx.AssignCache()
	} else {
		return nil
	}
}

func SpanGroupPath(span *Span) GoGroupPath {
	return GoGroupPath{
		Group: KoPkgGroup,
		Path:  ChamberPath(span).Slash(),
	}
}
