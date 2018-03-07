package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Bug", EvalBugMacro{})
}

type EvalBugMacro struct{}

func (m EvalBugMacro) MacroID() string { return m.Help() }

func (m EvalBugMacro) Label() string { return "bug" }

func (m EvalBugMacro) MacroSheathString() *string { return PtrString("Bug") }

func (m EvalBugMacro) Help() string { return "Bug" }

func (EvalBugMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		span.Fatalf(nil, "BUG %v\n", returns)
		return
	}
}
