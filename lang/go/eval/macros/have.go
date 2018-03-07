package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Have", new(EvalHaveMacro))
}

type EvalHaveMacro struct{}

func (m EvalHaveMacro) MacroID() string { return m.Help() }

func (m EvalHaveMacro) Label() string { return "have" }

func (m EvalHaveMacro) MacroSheathString() *string { return PtrString("Have") }

func (m EvalHaveMacro) Help() string { return "Have" }

func (EvalHaveMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if !IsEmptySymbol(arg.(*StructSymbol).SelectMonadic()) {
		return BasicTrue, nil, nil
	} else {
		return BasicFalse, nil, nil
	}
}
