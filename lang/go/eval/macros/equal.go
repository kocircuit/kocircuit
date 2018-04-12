package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Equal", new(EvalEqualMacro))
}

type EvalEqualMacro struct{}

func (m EvalEqualMacro) MacroID() string { return m.Help() }

func (m EvalEqualMacro) Label() string { return "equal" }

func (m EvalEqualMacro) MacroSheathString() *string { return PtrString("Equal") }

func (m EvalEqualMacro) Help() string { return "Equal" }

func (EvalEqualMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	series := arg.(*StructSymbol).SelectMonadic().LiftToSeries(span)
	for i := 1; i < len(series.Elem); i++ {
		if !series.Elem[i-1].Equal(span, series.Elem[i]) {
			return BasicFalse, nil, nil
		}
	}
	return BasicTrue, nil, nil
}
