package macros

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Disassemble", new(EvalDisassembleMacro))
}

type EvalDisassembleMacro struct{}

func (m EvalDisassembleMacro) MacroID() string { return m.Help() }

func (m EvalDisassembleMacro) Label() string { return "disassemble" }

func (m EvalDisassembleMacro) MacroSheathString() *string { return PtrString("Disassemble") }

func (m EvalDisassembleMacro) Help() string { return "Disassemble" }

func (EvalDisassembleMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		if dis := returns.(Symbol).Disassemble(span); dis != nil {
			wrapped := reflect.ValueOf(dis).Convert(TypeOfInterface)
			return Deconstruct(span, wrapped), nil, nil
		} else {
			return Deconstruct(span, reflect.Zero(TypeOfInterface)), nil, nil
		}
	}
}
