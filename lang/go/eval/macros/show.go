package macros

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Show", EvalShowMacro{})
	RegisterEvalMacro("ShowType", EvalShowTypeMacro{})
}

const showDoc = `
Show pretty-prints the values passed to it, whereas
ShowType pretty-prints the types of the values passed to it.`

type EvalShowMacro struct{}

func (m EvalShowMacro) MacroID() string { return m.Help() }

func (m EvalShowMacro) Label() string { return "show" }

func (m EvalShowMacro) MacroSheathString() *string { return PtrString("Show") }

func (m EvalShowMacro) Help() string { return "Show" }

func (m EvalShowMacro) Doc() string { return showDoc }

func (EvalShowMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		fmt.Printf("%v\n", returns)
		return
	}
}

type EvalShowTypeMacro struct{}

func (m EvalShowTypeMacro) MacroID() string { return m.Help() }

func (m EvalShowTypeMacro) Label() string { return "showtype" }

func (m EvalShowTypeMacro) MacroSheathString() *string { return PtrString("ShowType") }

func (m EvalShowTypeMacro) Help() string { return "ShowType" }

func (m EvalShowTypeMacro) Doc() string { return showDoc }

func (EvalShowTypeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		fmt.Printf("%v\n", returns.(Symbol).Type())
		return
	}
}
