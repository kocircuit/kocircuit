package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// GoEmpty is a unique type indicating non-existing fields (not passed arguments or selection into missing fields).
// (1) Everything assigns to GoEmpty (in assignToEmpty)
// (2) GoEmpty assigns only to GoEmpty (in assignMatrix), GoPtr (in assignCross) and GoSlice (in assignCross).
// (3) GoEmpty generalizes (unifies) with GoEmpty (in Generalize), GoPtr (in Generalize) and GoSlice (in Generalize),
// and with the remaining types by making them optional.
type GoEmpty struct {
	ID     string `ko:"name=id"`
	Origin *Span  `ko:"name=origin"`
}

func NewGoEmpty(span *Span) *GoEmpty {
	return &GoEmpty{ID: Mix("empty"), Origin: span}
}

func (empty *GoEmpty) StructureField() []*GoField { return nil }

func (empty *GoEmpty) TypeID() string { return empty.ID }

func (empty *GoEmpty) Doc() string { return Sprint(empty) }

func (empty *GoEmpty) String() string { return Sprint(empty) }

func (empty *GoEmpty) Sketch(ctx *GoSketchCtx) interface{} {
	return "()"
}

func (empty *GoEmpty) RenderRef(GoFileContext) string { return "struct{/*empty*/}" }

func (empty *GoEmpty) RenderDef(GoFileContext) string { return "struct{/*empty*/}" }

func (empty *GoEmpty) RenderZero(GoFileContext) string { return "struct{/*empty*/}{}" }

func (empty *GoEmpty) Tag() []*GoTag { return nil }

// EraseShaper converts any type to empty (erasing all information).
type EraseShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *EraseShaper) ShaperID() string { return Mix("EraseShaper", sh.ShapingID()) }

func (sh *EraseShaper) String() string { return Sprint(sh) }

func (sh *EraseShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	return sh.Shaping.To.(*GoEmpty).RenderZero(fileCtx)
}

func (sh *EraseShaper) Reverse() Shaper {
	return &UneraseShaper{Shaping: sh.Flip()}
}

// UneraseShaper is the inverse of EraseShaper.
type UneraseShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *UneraseShaper) ShaperID() string { return Mix("UneraseShaper", sh.ShapingID()) }

func (sh *UneraseShaper) String() string { return Sprint(sh) }

func (sh *UneraseShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	switch u := sh.To.(type) {
	case *GoPtr, *GoNeverNilPtr, *GoSlice:
		return fmt.Sprintf("(%s)(nil)", sh.To.RenderRef(fileCtx))
	case *GoStruct:
		if u.Len() == 0 {
			return u.RenderZero(fileCtx)
		}
	}
	panic("o")
}

func (sh *UneraseShaper) Reverse() Shaper {
	return &EraseShaper{Shaping: sh.Flip()}
}

// IrreversibleEraseShaper converts any type to empty (erasing all information).
type IrreversibleEraseShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *IrreversibleEraseShaper) ShaperID() string {
	return Mix("IrreversibleEraseShaper", sh.ShapingID())
}

func (sh *IrreversibleEraseShaper) String() string { return Sprint(sh) }

func (sh *IrreversibleEraseShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	return sh.Shaping.To.(*GoEmpty).RenderZero(fileCtx)
}

func (sh *IrreversibleEraseShaper) Reverse() Shaper { panic("irreversible") }
