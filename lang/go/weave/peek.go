package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Peek", new(GoPeekMacro))
}

// GoPeekMacro is a Ko operator that joins its arguments into a structure and displays this structure during weaving.
type GoPeekMacro struct{}

func (m GoPeekMacro) MacroID() string { return m.Help() }

func (m GoPeekMacro) Label() string { return "peek" }

func (m GoPeekMacro) MacroSheathString() *string { return PtrString("Peek") }

func (m GoPeekMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoPeekMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := GoJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		span.Printf("PEEK %s\n", Sprint(returns))
		return
	}
}
