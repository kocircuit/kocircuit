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
	RegisterGoMacro("And", new(GoAndMacro))
}

type GoAndMacro struct{}

func (m GoAndMacro) Label() string { return "and" }

func (m GoAndMacro) MacroSheathString() *string { return PtrString("And") }

func (m GoAndMacro) MacroID() string { return m.Help() }

func (m GoAndMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoAndMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveAnd(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving and")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type AndSolnArg struct {
	Arg GoStructure     `ko:"name=arg"`
	And []TypeExtractor `ko:"name=and"`
}

func (arg *AndSolnArg) CircuitEffect() *GoCircuitEffect {
	effect := make([]*GoCircuitEffect, len(arg.And))
	for i, arg := range arg.And {
		effect[i] = arg.Extractor.CircuitEffect()
	}
	return AggregateCircuitEffects(effect...)
}

func SolveAndArg(span *Span, arg GoStructure) (soln *AndSolnArg, err error) {
	soln = &AndSolnArg{Arg: arg, And: make([]TypeExtractor, len(arg.StructureField()))}
	for i, field := range arg.StructureField() {
		soln.And[i].Name = field.KoName()
		if soln.And[i].Extractor, soln.And[i].Type, err = GoSelectSimplify(span, Path{field.KoName()}, arg); err != nil {
			panic("o")
		}
	}
	return soln, nil
}

type AndSoln struct {
	Origin   *Span        `ko:"name=origin"`
	Arg      *AndSolnArg  `ko:"name=arg"`
	VarIndex []int        `ko:"name=var_index"` // indices of variable (non GoTrue or GoFalse) arguments
	Returns  GoType       `ko:"name=returns"`   // GoBool, GoTrue or GoFalse
	Form     GoSlotForm   `ko:"name=form"`
	Cached_  *AssignCache `ko:"name=cached"`
}

func (soln *AndSoln) String() string { return Sprint(soln) }

func (soln *AndSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *AndSoln) CircuitEffect() *GoCircuitEffect { return soln.Arg.CircuitEffect() }

func (soln *AndSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveAnd(span *Span, arg GoStructure) (soln *AndSoln, err error) {
	assn := NewAssignCtx(span)
	soln = &AndSoln{Origin: span}
	if soln.Arg, err = SolveAndArg(span, arg); err != nil {
		return nil, span.Errorf(err, "and argument")
	}
	falseClauses := false
	for i, arg := range soln.Arg.And { // arg is TypeExtractor
		if _, err = assn.Assign(arg.Type, GoTrue); err == nil {
			// nop
		} else if _, err = assn.Assign(arg.Type, GoFalse); err == nil {
			falseClauses = true
		} else if _, err = assn.Assign(arg.Type, GoBool); err == nil {
			soln.VarIndex = append(soln.VarIndex, i)
		} else {
			return nil, span.Errorf(err, "and argument not boolean")
		}
	}
	soln.Cached_ = assn.Flush()
	if falseClauses {
		soln.Form = &GoInvariantForm{FalseExpr}
		soln.Returns = GoFalse
	} else if len(soln.VarIndex) > 0 {
		soln.Form = (*AndSolnForm)(soln)
		soln.Returns = GoBool
	} else {
		soln.Form = &GoInvariantForm{TrueExpr}
		soln.Returns = GoTrue
	}
	return
}

func (soln *AndSoln) FormExpr(arg ...*GoSlotExpr) GoExpr { return soln.Form.FormExpr(arg...) }

type AndSolnForm AndSoln

func (soln *AndSolnForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	var andExpr GoExpr
	for _, varIndex := range soln.VarIndex {
		argExtractor := &GoShapeExpr{
			Shaper: soln.Arg.And[varIndex].Extractor,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		}
		if andExpr == nil {
			andExpr = argExtractor
		} else {
			andExpr = &GoAndExpr{
				Left:  argExtractor,
				Right: andExpr,
			}
		}
	}
	return andExpr
}
