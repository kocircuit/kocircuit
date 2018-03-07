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
	RegisterGoMacro("Or", new(GoOrMacro))
}

type GoOrMacro struct{}

func (m GoOrMacro) Label() string { return "or" }

func (m GoOrMacro) MacroSheathString() *string { return PtrString("Or") }

func (m GoOrMacro) MacroID() string { return m.Help() }

func (m GoOrMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoOrMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveOr(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving or")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type OrSolnArg struct {
	Arg GoStructure     `ko:"name=arg"`
	Or  []TypeExtractor `ko:"name=or"`
}

func (arg *OrSolnArg) CircuitEffect() *GoCircuitEffect {
	effect := make([]*GoCircuitEffect, len(arg.Or))
	for i, arg := range arg.Or {
		effect[i] = arg.Extractor.CircuitEffect()
	}
	return AggregateCircuitEffects(effect...)
}

func SolveOrArg(span *Span, arg GoStructure) (soln *OrSolnArg, err error) {
	soln = &OrSolnArg{Arg: arg, Or: make([]TypeExtractor, len(arg.StructureField()))}
	for i, field := range arg.StructureField() {
		soln.Or[i].Name = field.KoName()
		if soln.Or[i].Extractor, soln.Or[i].Type, err = GoSelectSimplify(span, Path{field.KoName()}, arg); err != nil {
			panic("o")
		}
	}
	return soln, nil
}

type OrSoln struct {
	Origin   *Span        `ko:"name=origin"`
	Arg      *OrSolnArg   `ko:"name=arg"`
	VarIndex []int        `ko:"name=var_index"` // indices of variable (non GoTrue or GoFalse) arguments
	Returns  GoType       `ko:"name=returns"`   // GoBool, GoTrue or GoFalse
	Form     GoSlotForm   `ko:"name=form"`
	Cached_  *AssignCache `ko:"name=cached"`
}

func (soln *OrSoln) String() string { return Sprint(soln) }

func (soln *OrSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *OrSoln) CircuitEffect() *GoCircuitEffect { return soln.Arg.CircuitEffect() }

func (soln *OrSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveOr(span *Span, arg GoStructure) (soln *OrSoln, err error) {
	assn := NewAssignCtx(span)
	soln = &OrSoln{Origin: span}
	if soln.Arg, err = SolveOrArg(span, arg); err != nil {
		return nil, span.Errorf(err, "or argument")
	}
	trueClauses := false
	for i, arg := range soln.Arg.Or { // arg is TypeExtractor
		if _, err = assn.Assign(arg.Type, GoTrue); err == nil {
			trueClauses = true
		} else if _, err = assn.Assign(arg.Type, GoFalse); err == nil {
			// nop
		} else if _, err = assn.Assign(arg.Type, GoBool); err == nil {
			soln.VarIndex = append(soln.VarIndex, i)
		} else {
			return nil, span.Errorf(err, "or argument not boolean")
		}
	}
	soln.Cached_ = assn.Flush()
	if trueClauses {
		soln.Form = &GoInvariantForm{TrueExpr}
		soln.Returns = GoTrue
	} else if len(soln.VarIndex) > 0 {
		soln.Form = (*OrSolnForm)(soln)
		soln.Returns = GoBool
	} else {
		soln.Form = &GoInvariantForm{FalseExpr}
		soln.Returns = GoFalse
	}
	return
}

func (soln *OrSoln) FormExpr(arg ...*GoSlotExpr) GoExpr { return soln.Form.FormExpr(arg...) }

type OrSolnForm OrSoln

func (soln *OrSolnForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	var orExpr GoExpr
	for _, varIndex := range soln.VarIndex {
		argExtractor := &GoShapeExpr{
			Shaper: soln.Arg.Or[varIndex].Extractor,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		}
		if orExpr == nil {
			orExpr = argExtractor
		} else {
			orExpr = &GoOrExpr{
				Left:  argExtractor,
				Right: orExpr,
			}
		}
	}
	return orExpr
}
