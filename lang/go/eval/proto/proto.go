package proto

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type EvalProtoMessageMacro struct {
	ProtoPkg  string       `ko:"name=protoPkg"`
	ProtoName string       `ko:"name=protoName"`
	MsgType   reflect.Type `ko:"name=msgType"` // MsgType is a proto message go struct
}

func (m *EvalProtoMessageMacro) MacroID() string { return m.Help() }

func (m *EvalProtoMessageMacro) Label() string { return "protoMessage" }

func (m *EvalProtoMessageMacro) MacroSheathString() *string { return PtrString(m.Help()) }

func (m *EvalProtoMessageMacro) Help() string {
	return fmt.Sprintf("Proto<%s.%s>", m.ProtoPkg, m.ProtoName)
}

func (m *EvalProtoMessageMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	msgValue, err := Integrate(span, arg.(*StructSymbol), m.MsgType)
	if err != nil {
		return nil, nil, err
	}
	msgSym, err := Deconstruct(span, msgValue)
	if err != nil {
		return nil, nil, err
	}
	return msgSym, nil, nil
}
