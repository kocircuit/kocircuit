package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Merge", new(EvalMergeMacro))
}

type EvalMergeMacro struct{}

func (m EvalMergeMacro) MacroID() string { return m.Help() }

func (m EvalMergeMacro) Label() string { return "merge" }

func (m EvalMergeMacro) MacroSheathString() *string { return PtrString("Merge") }

func (m EvalMergeMacro) Help() string { return "Merge" }

func (EvalMergeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	// filter subseries that are not empty
	ss := Symbols{}
	for _, a := range arg.(*StructSymbol).SelectMonadic().LiftToSeries(span).Elem {
		for _, e := range a.LiftToSeries(span).Elem {
			if IsEmptySymbol(e) {
				panic("o")
			}
			ss = append(ss, e)
		}
	}
	if merged, err := MakeSeriesSymbol(span, ss); err != nil {
		return nil, nil, span.Errorf(err, "merging sequences")
	} else {
		if merged.IsEmpty() {
			return EmptySymbol{}, nil, nil
		} else {
			return merged, nil, nil
		}
	}
}
