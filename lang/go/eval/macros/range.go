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
	RegisterEvalMacro("Range", new(EvalRangeMacro))
}

type EvalRangeMacro struct{}

func (m EvalRangeMacro) MacroID() string { return m.Help() }

func (m EvalRangeMacro) Label() string { return "range" }

func (m EvalRangeMacro) MacroSheathString() *string { return PtrString("Range") }

func (m EvalRangeMacro) Help() string { return "Range" }

func (EvalRangeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	with, ok := a.Walk("with").(*VarietySymbol) // with returns (emit: X, carry: Y) values
	if !ok {
		return nil, nil, span.Errorf(nil, "range with is not a variety")
	}
	over := a.Walk("over")
	carry := a.Walk("start")
	image := Symbols{}
	for i, elem := range over.LiftToSeries(span).Elem {
		iterSpan := RefineOutline(span, fmt.Sprintf("#%d", i))
		knot := Knot{
			{Name: "carry", Shape: carry, Effect: nil, Frame: iterSpan},
			{Name: "elem", Shape: elem, Effect: nil, Frame: iterSpan},
		}
		if iterReturns, _, err := with.Evoke(iterSpan, knot); err != nil {
			return nil, nil, err
		} else {
			iterResult := iterReturns.(Symbol)
			switch u := iterResult.(type) {
			case EmptySymbol:
				carry = EmptySymbol{}
			case *StructSymbol:
				emitted := u.Walk("emit")
				carry = u.Walk("carry")
				if !IsEmptySymbol(emitted) {
					image = append(image, emitted)
				}
			default:
				return nil, nil, iterSpan.Errorf(nil, "range iterator must return a structure or nothing")
			}
		}
	}
	if imageSeries, err := MakeSeriesSymbol(span, image); err != nil {
		return nil, nil, span.Errorf(err, "range unifying image elements")
	} else {
		return MakeStructSymbol(
			FieldSymbols{
				{Name: "image", Value: imageSeries},
				{Name: "residue", Value: carry},
			},
		), nil, nil
	}
}
