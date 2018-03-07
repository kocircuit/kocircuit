package macros

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type EvalConstantMacro struct {
	Value interface{} `ko:"name=value"`
}

func (m EvalConstantMacro) MacroID() string { return m.Help() }

func (m EvalConstantMacro) Label() string { return "constant" }

func (m EvalConstantMacro) MacroSheathString() *string { return PtrString(m.Help()) }

func (m EvalConstantMacro) Help() string {
	return fmt.Sprintf("Constant(%s)", Sprint(m.Value))
}

func (m EvalConstantMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if sym, err := Deconstruct(span, reflect.ValueOf(m.Value)); err != nil {
		return nil, nil, err
	} else {
		return sym, nil, nil
	}
}
