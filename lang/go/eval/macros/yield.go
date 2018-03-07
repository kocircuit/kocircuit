package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Yield", new(EvalYieldMacro))
}

type EvalYieldMacro struct{}

func (m EvalYieldMacro) MacroID() string { return m.Help() }

func (m EvalYieldMacro) Label() string { return "yield" }

func (m EvalYieldMacro) MacroSheathString() *string { return PtrString("Yield") }

func (m EvalYieldMacro) Help() string { return "Yield" }

func (EvalYieldMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	ifArg, ifBool := AsBasicBool(a.Walk("if"))
	if !ifBool {
		return nil, nil, span.Errorf(err, "yield if must be boolean")
	}
	if ifArg {
		return a.Walk("then"), nil, nil
	} else {
		return a.Walk("else"), nil, nil
	}
}
