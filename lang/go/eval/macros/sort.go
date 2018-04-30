package macros

import (
	"sort"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Sort", new(EvalSortMacro))
}

type EvalSortMacro struct{}

func (m EvalSortMacro) MacroID() string { return m.Help() }

func (m EvalSortMacro) Label() string { return "sort" }

func (m EvalSortMacro) MacroSheathString() *string { return PtrString("Sort") }

func (m EvalSortMacro) Help() string { return "Sort" }

func (m EvalSortMacro) Doc() string {
	return `Sort(over, less) sorts the sequence of values over, using the comparator less(left, right).`
}

// Sort(over, less)
func (EvalSortMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	lessArg := a.Walk("less")
	less, ok := lessArg.(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "sort less argument must be a variety, got %v", lessArg)
	}
	ss := &sortSymbols{
		symbols: FilterEmptySymbols(a.Walk("over").LiftToSeries(span).Elem),
		less:    less,
		span:    span,
	}
	defer func() {
		if r := recover(); r != nil {
			switch p := r.(type) {
			case *lessErrorPanic:
				returns, effect, err = nil, nil, p.Error
			case *EvalPanic:
				panic(p) // forward eval panics
			default:
				panic(p) // forward unknown panics
			}
		}
	}()
	sort.Sort(ss)
	if result, err := MakeSeriesSymbol(span, ss.symbols); err != nil {
		panic("o")
	} else {
		return result, nil, nil
	}
}

type sortSymbols struct {
	symbols Symbols
	less    *VarietySymbol
	span    *Span
}

func (ss *sortSymbols) Len() int {
	return len(ss.symbols)
}

func (ss *sortSymbols) Swap(i, j int) {
	ss.symbols[i], ss.symbols[j] = ss.symbols[j], ss.symbols[i]
}

type lessErrorPanic struct {
	Error error
}

func (ss *sortSymbols) Less(i, j int) bool {
	knot := Fields{
		{Name: "left", Shape: ss.symbols[i], Effect: nil, Frame: ss.span},
		{Name: "right", Shape: ss.symbols[j], Effect: nil, Frame: ss.span},
	}
	if lessReturns, _, err := ss.less.Evoke(ss.span, knot); err != nil {
		panic(&lessErrorPanic{Error: err})
	} else if boolValue, isBool := AsBasicBool(lessReturns); !isBool {
		panic(&lessErrorPanic{Error: ss.span.Errorf(nil, "less variety must return boolean, got %v", lessReturns)})
	} else {
		return boolValue
	}
}
