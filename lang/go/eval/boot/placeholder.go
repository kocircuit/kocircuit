package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type BootPlaceholderMacro struct{}

func (BootPlaceholderMacro) Splay() Tree { panic("boot") }

func (BootPlaceholderMacro) MacroID() string { panic("boot") }

func (BootPlaceholderMacro) Label() string { panic("boot") }

func (BootPlaceholderMacro) MacroSheathString() *string { panic("boot") }

func (BootPlaceholderMacro) Help() string { panic("boot") }

func (BootPlaceholderMacro) Doc() string { panic("boot") }

func (BootPlaceholderMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) { panic("boot") }
