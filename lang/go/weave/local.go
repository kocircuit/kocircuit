package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type GoLocal struct {
	Origin *Span `ko:"name=origin"`
	// Step identifies the step returning the un-shaped type.
	// GoLocal resulting from a boundary.Figure (constant) has nil step.
	Step *GoStep `ko:"name=step"`
	Type GoType  `ko:"name=type"`
	// Shaper represents a logical generalized subselection over the root type.
	Shaper Shaper       `ko:"name=shaper"`
	Cached *AssignCache `ko:"name=cached"`
}

// step can be nil, if typ is GoNumber or GoVariety.
func NewGoLocal(
	origin *Span,
	step *GoStep,
	typ GoType,
) *GoLocal {
	return &GoLocal{
		Origin: origin,
		Step:   step,
		Type:   typ,
		Shaper: IdentityShaper(origin, typ),
	}
}

func (local *GoLocal) Inherit(
	origin *Span,
	step *GoStep,
	typ GoType,
) *GoLocal {
	return &GoLocal{
		Origin: origin,
		Step:   step,
		Type:   typ,
		Shaper: IdentityShaper(origin, typ),
		Cached: AssignCacheUnion(local.Cached, step.Cached),
	}
}

func (local *GoLocal) Image() GoType {
	return local.Shaper.Shape(local.Type)
}

func (local *GoLocal) String() string {
	return Sprint(local)
}

func (local *GoLocal) Extend(span *Span, shaper ...Shaper) *GoLocal {
	return &GoLocal{
		Origin: span,
		Step:   local.Step,
		Type:   local.Type,
		Shaper: CompressShapers(span, local.Shaper, CompressShapers(span, shaper...)),
		Cached: local.Cached,
	}
}
