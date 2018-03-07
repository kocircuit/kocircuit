package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type InvokingCtx struct {
	Origin         *Span               `ko:"name=origin"`
	Varietal       GoVarietal          `ko:"name=varietal"`
	MacroTransform []*GoMacroTransform `ko:"name=macro_transform"` // macro transforms corresponding to each projection line
	VtyReturns     GoType              `ko:"name=vty_returns"`
}

type GoMacroTransform struct {
	Origin      *Span      `ko:"name=origin"` // invocation site of varietal, whose branch results in this transform
	Macro       Macro      `ko:"name=macro"`
	Arg         GoType     `ko:"name=arg"`         // type of argument expected by this expression transform
	Returns     GoType     `ko:"name=returns"`     // type returned from invoking the branch projection macro
	SlotForm    GoSlotForm `ko:"name=slot_form"`   // form an arg-type expression into a returns-type expression
	ExpandValve *GoValve   `ko:"name=expandValve"` // used by GoExpandMacro to pass a valve to GoFixMacro macro
}

func NewInvokingCtx(
	span *Span,
	varietal GoVarietal,
	macroTransform []*GoMacroTransform,
	vtyReturns GoType,
) *InvokingCtx {
	return &InvokingCtx{
		Origin:         span,
		Varietal:       varietal,
		MacroTransform: macroTransform,
		VtyReturns:     vtyReturns,
	}
}

// DuctFuncExpr returns an expression for the invocation function of a variety:
//
//	func invoke_VtyID(step_ctx *runtime.Context, proj *VtyProj) *VtyCallReturns {
//		switch {
//		case proj.Projecion_1_0_7_join != nil:
//			return shapeReturn(
//				shapeArg(
//					proj.Projecion_1_0_7_join,
//				).Play(step_ctx),
//			)
//		█
//		default:
//			panic("o")
//		}
//	}
func (ctx *InvokingCtx) FuncExprEffect() (*GoDuctFuncExpr, *GoCircuitEffect) {
	switchCaseExpr, effect := ctx.SwitchCaseExpr()
	return &GoDuctFuncExpr{
		Comment:    fmt.Sprintf("(%s)", ctx.Origin.SourceLine()),
		FuncName:   ctx.FuncName(),
		ArgName:    &GoVerbatimExpr{ctx.ProjArgName()},
		ArgType:    VarietalProjectionReal(ctx.Varietal),
		ReturnType: ctx.VtyReturns,
		Line: []GoExpr{
			&GoSwitchExpr{
				Over:    EmptyExpr,
				Case:    switchCaseExpr,
				Default: &GoPanicExpr{},
			},
		},
	}, effect
}

func (ctx *InvokingCtx) FuncExpr() *GoDuctFuncExpr {
	funcExpr, _ := ctx.FuncExprEffect()
	return funcExpr
}

func (ctx *InvokingCtx) CircuitEffect() *GoCircuitEffect {
	funcExpr, effect := ctx.FuncExprEffect()
	return effect.Aggregate(
		&GoCircuitEffect{
			DuctType: []GoType{VarietalProjectionReal(ctx.Varietal), ctx.VtyReturns},
			DuctFunc: []GoFuncExpr{funcExpr},
		},
	)
}

func (ctx *InvokingCtx) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: ctx.Origin,
		Name:   fmt.Sprintf("invoke_%s", ctx.Varietal.TypeID()),
	}
}

func (ctx *InvokingCtx) ProjArgName() string { return "proj" }

func (ctx *InvokingCtx) VtyProjExpr() GoExpr { return &GoVerbatimExpr{ctx.ProjArgName()} }

func (ctx *InvokingCtx) SwitchCaseExpr() ([]*GoSwitchCaseExpr, *GoCircuitEffect) {
	line := VarietalProjectionLine(ctx.Varietal) // line is one-to-one index-wise with ctx.MacroTransform
	switchCase, effect := make([]*GoSwitchCaseExpr, len(line)), &GoCircuitEffect{}
	for i, line := range line {
		lineField := line.ProjectionRealField()
		lineFieldSelectExpr := &GoSelectExpr{ // proj.Line_XYZ
			Into:  ctx.VtyProjExpr(),
			Field: lineField.Name,
		}
		macroExpr, macroEffect := ctx.MacroTransform[i].Form(
			ctx.Origin,
			lineFieldSelectExpr,
			lineField.Type,
			ctx.VtyReturns,
		)
		effect = effect.Aggregate(macroEffect)
		switchCase[i] = &GoSwitchCaseExpr{
			Predicate: &GoInequalityExpr{ // case proj.Line_XYZ != nil:
				Left:  lineFieldSelectExpr,
				Right: NilExpr,
			},
			Expr: &GoBlockExpr{ // return MacroExpr
				Line: []GoExpr{&GoReturnExpr{macroExpr}},
			},
		}
	}
	return switchCase, effect
}

// 	shapeReturn(
// 		shapeArg(
// 			proj.Projecion_1_0_7_join,
// 		).Play(step_ctx),
// 	)
//	shapeReturn(
//		Transform(
//			shapeArg(
//				proj.Projecion_1_0_7_join,
//			),
//		),
//	)
func (macroTransform *GoMacroTransform) Form(span *Span, argExpr GoExpr, arg GoType, returns GoType) (
	GoExpr,
	*GoCircuitEffect,
) {
	argShaper := macroTransform.ArgShaper(span, arg)
	returnShaper := macroTransform.ReturnShaper(span, returns)
	return &GoShapeExpr{
			Shaper: returnShaper,
			Expr: macroTransform.SlotForm.FormExpr(
				&GoSlotExpr{
					Slot: RootSlot{},
					Expr: &GoShapeExpr{
						Shaper: argShaper,
						Expr:   argExpr,
					},
				},
			),
		}, argShaper.CircuitEffect().
			Aggregate(returnShaper.CircuitEffect()).
			AggregateDuctType(arg, returns)
}

func (macroTransform *GoMacroTransform) ArgShaper(span *Span, arg GoType) Shaper {
	shaper, _, err := Assign(span, arg, macroTransform.Arg)
	if err != nil {
		panic(err)
	}
	// MustAssign(shaper.Shape(arg), macroTransform.Arg)
	return shaper
}

func (macroTransform *GoMacroTransform) ReturnShaper(span *Span, returns GoType) Shaper {
	shaper, _, err := Assign(span, macroTransform.Returns, returns)
	if err != nil {
		panic(err)
	}
	// MustAssign(shaper.Shape(macroTransform.Returns), returns)
	return shaper
}
