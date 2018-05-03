package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type WeavePlaceholderMacro struct{}

func (WeavePlaceholderMacro) Splay() Tree { panic("weave") }

func (WeavePlaceholderMacro) MacroID() string { panic("weave") }

func (WeavePlaceholderMacro) Label() string { panic("weave") }

func (WeavePlaceholderMacro) MacroSheathString() *string { panic("weave") }

func (WeavePlaceholderMacro) Help() string { panic("weave") }

func (WeavePlaceholderMacro) Doc() string { panic("weave") }

func (WeavePlaceholderMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) { panic("weave") }
