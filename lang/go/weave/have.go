package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Have", new(GoHaveMacro))
}

type GoHaveMacro struct{}

func (m GoHaveMacro) MacroID() string { return m.Help() }

func (m GoHaveMacro) Label() string { return "have" }

func (m GoHaveMacro) MacroSheathString() *string { return PtrString("Have") }

func (m GoHaveMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoHaveMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveHave(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving have")
	}
	if soln.IsUnknown() {
		return GoUnknownMacro{}.Invoke(span, arg)
	}
	return soln.Returns(), SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type HaveSolnArg struct {
	Arg     GoStructure   `ko:"name=arg"`
	Monadic TypeExtractor `ko:"name=monadic"`
}

func (arg *HaveSolnArg) CircuitEffect() *GoCircuitEffect {
	return arg.Monadic.Extractor.CircuitEffect()
}

type HaveSoln struct {
	Origin *Span          `ko:"name=origin"`
	Arg    *HaveSolnArg   `ko:"name=arg"`
	Form   GoSlotFormExpr `ko:"name=form"`
}

func (soln *HaveSoln) String() string { return Sprint(soln) }

func (soln *HaveSoln) Cached() *AssignCache { return nil }

func (soln *HaveSoln) CircuitEffect() *GoCircuitEffect { return soln.Arg.CircuitEffect() }

func (soln *HaveSoln) ProgramEffect() *GoProgramEffect { return nil }

func (soln *HaveSoln) Returns() GoType {
	switch soln.Arg.Monadic.Type.(type) {
	case Unknown:
		panic("o")
	case *GoEmpty:
		return GoFalse
	case *GoPtr:
		return GoBool
	case *GoSlice:
		return GoBool
	default:
		return GoTrue
	}
}

func SolveHaveArg(span *Span, arg GoStructure) (soln *HaveSolnArg, err error) {
	soln = &HaveSolnArg{Arg: arg}
	if monadic := StructureMonadicField(arg); monadic != nil {
		if soln.Monadic.Extractor, soln.Monadic.Type, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg); err != nil {
			return nil, span.Errorf(err, "have expects a monadic argument")
		}
	} else {
		empty := NewGoEmpty(span)
		soln.Monadic.Type = empty
		soln.Monadic.Extractor = &IrreversibleEraseShaper{
			Shaping: Shaping{Origin: span, From: empty, To: empty},
		}
	}
	return
}

func (soln *HaveSoln) IsUnknown() bool {
	_, ok := soln.Arg.Monadic.Type.(Unknown)
	return ok
}

func SolveHave(span *Span, arg GoStructure) (soln *HaveSoln, err error) {
	soln = &HaveSoln{Origin: span}
	if soln.Arg, err = SolveHaveArg(span, arg); err != nil {
		return nil, err
	}
	return
}

func (soln *HaveSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	switch soln.Arg.Monadic.Type.(type) {
	case Unknown:
		panic("o")
	case *GoEmpty:
		return FalseExpr
	case *GoPtr:
		return &GoInequalityExpr{
			Left:  FindSlotExpr(arg, RootSlot{}),
			Right: NilExpr,
		}
	case *GoSlice:
		return &GoInequalityExpr{
			Left: &GoCallExpr{
				Func: LenExpr,
				Arg:  []GoExpr{FindSlotExpr(arg, RootSlot{})},
			},
			Right: ZeroExpr,
		}
	default:
		return FalseExpr
	}
}
