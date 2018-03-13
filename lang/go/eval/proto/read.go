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

type EvalReadProtoMacro struct {
	ProtoPkg  string       `ko:"name=protoPkg"`
	ProtoName string       `ko:"name=protoName"`
	MsgType   reflect.Type `ko:"name=msgType"` // MsgType is a proto message go struct
}

func (m *EvalReadProtoMacro) MacroID() string { return m.Help() }

func (m *EvalReadProtoMacro) Label() string { return "readproto" }

func (m *EvalReadProtoMacro) MacroSheathString() *string { return PtrString(m.Help()) }

func (m *EvalReadProtoMacro) Help() string {
	return fmt.Sprintf("ReadProto<%s.%s>", m.ProtoPkg, m.ProtoName)
}

var (
	someBytes   = []byte{}
	typeOfBytes = reflect.TypeOf(someBytes)
)

func (m *EvalReadProtoMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	bytesArg := a.Walk("bytes")
	bytesValue, err := Integrate(span, bytesArg, typeOfBytes)
	if err != nil {
		return nil, nil, span.Errorf(nil,
			"%s exepcts a bytes argument of type sequence of bytes, got %v",
			m.Help(), bytesArg,
		)
	}
	msg := reflect.New(m.MsgType).Interface().(proto.Message)
	if err := proto.Unmarshal(bytesValue.Interface().([]byte), msg); err != nil {
		panic(
			NewEvalPanic(
				span,
				MakeStructSymbol(
					FieldSymbols{{Name: "protoUnmarshal", Value: BasicStringSymbol(err.Error())}},
				),
			),
		)
	}
	result, err := Deconstruct(span, reflect.ValueOf(msg))
	if err != nil {
		panic("o")
		return nil, nil, span.Errorf(err, "deconstructing proto")
	}
	return result, nil, nil
}
