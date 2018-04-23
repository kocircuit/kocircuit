package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Placeholder", new(EvalPlaceholderMacro))
}

type EvalPlaceholderMacro struct{}

func (m EvalPlaceholderMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalPlaceholderMacro) MacroID() string { return m.Help() }

func (m EvalPlaceholderMacro) Label() string { return "placeholder" }

func (m EvalPlaceholderMacro) MacroSheathString() *string { return PtrString("Placeholder") }

func (m EvalPlaceholderMacro) Help() string { return "Placeholder" }

func (m EvalPlaceholderMacro) Doc() string { return "Placeholder." }

func (EvalPlaceholderMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	panic("placeholder macro")
}
