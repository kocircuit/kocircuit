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
	RegisterGoMacro("Yield", new(GoYieldMacro))
}

type GoYieldMacro struct{}

func (m GoYieldMacro) MacroID() string { return m.Help() }

func (m GoYieldMacro) Label() string { return "yield" }

func (m GoYieldMacro) MacroSheathString() *string { return PtrString("Yield") }

func (m GoYieldMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

// Emitted effect is GoMacroEffect.
func (GoYieldMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveYield(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving yield")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type YieldSolnArg struct {
	Arg  GoStructure   `ko:"name=arg"`
	If   TypeExtractor `ko:"name=if"`
	Then TypeExtractor `ko:"name=then"`
	Else TypeExtractor `ko:"name=else"`
}

func (arg *YieldSolnArg) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		arg.If.Extractor.CircuitEffect(),
		arg.Then.Extractor.CircuitEffect(),
		arg.Else.Extractor.CircuitEffect(),
	)
}

func SolveYieldArg(span *Span, arg GoStructure) (soln *YieldSolnArg, err error) {
	soln = &YieldSolnArg{Arg: arg}
	if soln.If.Extractor, soln.If.Type, err = GoSelectSimplify(span, Path{"if"}, arg); err != nil {
		return nil, span.Errorf(err, "yield expects an “if” argument")
	}
	if soln.Then.Extractor, soln.Then.Type, err = GoSelectSimplify(span, Path{"then"}, arg); err != nil {
		return nil, span.Errorf(err, "yield expects a “then” argument")
	}
	if soln.Else.Extractor, soln.Else.Type, err = GoSelectSimplify(span, Path{"else"}, arg); err != nil {
		return nil, span.Errorf(err, "yield expects an “else” argument")
	}
	return soln, nil
}

type YieldSoln struct {
	Origin      *Span         `ko:"name=origin"`
	Arg         *YieldSolnArg `ko:"name=arg"`
	Returns     GoType        `ko:"name=returns"`
	IfShaper    Shaper        `ko:"name=if_shaper"`
	Determined  *bool         `ko:"name=determined"`
	ThenUnifier Shaper        `ko:"name=then_unifier"`
	ElseUnifier Shaper        `ko:"name=else_unifier"`
	Cached_     *AssignCache  `ko:"name=cached"`
}

func (soln *YieldSoln) String() string { return Sprint(soln) }

func (soln *YieldSoln) Cached() *AssignCache { return soln.Cached_ }

func (soln *YieldSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		soln.Arg.CircuitEffect(),
		soln.IfShaper.CircuitEffect(),
		CircuitEffectIfNotNil(soln.ThenUnifier),
		CircuitEffectIfNotNil(soln.ElseUnifier),
	).AggregateDuctType(soln.Returns).AggregateDuctFunc(soln.FuncExpr())
}

func (soln *YieldSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveYield(span *Span, argStruct GoStructure) (soln *YieldSoln, err error) {
	soln = &YieldSoln{Origin: span}
	assn := NewAssignCtx(span)
	defer func() {
		if err == nil {
			soln.Cached_ = assn.Flush()
		}
	}()
	if soln.Arg, err = SolveYieldArg(span, argStruct); err != nil {
		return nil, err
	}
	if soln.IfShaper, err = assn.Assign(soln.Arg.If.Type, GoBool); err != nil {
		return nil, span.Errorf(err, "yield if argument not assignable to bool")
	}
	// check whether if-condition determinable at weaving time
	if _, err = assn.Assign(soln.Arg.If.Type, GoTrue); err == nil {
		soln.Determined = PtrBool(true)
		soln.Returns = soln.Arg.Then.Type
		return soln, nil
	}
	if _, err = assn.Assign(soln.Arg.If.Type, GoFalse); err == nil {
		soln.Determined = PtrBool(false)
		soln.Returns = soln.Arg.Else.Type
		return soln, nil
	}
	// generalize then and else
	if soln.Returns, err = Generalize(span, soln.Arg.Then.Type, soln.Arg.Else.Type); err != nil {
		return nil, span.Errorf(err, "yield then-else do not generalize")
	}
	if soln.ThenUnifier, err = assn.Assign(soln.Arg.Then.Type, soln.Returns); err != nil {
		panic("o")
	}
	if soln.ElseUnifier, err = assn.Assign(soln.Arg.Else.Type, soln.Returns); err != nil {
		panic("o")
	}
	return
}

func (soln *YieldSoln) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: soln.Origin,
		Name:   fmt.Sprintf("yield_%s", soln.Origin.SpanID().String()),
	}
}

// FormExpr returns an invocation of the yield_ function with arg[0]
func (soln *YieldSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: soln.FuncExpr(),
		Arg: []GoExpr{
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}

func (soln *YieldSoln) FuncExpr() GoFuncExpr {
	if soln.Determined != nil {
		if *soln.Determined { // if = true
			return soln.thenFuncExpr()
		} else { // if = false
			return soln.elseFuncExpr()
		}
	} else {
		return soln.thenElseFuncExpr()
	}
}

func (soln *YieldSoln) argExpr() GoExpr {
	return &GoVerbatimExpr{"arg"}
}

func (soln *YieldSoln) ifExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.If.Extractor, Expr: soln.argExpr()}
}

func (soln *YieldSoln) thenExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.Then.Extractor, Expr: soln.argExpr()}
}

func (soln *YieldSoln) elseExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.Else.Extractor, Expr: soln.argExpr()}
}

func (soln *YieldSoln) buildFuncExpr(line []GoExpr) GoFuncExpr {
	return &GoShaperFuncExpr{
		Comment:    fmt.Sprintf("(%s)", soln.Origin.SourceLine()),
		FuncName:   soln.FuncName(),
		ArgName:    soln.argExpr(),
		ArgType:    soln.Arg.Arg,
		ReturnType: soln.Returns,
		Line:       line,
	}
}

func (soln *YieldSoln) thenFuncExpr() GoFuncExpr {
	return soln.buildFuncExpr(
		[]GoExpr{
			&GoReturnExpr{soln.thenExpr()},
		},
	)
}

func (soln *YieldSoln) elseFuncExpr() GoFuncExpr {
	return soln.buildFuncExpr(
		[]GoExpr{
			&GoReturnExpr{soln.elseExpr()},
		},
	)
}

func (soln *YieldSoln) thenElseFuncExpr() GoFuncExpr {
	return soln.buildFuncExpr(
		[]GoExpr{
			&GoIfThenElseExpr{
				If: &GoShapeExpr{Shaper: soln.IfShaper, Expr: soln.ifExpr()},
				Then: []GoExpr{
					&GoReturnExpr{
						&GoShapeExpr{Shaper: soln.ThenUnifier, Expr: soln.thenExpr()},
					},
				},
				Else: []GoExpr{
					&GoReturnExpr{
						&GoShapeExpr{Shaper: soln.ElseUnifier, Expr: soln.elseExpr()},
					},
				},
			},
		},
	)
}
