package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("", new(EvalJoinMacro))
}

type EvalJoinMacro struct{}

func (m EvalJoinMacro) MacroID() string { return m.Help() }

func (m EvalJoinMacro) Label() string { return "join" }

func (m EvalJoinMacro) MacroSheathString() *string { return PtrString("Join") }

func (m EvalJoinMacro) Help() string { return "Join" }

func (m EvalJoinMacro) Doc() string {
	return `Join returns its argument if called with a default (unnamed argument).
Otherwise Join returns a structure containing the named arguments passed to it.`
}

func (EvalJoinMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return Construct(arg.(*StructSymbol)), nil, nil
}

func Construct(a *StructSymbol) Symbol {
	if monadic, hasMonadic := a.GetMonadic(); hasMonadic {
		return monadic
	} else {
		return FilterEmptyStructFields(a)
	}
}
