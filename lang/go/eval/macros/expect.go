package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Expect", EvalExpectMacro{})
	RegisterEvalMacro("Require", EvalRequireMacro{})
}

type EvalExpectMacro struct{}

func (m EvalExpectMacro) MacroID() string { return m.Help() }

func (m EvalExpectMacro) Label() string { return "expect" }

func (m EvalExpectMacro) MacroSheathString() *string { return PtrString("Expect") }

func (m EvalExpectMacro) Help() string { return "Expect" }

func (EvalExpectMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol).SelectMonadic()
	if IsEmptySymbol(a) {
		panic(
			MakeStructSymbol(
				FieldSymbols{
					{Name: "expected", Value: BasicTrue}, // indicate panic originates from Expect macro
				},
			),
		)
	} else {
		return a, nil, nil
	}
}

type EvalRequireMacro struct{}

func (m EvalRequireMacro) MacroID() string { return m.Help() }

func (m EvalRequireMacro) Label() string { return "require" }

func (m EvalRequireMacro) MacroSheathString() *string { return PtrString("Require") }

func (m EvalRequireMacro) Help() string { return "Require" }

func (EvalRequireMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol).SelectMonadic()
	if IsEmptySymbol(a) {
		return nil, nil, span.Errorf(nil, "required value missing")
	} else {
		return a, nil, nil
	}
}
