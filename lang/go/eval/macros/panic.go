package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Panic", new(EvalPanicMacro))
	RegisterEvalMacro("Recover", new(EvalRecoverMacro))
}

type EvalPanicMacro struct{}

func (m EvalPanicMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalPanicMacro) MacroID() string { return m.Help() }

func (m EvalPanicMacro) Label() string { return "panic" }

func (m EvalPanicMacro) MacroSheathString() *string { return PtrString("Panic") }

func (m EvalPanicMacro) Help() string {
	return "Panic"
}

func (EvalPanicMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	panic(arg.(*StructSymbol).SelectMonadic())
}

type EvalRecoverMacro struct{}

func (m EvalRecoverMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalRecoverMacro) MacroID() string { return m.Help() }

func (m EvalRecoverMacro) Label() string { return "recover" }

func (m EvalRecoverMacro) MacroSheathString() *string { return PtrString("Recover") }

func (m EvalRecoverMacro) Help() string {
	return "Recover"
}

// Recover(invoke, panic)
func (EvalRecoverMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	invokeVty, ok := a.Walk("invoke").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "recover invoke is not a variety")
	}
	panicVty, ok := a.Walk("panic").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "recover panic is not a variety")
	}
	defer func() {
		if r := recover(); r != nil {
			panicKnot := Knot{{Name: "", Shape: r.(Symbol), Effect: nil, Frame: span}}
			returns, effect, err = panicVty.Evoke(span, panicKnot)
			return
		}
	}()
	return invokeVty.Invoke(span)
}
