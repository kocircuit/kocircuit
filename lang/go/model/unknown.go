package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Unknown interface {
	GoVarietal
	Sequence
	IsUnknown()
}

type GoUnknown struct {
	ID     string `ko:"name=id"`
	Origin *Span  `ko:"name=origin"`
}

func NewGoUnknown(span *Span) *GoUnknown {
	return &GoUnknown{ID: Mix("unknown", span.SpanID().String()), Origin: span}
}

func (u *GoUnknown) TypeID() string { return u.ID }

func (u *GoUnknown) IsUnknown() {}

func (u *GoUnknown) Doc() string { return Sprint(u) }

func (u *GoUnknown) String() string { return Sprint(u) }

func (u *GoUnknown) Sketch(ctx *GoSketchCtx) interface{} {
	return K{"!unknown": u.Origin.SourceLine()}
}

func (u *GoUnknown) Tag() []*GoTag { return nil }

func (u *GoUnknown) RenderRef(GoFileContext) string { panic("o") } // prohibit rendering

func (u *GoUnknown) RenderDef(GoFileContext) string { panic("o") } // prohibit rendering

func (u *GoUnknown) RenderZero(GoFileContext) string { panic("o") } // prohibit rendering

func (u *GoUnknown) Real() GoType {
	return NewGoAlias(u.UnknownAddress(), NewGoStruct())
}

func (u *GoUnknown) SequenceElem() GoType { return u }

func (u *GoUnknown) VarietyMacro() Macro { return GoUnknownMacro{} }

func (u *GoUnknown) VarietyBranch() []*GoBranch { return nil }

func (u *GoUnknown) VarietyExpr() GoExpr { return GoUnknownExpr{} }

func (u *GoUnknown) VarietyAddress() *GoAddress { return u.UnknownAddress() }

func (u *GoUnknown) VarietyProjectionAddress() *GoAddress { return u.UnknownAddress() }

func (u *GoUnknown) UnknownAddress() *GoAddress {
	return &GoAddress{
		Comment:   "unknown address",
		Span:      u.Origin,
		GroupPath: GoGroupPath{Group: KoPkgGroup, Path: "unknown"},
		Name:      "UnknownAddressXXX",
	}
}

func (u *GoUnknown) Select(*Span, Path) (Shape, Effect, error) {
	panic("o")
}

func (u *GoUnknown) Invoke(*Span) (Shape, Effect, error) {
	panic("o")
}

func (u *GoUnknown) Augment(span *Span, knot Knot) (Shape, Effect, error) {
	panic("o")
}

type UnknownShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *UnknownShaper) RenderExprShaping(GoFileContext, GoExpr) string { panic("o") }

func (sh *UnknownShaper) ShaperID() string { return Mix("UnknownShaper", sh.ShapingID()) }

func (sh *UnknownShaper) String() string { return Sprint(sh) }

type GoUnknownLogic struct{}

func (logic *GoUnknownLogic) Render(_ GoFileContext, _ []*GoSlotExpr) string { panic("o") }

type GoUnknownSlotForm struct{}

func (*GoUnknownSlotForm) FormExpr(...*GoSlotExpr) GoExpr { panic("o") }

type GoUnknownMacro struct{}

func (m GoUnknownMacro) MacroID() string { return m.Help() }

func (GoUnknownMacro) Label() string { return "unknown" }

func (GoUnknownMacro) MacroSheathString() *string { return PtrString("Unknown") }

func (GoUnknownMacro) Help() string { return "Unknown" }

func (GoUnknownMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return NewGoUnknown(span), nil, nil
}

type GoUnknownExpr struct{}

func (GoUnknownExpr) RenderExpr(GoFileContext) string {
	panic("o") // prohibit rendering
}
