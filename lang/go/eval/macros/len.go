package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Len", new(EvalLenMacro))
}

type EvalLenMacro struct{}

func (m EvalLenMacro) MacroID() string { return m.Help() }

func (m EvalLenMacro) Label() string { return "len" }

func (m EvalLenMacro) MacroSheathString() *string { return PtrString("Len") }

func (m EvalLenMacro) Help() string { return "Len" }

func (m EvalLenMacro) Doc() string {
	return `Len returns the length of the sequence passed to it.`
}

func (EvalLenMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	series := arg.(*StructSymbol).SelectMonadic().LiftToSeries(span)
	return BasicInt64Symbol(int64(series.Len())), nil, nil
}
