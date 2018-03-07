package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Cut", new(GoCutMacro))
}

type GoCutMacro struct{}

func (m GoCutMacro) MacroID() string { return m.Help() }

func (GoCutMacro) Label() string { return "cut" }

func (m GoCutMacro) MacroSheathString() *string { return PtrString("Cut") }

func (m GoCutMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoCutMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveCut(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving cut")
	}
	if soln.IntoIsUnknown() {
		return GoUnknownMacro{}.Invoke(span, arg)
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type CutSolnArg struct {
	Arg        GoStructure   `ko:"name=arg"`
	Into       TypeExtractor `ko:"name=into"`
	IntoLifted GoType        `ko:"name=intoLifted"`
	IntoLift   Shaper        `ko:"name=intoLift"`
	Otherwise  TypeExtractor `ko:"name=otherwise"`
}

func (arg *CutSolnArg) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		arg.Into.Extractor.CircuitEffect(),
		arg.Otherwise.Extractor.CircuitEffect(),
		arg.IntoLift.CircuitEffect(),
	)
}

type CutSoln struct {
	Origin           *Span           `ko:"name=origin"`
	Arg              *CutSolnArg     `ko:"name=arg"`
	Returns          GoType          `ko:"name=returns"`
	IntoElemUnifier  Shaper          `ko:"name=intoElemUnifier"`
	OtherwiseUnifier Shaper          `ko:"name=otherwiseUnifier"`
	BodyGoFunc       func() []GoExpr `ko:"name=bodyGoFunc"`
	Cached_          *AssignCache    `ko:"name=cached"`
}

func (soln *CutSoln) String() string { return Sprint(soln) }

func (soln *CutSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *CutSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		soln.Arg.CircuitEffect(),
		soln.IntoElemUnifier.CircuitEffect(),
		soln.OtherwiseUnifier.CircuitEffect(),
	).AggregateDuctFunc(soln.FuncExpr())
}

func (soln *CutSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveCutArg(span *Span, arg GoStructure) (soln *CutSolnArg, err error) {
	soln = &CutSolnArg{Arg: arg}
	if soln.Into.Extractor, soln.Into.Type, err = GoSelectSimplify(span, Path{"into"}, arg); err != nil {
		return nil, span.Errorf(err, "cut expects an over argument")
	}
	if soln.IntoLifted, soln.IntoLift, err = LiftToSequence(span, soln.Into.Type); err != nil {
		return nil, span.Errorf(err, "cut lifting into to series")
	}
	if soln.Otherwise.Extractor, soln.Otherwise.Type, err = GoSelect(span, Path{"otherwise"}, arg); err != nil {
		return nil, span.Errorf(err, "cut expects an otherwise argument")
	}
	return
}

func (soln *CutSoln) IntoIsUnknown() bool {
	_, ok := soln.Arg.IntoLifted.(Unknown)
	return ok
}

func SolveCut(span *Span, arg GoStructure) (soln *CutSoln, err error) {
	assn := NewAssignCtx(span)
	soln = &CutSoln{Origin: span}
	if soln.Arg, err = SolveCutArg(span, arg); err != nil {
		return nil, err
	}
	switch v := soln.Arg.IntoLifted.(type) {
	case Unknown:
		return
	case *GoSlice: // conditionally return otherwise or slice's first element
		soln.BodyGoFunc = soln.cutSlice
		if soln.Returns, err = Generalize(span, v.Elem, soln.Arg.Otherwise.Type); err != nil {
			return nil, err
		}
		if soln.IntoElemUnifier, err = assn.Assign(v.Elem, soln.Returns); err != nil {
			return nil, err
		}
		if soln.OtherwiseUnifier, err = assn.Assign(soln.Arg.Otherwise.Type, soln.Returns); err != nil {
			return nil, err
		}
	case *GoArray:
		if v.Len == 0 { // return otherwise unconditionally
			soln.BodyGoFunc = soln.cutEmptyArray
			soln.Returns = soln.Arg.Otherwise.Type
		} else { // return array's first element unconditionally
			soln.BodyGoFunc = soln.cutNonEmptyArray
			if soln.Returns, err = Generalize(span, v.Elem, soln.Arg.Otherwise.Type); err != nil {
				return nil, err
			}
		}
		if soln.IntoElemUnifier, err = assn.Assign(v.Elem, soln.Returns); err != nil {
			return nil, err
		}
		if soln.OtherwiseUnifier, err = assn.Assign(soln.Arg.Otherwise.Type, soln.Returns); err != nil {
			return nil, err
		}
	}
	soln.Cached_ = assn.Flush()
	return
}

func (soln *CutSoln) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: soln.Origin,
		Name:   fmt.Sprintf("cut_%s", soln.Origin.SpanID().String()),
	}
}

func (soln *CutSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: soln.FuncExpr(),
		Arg: []GoExpr{
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}

func (soln *CutSoln) FuncExpr() GoFuncExpr {
	return &GoShaperFuncExpr{
		Comment:    fmt.Sprintf("(%s)", soln.Origin.SourceLine()),
		FuncName:   soln.FuncName(),
		ArgName:    soln.argExpr(),
		ArgType:    soln.Arg.Arg,
		ReturnType: soln.Returns,
		Line:       soln.BodyGoFunc(),
	}
}

func (soln *CutSoln) argExpr() GoExpr { return &GoVerbatimExpr{"arg"} }

func (soln *CutSoln) liftedIntoExpr() GoExpr {
	return &GoShapeExpr{
		Shaper: soln.Arg.IntoLift,
		Expr: &GoShapeExpr{
			Shaper: soln.Arg.Into.Extractor,
			Expr:   soln.argExpr(),
		},
	}
}

func (soln *CutSoln) otherwiseExpr() GoExpr {
	return &GoShapeExpr{
		Shaper: soln.Arg.Otherwise.Extractor,
		Expr:   soln.argExpr(),
	}
}

func (soln *CutSoln) cutEmptyArray() []GoExpr { return []GoExpr{soln.returnOtherwise()} }

func (soln *CutSoln) returnOtherwise() GoExpr {
	return &GoReturnExpr{
		&GoShapeExpr{
			Shaper: soln.OtherwiseUnifier,
			Expr:   soln.otherwiseExpr(),
		},
	}
}

func (soln *CutSoln) cutNonEmptyArray() []GoExpr { return []GoExpr{soln.returnCut()} }

func (soln *CutSoln) returnCut() GoExpr {
	return &GoReturnExpr{
		&GoShapeExpr{
			Shaper: soln.IntoElemUnifier,
			Expr: &GoIndexExpr{
				Container: soln.liftedIntoExpr(),
				Index:     ZeroExpr,
			},
		},
	}
}

func (soln *CutSoln) cutSlice() []GoExpr {
	vExpr := &GoVerbatimExpr{"x"}
	return []GoExpr{
		&GoColonAssignExpr{
			Left:  vExpr,
			Right: soln.liftedIntoExpr(),
		},
		&GoIfThenElseExpr{
			If: &GoEqualityExpr{
				Left:  &GoCallExpr{Func: LenExpr, Arg: []GoExpr{vExpr}},
				Right: ZeroExpr,
			},
			Then: []GoExpr{soln.returnOtherwise()},
			Else: []GoExpr{soln.returnCut()},
		},
	}
}
