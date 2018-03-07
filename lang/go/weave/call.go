package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoCallMacro struct {
	Valve *GoValve `ko:"name=valve"`
}

func (m *GoCallMacro) Splay() Tree { return Quote{m.Help()} }

func (m *GoCallMacro) MacroID() string { return m.Help() }

func (m *GoCallMacro) Label() string { return "call" }

func (m *GoCallMacro) MacroSheathString() *string { return PtrString("Call") }

func (m *GoCallMacro) Help() string {
	return fmt.Sprintf("Call(%s)", m.Valve.Address.String())
}

// Argument arg is a GoStructure. Emitted effect is GoMacroEffect.
func (call *GoCallMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if _, cached, err := Assign(span, arg.(GoType), call.Valve.Real()); err != nil {
		return nil, nil, span.Errorf(err, "arguments do not meet valve")
	} else {
		return call.Valve.Returns,
			&GoMacroEffect{
				Arg:      call.Valve.Real(),
				SlotForm: &GoCallForm{},
				Cached:   cached,
			}, nil
	}
}

type GoCallForm struct{}

func (playForm *GoCallForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoFmtExpr{
		Fmt: "%s.Play(step_ctx)",
		Arg: []GoExpr{
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}
