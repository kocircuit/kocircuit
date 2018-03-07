package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Hash", new(EvalHashMacro))
}

type EvalHashMacro struct{}

func (m EvalHashMacro) MacroID() string { return m.Help() }

func (m EvalHashMacro) Label() string { return "hash" }

func (m EvalHashMacro) MacroSheathString() *string { return PtrString("Hash") }

func (m EvalHashMacro) Help() string { return "Hash" }

func (EvalHashMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		return BasicStringSymbol(returns.(Symbol).Hash()), nil, nil
	}
}
