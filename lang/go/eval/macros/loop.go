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

func (m EvalLoopMacro) Doc() string {
	return `
The builtin function Loop is a mechanism for running a user function in a loop, providing ways
to carry values from one invocation to the next and an optional stopping condition.

Loop expects three arguments named: start, step and stop.

* start is the initial carry value,
* step is a variety (functional value) which accepts a default (aka monadic) argument,
* stop is an optional variety which accepts a default argument and returns a boolean.

Loop invokes step in a loop. On the first invocation, Loop passes the value
of start to the default argument of step. On following invocations, Loop passes
the value returned by the previous invocation of step. We call the value returned 
by step a carry value.

If stop is provided, Loop will call stop after each invocation of step passing
it the value returned by the immediately preceding step invocation.

* If stop returns true, looping stops and Loop returns the last carry value
(which triggered stop to return true.)

* If stop returns false, looping continues.

If stop is not provided, looping will never end.`
}

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
