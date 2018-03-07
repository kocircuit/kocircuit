package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Take", new(EvalTakeMacro))
}

type EvalTakeMacro struct{}

func (m EvalTakeMacro) MacroID() string { return m.Help() }

func (m EvalTakeMacro) Label() string { return "take" }

func (m EvalTakeMacro) MacroSheathString() *string { return PtrString("Take") }

func (m EvalTakeMacro) Help() string { return "Take" }

func (EvalTakeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	if fromSeries := LiftToSeries(span, a.Walk("from")); fromSeries.Len() > 0 {
		return fromSeries.Elem[0], nil, nil
	} else {
		return a.Walk("otherwise"), nil, nil
	}
}
