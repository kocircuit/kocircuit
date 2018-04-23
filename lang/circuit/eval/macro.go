package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Macro interface {
	Figure
	MacroID() string // MacroID uniquely identifies the meaning of this macro
	// The top-level circuit step in the frame argument is the circuit invocation step
	// corresponding to this operator invocation.
	Invoke(*Span, Arg) (Return, Effect, error)
	Label() string              // string identifier
	Help() string               // human-readable line of text, used in "ko list"
	MacroSheathString() *string // human-readable line ot text show in debug frames (if not nil)
	Doc() string
}

func RefineMacro(span *Span, macro Macro) *Span {
	return span.Refine(MacroSheath{macro})
}

type MacroSheath struct {
	Macro Macro `ko:"name=macro"`
}

func (sh MacroSheath) SheathID() *ID {
	return PtrID(StringID(sh.Macro.MacroID()))
}

func (sh MacroSheath) SheathLabel() *string {
	return PtrString(sh.Macro.Label())
}

func (sh MacroSheath) SheathString() *string {
	return sh.Macro.MacroSheathString()
}
