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
	if fromSeries := ExtractMonadicOrNamed(arg, "from").LiftToSeries(span); fromSeries.Len() > 0 {
		if remainder, err := MakeSeriesSymbol(span, fromSeries.Elem[1:]); err != nil {
			panic(err)
		} else {
			return MakeStructSymbol(
				FieldSymbols{
					{Name: "first", Value: fromSeries.Elem[0]},
					{Name: "remainder", Value: remainder},
				},
			), nil, nil
		}
	} else {
		return EmptySymbol{}, nil, nil
	}
}

func ExtractNamed(arg Arg, name string) Symbol {
	a := arg.(*StructSymbol)
	return a.Walk(name)
}

func ExtractMonadicOrNamed(arg Arg, name string) Symbol {
	a := arg.(*StructSymbol)
	if v := a.SelectMonadic(); !IsEmptySymbol(v) {
		return v
	} else {
		return a.Walk("name")
	}
}
