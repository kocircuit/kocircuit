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

func (m EvalRangeMacro) Help() string { return "Range(start, over, with, stop)" }

func (m EvalRangeMacro) Doc() string {
	return `Range iterates sequentially over the elements of an input sequence.

A user-supplied iterator function is invoked for each sequence element.
The iterator function is expected to return a structure with two fields: emit and carry.

Values stored in emit are collected across all invocations of the user
function and merged into a new sequence and returned by Range.

The value stored in carry is passed to the next invocation of the iterator function.
The carry returned by the last invocation of the iterator is returned by Range.

Range expects three input arguments:
* over holds the sequence value to be ranged over,
* with holds the user-supplied iterator function,
* start holds the carry value to be passed to the first invocation of the iterator.

Range passes two arguments to the user-supplied iterator:
* elem holds the element currently being iterated over
* carry holds the carry value from the previous invocation of the iterator,
or in the case of the first iteration, it holds the value of start

Range returns a structure with two fields: image and residue.
* image holds a merged sequence of all values emitted by the per-element
invocations to the iterator,
* residue holds the value of the carry returned by the last iterator invocation;
if the input sequence is empty, residue holds the value of start.`
}

func (EvalRangeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	with, ok := a.Walk("with").(*VarietySymbol) // with returns (emit: X, carry: Y) values
	if !ok {
		return nil, nil, span.Errorf(nil, "range with is not a variety")
	}
	over := a.Walk("over")
	carry := a.Walk("start")
	stopArg := a.Walk("stop")
	stop, stopIsVty := stopArg.(*VarietySymbol) // with returns (emit: X, carry: Y) values
	if !IsEmptySymbol(stopArg) && !stopIsVty {
		return nil, nil, span.Errorf(nil, "range stop is not a variety")
	}
	//
	image := Symbols{}
	for i, elem := range over.LiftToSeries(span).Elem {
		iterSpan := RefineOutline(span, fmt.Sprintf("#%d", i))
		fields := Fields{
			{Name: "carry", Shape: carry},
			{Name: "elem", Shape: elem},
		}
		if iterReturns, _, err := with.Evoke(iterSpan, fields); err != nil {
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
		// stop
		if stop != nil {
			stopFields := Fields{{Name: "", Shape: carry}}
			if stopReturns, _, err := stop.Evoke(iterSpan, stopFields); err != nil {
				return nil, nil, err
			} else {
				if stopFlag, ok := AsBasicBool(stopReturns.(Symbol)); !ok {
					return nil, nil, iterSpan.Errorf(nil, "range stop must return boolean, got %v", stopReturns)
				} else {
					if stopFlag {
						break // for loop
					}
				}
			}
		}
	} // for loop
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
