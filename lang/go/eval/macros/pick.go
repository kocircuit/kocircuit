package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Pick", new(EvalPickMacro))
}

type EvalPickMacro struct{}

func (m EvalPickMacro) MacroID() string { return m.Help() }

func (m EvalPickMacro) Label() string { return "pick" }

func (m EvalPickMacro) MacroSheathString() *string { return PtrString("Pick") }

func (m EvalPickMacro) Help() string { return "Pick" }

func (EvalPickMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	if either := a.Walk("either"); !IsEmptySymbol(either) {
		return either, nil, nil
	} else {
		return a.Walk("or"), nil, nil
	}
}
