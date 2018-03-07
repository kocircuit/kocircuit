package model

import (
	"fmt"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type GoSlotForm interface {
	FormExpr(...*GoSlotExpr) GoExpr
}

type GoSlotExpr struct {
	Slot Slot   `ko:"name=slot"`
	Expr GoExpr `ko:"name=expr"`
}

type GoSlotFormExpr struct {
	SlotExpr []*GoSlotExpr `ko:"name=slot_expr"`
	Form     GoSlotForm    `ko:"name=form"`
}

func (expr *GoSlotFormExpr) RenderExpr(fileCtx GoFileContext) string {
	return expr.Form.FormExpr(expr.SlotExpr...).RenderExpr(fileCtx)
}

func FindSlotExpr(arg []*GoSlotExpr, find Slot) GoExpr {
	for _, arg := range arg {
		if arg.Slot == find {
			return arg.Expr
		}
	}
	return nil
}

type GoUnknownForm struct{}

func (*GoUnknownForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return BrokenExpr
}

type GoInvariantForm struct {
	Expr GoExpr `ko:"name=expr"`
}

func (form *GoInvariantForm) FormExpr(...*GoSlotExpr) GoExpr { return form.Expr }

type GoIdentityForm struct{}

func (*GoIdentityForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return FindSlotExpr(arg, RootSlot{})
}

type GoNegativeForm struct {
	Shaper Shaper `ko:"name=shaper"`
}

func (form *GoNegativeForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoNegativeExpr{
		&GoShapeExpr{
			Shaper: form.Shaper,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		},
	}
}

type GoShapeForm struct {
	Shaper Shaper `ko:"name=shaper"`
}

func (form *GoShapeForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoShapeExpr{
		Shaper: form.Shaper,
		Expr:   FindSlotExpr(arg, RootSlot{}),
	}
}

// GoVarietyFitForm represents the form:
// 	&Variety{
// 		BranchN: RootSlot,
// 	}
type GoVarietyFitForm struct {
	Varietal GoVarietal `ko:"name=varietal"`
	Branch   *GoBranch  `ko:"name=branch"`
	Augments GoSlotForm `ko:"name=augments"`
}

func (fit *GoVarietyFitForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoMakeStructExpr{
		For: fit.Varietal.Real(),
		Field: []*GoFieldExpr{{
			Field: fit.Branch.BranchRealField(),
			Expr: &GoMakeStructExpr{
				For: fit.Branch.BranchRealField().Type,
				Field: []*GoFieldExpr{{
					Field: fit.Branch.AugmentsField(),
					Expr: fit.Augments.FormExpr(
						&GoSlotExpr{
							Slot: RootSlot{},
							Expr: FindSlotExpr(arg, RootSlot{}),
						},
					),
				}},
			},
		}},
	}
}

// GoExtendForm wraps a form into a function.
type GoExtendForm struct {
	Origin  *Span      `ko:"name=origin"`
	Prefix  string     `ko:"name=prefix"`
	Form    GoSlotForm `ko:"name=form"`
	Arg     GoType     `ko:"name=arg"`
	Returns GoType     `ko:"name=returns"`
}

func (ff *GoExtendForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: ff.funcExpr(),
		Arg: []GoExpr{
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}

func (ff *GoExtendForm) funcExpr() GoFuncExpr {
	return ff.wrapFuncExpr(
		&GoSlotFormExpr{
			SlotExpr: []*GoSlotExpr{{Slot: RootSlot{}, Expr: ff.argExpr()}},
			Form:     ff.Form,
		},
	)
}

func (ff *GoExtendForm) CircuitEffect() *GoCircuitEffect {
	return &GoCircuitEffect{
		DuctFunc: []GoFuncExpr{ff.funcExpr()},
	}
}

func (ff *GoExtendForm) funcName() *GoNameExpr {
	return &GoNameExpr{
		Origin: ff.Origin,
		Name:   fmt.Sprintf("%s_%s", ff.Prefix, ff.Origin.SpanID().String()),
	}
}

func (ff *GoExtendForm) argExpr() GoExpr {
	return &GoVerbatimExpr{"arg"}
}

func (ff *GoExtendForm) wrapFuncExpr(line ...GoExpr) GoFuncExpr {
	return &GoShaperFuncExpr{
		Comment:    fmt.Sprintf("(%s)", ff.Origin.SourceLine()),
		FuncName:   ff.funcName(),
		ArgName:    ff.argExpr(),
		ArgType:    ff.Arg,
		ReturnType: ff.Returns,
		Line:       line,
	}
}
