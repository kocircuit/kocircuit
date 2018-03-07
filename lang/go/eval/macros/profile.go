package macros

import (
	"time"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Profile", new(EvalProfileMacro))
}

type EvalProfileMacro struct{}

func (m EvalProfileMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalProfileMacro) MacroID() string { return m.Help() }

func (m EvalProfileMacro) Label() string { return "profile" }

func (m EvalProfileMacro) MacroSheathString() *string { return PtrString("Profile") }

func (m EvalProfileMacro) Help() string { return "Profile" }

func (EvalProfileMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol).SelectMonadic()
	if vty, ok := a.(*VarietySymbol); !ok {
		return nil, nil, span.Errorf(nil, "profile cannot be applied to a non-variety %v", vty)
	} else {
		t0 := time.Now()
		returns, effect, err = vty.Invoke(span)
		t1 := time.Now()
		dur := t1.Sub(t0)
		if err == nil {
			span.Printf("%1.3fs\n", dur.Seconds())
		}
		return returns, effect, err
	}
}
