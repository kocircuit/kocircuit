package boolean

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterGoMacro("Not", new(GoNotMacro))
}

type GoNotMacro struct{}

func (m GoNotMacro) Label() string { return "not" }

func (m GoNotMacro) MacroSheathString() *string { return PtrString("Not") }

func (m GoNotMacro) MacroID() string { return m.Help() }

func (m GoNotMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoNotMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveNot(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving not")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type NotSolnArg struct {
	Arg GoStructure   `ko:"name=arg"`
	Not TypeExtractor `ko:"name=not"`
}

func (arg *NotSolnArg) CircuitEffect() *GoCircuitEffect {
	return arg.Not.Extractor.CircuitEffect()
}

func SolveNotArg(span *Span, arg GoStructure) (soln *NotSolnArg, err error) {
	if monadic := StructureMonadicField(arg); monadic == nil {
		return nil, span.Errorf(nil, "not expects a monadic argument, got %s", Sprint(arg))
	} else {
		soln = &NotSolnArg{Arg: arg}
		if soln.Not.Extractor, soln.Not.Type, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg); err != nil {
			return nil, span.Errorf(err, "not expects an etc argument")
		}
		return soln, nil
	}
}

type NotSoln struct {
	Origin  *Span        `ko:"name=origin"`
	Arg     *NotSolnArg  `ko:"name=arg"`
	Returns GoType       `ko:"name=returns"`
	Boolify Shaper       `ko:"name=boolify"`
	Form    GoSlotForm   `ko:"name=form"`
	Cached_ *AssignCache `ko:"name=cached"`
}

func (soln *NotSoln) String() string { return Sprint(soln) }

func (soln *NotSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *NotSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), CircuitEffectIfNotNil(soln.Boolify))
}

func (soln *NotSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveNot(span *Span, arg GoStructure) (soln *NotSoln, err error) {
	assn := NewAssignCtx(span)
	soln = &NotSoln{Origin: span}
	if soln.Arg, err = SolveNotArg(span, arg); err != nil {
		return nil, span.Errorf(err, "not argument")
	}
	if _, err = assn.Assign(soln.Arg.Not.Type, GoTrue); err == nil {
		soln.Form = &GoInvariantForm{FalseExpr}
		soln.Returns = GoFalse
	} else if _, err = assn.Assign(soln.Arg.Not.Type, GoFalse); err == nil {
		soln.Form = &GoInvariantForm{TrueExpr}
		soln.Returns = GoTrue
	} else if soln.Boolify, err = assn.Assign(soln.Arg.Not.Type, GoBool); err == nil {
		soln.Form = (*NotSolnForm)(soln)
		soln.Returns = GoBool
	} else {
		return nil, span.Errorf(err, "not argument not boolean")
	}
	soln.Cached_ = assn.Flush()
	return
}

func (soln *NotSoln) FormExpr(arg ...*GoSlotExpr) GoExpr { return soln.Form.FormExpr(arg...) }

type NotSolnForm NotSoln

func (soln *NotSolnForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoNotExpr{
		&GoShapeExpr{
			Shaper: soln.Boolify,
			Expr: &GoShapeExpr{
				Shaper: soln.Arg.Not.Extractor,
				Expr:   FindSlotExpr(arg, RootSlot{}),
			},
		},
	}
}
