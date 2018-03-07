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
	RegisterEvalMacro("Loop", new(EvalLoopMacro))
}

type EvalLoopMacro struct{}

func (m EvalLoopMacro) MacroID() string { return m.Help() }

func (m EvalLoopMacro) Label() string { return "loop" }

func (m EvalLoopMacro) MacroSheathString() *string { return PtrString("Loop") }

func (m EvalLoopMacro) Help() string { return "Loop" }

// Loop(start:█, step:█, stop:█)
func (EvalLoopMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	step, ok := a.Walk("step").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "loop step is not a variety")
	}
	var stop *VarietySymbol
	switch u := a.Walk("stop").(type) {
	case EmptySymbol:
	case *VarietySymbol:
		stop = u
	default:
		return nil, nil, span.Errorf(nil, "loop stop is not a variety, it is %v", u)
	}
	carry := a.Walk("start")
	for j := 0; true; j++ {
		iterSpan := RefineOutline(span, fmt.Sprintf("#%d", j))
		// step
		stepKnot := Knot{{Name: "", Shape: carry, Effect: nil, Frame: iterSpan}}
		if stepReturns, _, err := step.Evoke(iterSpan, stepKnot); err != nil {
			return nil, nil, err
		} else {
			carry = stepReturns.(Symbol)
			// stop
			if stop != nil {
				stopKnot := Knot{{Name: "", Shape: carry, Effect: nil, Frame: iterSpan}}
				if stopReturns, _, err := stop.Evoke(iterSpan, stopKnot); err != nil {
					return nil, nil, err
				} else {
					if stopFlag, ok := AsBasicBool(stopReturns.(Symbol)); !ok {
						return nil, nil, iterSpan.Errorf(nil, "loop stop must return boolean, got %v", stopReturns)
					} else {
						if stopFlag {
							break // for loop
						}
					}
				}
			}
		}
	}
	return carry, nil, nil
}
