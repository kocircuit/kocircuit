package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Hang", new(EvalHangMacro))
}

type EvalHangMacro struct{}

func (m EvalHangMacro) MacroID() string { return m.Help() }

func (m EvalHangMacro) Label() string { return "hang" }

func (m EvalHangMacro) MacroSheathString() *string { return PtrString("Hang") }

func (m EvalHangMacro) Help() string { return "Hang" }

func (EvalHangMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	select {} // wait forever
}
