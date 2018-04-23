package proto

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"

	"github.com/golang/protobuf/proto"
)

type EvalMarshalProtoMacro struct {
	ProtoPkg  string       `ko:"name=protoPkg"`
	ProtoName string       `ko:"name=protoName"`
	MsgType   reflect.Type `ko:"name=msgType"` // MsgType is a proto message go struct
}

func (m *EvalMarshalProtoMacro) MacroID() string { return m.Help() }

func (m *EvalMarshalProtoMacro) Label() string { return "marshalProto" }

func (m *EvalMarshalProtoMacro) MacroSheathString() *string { return PtrString(m.Help()) }

func (m *EvalMarshalProtoMacro) Help() string {
	return fmt.Sprintf("MarshalProto<%s.%s>", m.ProtoPkg, m.ProtoName)
}

func (m *EvalMarshalProtoMacro) Doc() string {
	return m.Help()
}

func (m *EvalMarshalProtoMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	protoArg := a.Walk("proto")
	protoValue, err := Integrate(span, protoArg, reflect.PtrTo(m.MsgType)) // *MsgType
	if err != nil {
		return nil, nil, span.Errorf(nil,
			"%s exepcts a proto argument of type %v, got %v",
			m.Help(), m.MsgType, protoArg,
		)
	}
	buf, err := proto.Marshal(protoValue.Interface().(proto.Message))
	if err != nil {
		panic(
			NewEvalPanic(
				span,
				MakeStructSymbol(
					FieldSymbols{{Name: "marshalProto", Value: BasicStringSymbol(err.Error())}},
				),
			),
		)
	}
	return Deconstruct(span, reflect.ValueOf(buf)), nil, nil
}
