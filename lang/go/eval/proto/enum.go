package proto

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type EvalProtoEnumValueMacro struct {
	ProtoPkg  string `ko:"name=protoPkg"`
	EnumName  string `ko:"name=enumName"`
	ValueName string `ko:"name=valueName"`
	Number    int32  `ko:"name=number"`
}

func (m *EvalProtoEnumValueMacro) MacroID() string { return m.Help() }

func (m *EvalProtoEnumValueMacro) Label() string { return "enumValue" }

func (m *EvalProtoEnumValueMacro) MacroSheathString() *string { return PtrString(m.Help()) }

func (m *EvalProtoEnumValueMacro) Help() string {
	return fmt.Sprintf("ProtoEnumValue<%s.%s.%s>", m.ProtoPkg, m.EnumName, m.ValueName)
}

func (m *EvalProtoEnumValueMacro) Doc() string {
	return m.Help()
}

func (m *EvalProtoEnumValueMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return BasicInt32Symbol(m.Number), nil, nil
}
