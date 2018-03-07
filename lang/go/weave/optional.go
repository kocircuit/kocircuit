package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Optional", new(GoOptionalMacro))
}

type GoOptionalMacro struct{}

func (m GoOptionalMacro) MacroID() string { return m.Help() }

func (m GoOptionalMacro) Label() string { return "optional" }

func (m GoOptionalMacro) MacroSheathString() *string { return PtrString("Optional") }

func (m GoOptionalMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

// eg: (
// 	field1: Optional(Bool(true))
// 	field2: Repeated(Int64(0))
// )
func (GoOptionalMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveOptional(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving optional")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type OptionalSolnArg struct {
	Arg      GoStructure   `ko:"name=arg"`
	Optional TypeExtractor `ko:"name=optional"`
}

func (arg *OptionalSolnArg) CircuitEffect() *GoCircuitEffect {
	return arg.Optional.Extractor.CircuitEffect()
}

func SolveOptionalArg(span *Span, arg GoStructure) (soln *OptionalSolnArg, err error) {
	if monadic := StructureMonadicField(arg); monadic == nil {
		return nil, span.Errorf(nil, "optional expects a monadic argument, got %s", Sprint(arg))
	} else {
		soln = &OptionalSolnArg{Arg: arg}
		if soln.Optional.Extractor, soln.Optional.Type, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg); err != nil {
			return nil, span.Errorf(err, "optional expects a monadic argument")
		}
		return soln, nil
	}
}

type OptionalSoln struct {
	Origin  *Span            `ko:"name=origin"`
	Arg     *OptionalSolnArg `ko:"name=arg"`
	Returns GoType           `ko:"name=returns"`
	Lift    Shaper           `ko:"name=lift"`
	Cached_ *AssignCache     `ko:"name=cached"`
}

func (soln *OptionalSoln) String() string { return Sprint(soln) }

func (soln *OptionalSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *OptionalSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), CircuitEffectIfNotNil(soln.Lift))
}

func (soln *OptionalSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveOptional(span *Span, arg GoStructure) (soln *OptionalSoln, err error) {
	assn := NewAssignCtx(span)
	soln = &OptionalSoln{Origin: span}
	if soln.Arg, err = SolveOptionalArg(span, arg); err != nil {
		return nil, span.Errorf(err, "optional argument")
	}
	switch u := soln.Arg.Optional.Type.(type) {
	case *GoPtr, *GoSlice:
		soln.Returns = u
	default:
		soln.Returns = NewGoPtr(u)
		if soln.Lift, err = assn.Assign(soln.Arg.Optional.Type, soln.Returns); err != nil {
			panic("o")
		}
	}
	return
}

func (soln *OptionalSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	if soln.Lift != nil {
		return &GoShapeExpr{
			Shaper: soln.Lift,
			Expr: &GoShapeExpr{
				Shaper: soln.Arg.Optional.Extractor,
				Expr:   FindSlotExpr(arg, RootSlot{}),
			},
		}
	} else {
		return &GoShapeExpr{
			Shaper: soln.Arg.Optional.Extractor,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		}
	}
}
