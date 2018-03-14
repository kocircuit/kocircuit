package macros

import (
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Format", new(EvalFormatMacro))
}

type EvalFormatMacro struct{}

func (m EvalFormatMacro) MacroID() string { return m.Help() }

func (m EvalFormatMacro) Label() string { return "format" }

func (m EvalFormatMacro) MacroSheathString() *string { return PtrString("Format") }

func (m EvalFormatMacro) Help() string { return "Format" }

// Format(format:█, args:█, withString:█, withArg:█)
func (EvalFormatMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	// parse arguments
	a := arg.(*StructSymbol)
	format, ok := AsBasicString(a.Walk("format"))
	if !ok {
		return nil, nil, span.Errorf(nil, "format format argument must be string")
	}
	args := a.Walk("args").LiftToSeries(span)
	withString, ok := a.Walk("withString").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "format withString is not a variety")
	}
	withArg, ok := a.Walk("withArg").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "format withArg is not a variety")
	}
	// orchestrate transformation
	parts := strings.Split(format, `%%`)
	if len(parts) != args.Len()+1 {
		return nil, nil, span.Errorf(nil, "format expects %d arguments, got %d", len(parts)-1, args.Len())
	}
	result := make(Symbols, 0, len(parts)+args.Len())
	for i, p := range parts {
		stringKnot := Knot{{Name: "", Shape: BasicStringSymbol(p), Effect: nil, Frame: span}}
		if formResult, _, err := withString.Evoke(span, stringKnot); err != nil {
			return nil, nil, err
		} else {
			result = append(result, formResult.(Symbol))
		}
		if i+1 < len(parts) { // if not the last string
			argKnot := Knot{{Name: "", Shape: args.Elem[i], Effect: nil, Frame: span}}
			if formResult, _, err := withArg.Evoke(span, argKnot); err != nil {
				return nil, nil, err
			} else {
				result = append(result, formResult.(Symbol))
			}
		}
	}
	if series, err := MakeSeriesSymbol(span, result); err != nil {
		return nil, nil, span.Errorf(err, "format unifying elements")
	} else {
		return series, nil, nil
	}
}
