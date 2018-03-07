package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Peek", EvalPeekMacro{})
	RegisterEvalMacro("PeekType", EvalPeekTypeMacro{})
}

type EvalPeekMacro struct{}

func (m EvalPeekMacro) MacroID() string { return m.Help() }

func (m EvalPeekMacro) Label() string { return "peek" }

func (m EvalPeekMacro) MacroSheathString() *string { return PtrString("Peek") }

func (m EvalPeekMacro) Help() string { return "Peek" }

func (EvalPeekMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		span.Printf("%v\n", returns)
		return
	}
}

type EvalPeekTypeMacro struct{}

func (m EvalPeekTypeMacro) MacroID() string { return m.Help() }

func (m EvalPeekTypeMacro) Label() string { return "peektype" }

func (m EvalPeekTypeMacro) MacroSheathString() *string { return PtrString("PeekType") }

func (m EvalPeekTypeMacro) Help() string { return "PeekType" }

func (EvalPeekTypeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		span.Printf("%v\n", returns.(Symbol).Type())
		return
	}
}
