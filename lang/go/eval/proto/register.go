package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/go/eval"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func RegisterEvalProtoFile(protoFile string) {
	if err := registerEvalProtoFile(protoFile); err != nil {
		log.Fatalf("registering proto file %q with runtime (%v)", protoFile, err)
	}
}

func registerEvalProtoFile(protoFile string) error {
	fileDesc, err := decodeRegisteredProtoFile(protoFile)
	if err != nil {
		return err
	}
	return registerEvalFileDescriptor(fileDesc)
}

func RegisterEvalFileDescriptor(fileDesc *descriptor.FileDescriptorProto) {
	if err := registerEvalFileDescriptor(fileDesc); err != nil {
		log.Fatalf("registering file descriptor with runtime (%v)", err)
	}
}

func registerEvalFileDescriptor(fileDesc *descriptor.FileDescriptorProto) error {
	ctx := &registerCtx{
		ProtoPkg:   fileDesc.GetPackage(),
		KoProtoPkg: path.Join("proto", fileDesc.GetPackage()),
	}
	for _, msgDesc := range fileDesc.MessageType {
		if err := ctx.registerMessage(msgDesc); err != nil {
			return err
		}
	}
	for _, enumDesc := range fileDesc.EnumType {
		if err := ctx.registerEnum(enumDesc); err != nil {
			return err
		}
	}
	return nil
}

type registerCtx struct {
	ProtoPkg   string   `ko:"name=protoPkg"`
	KoProtoPkg string   `ko:"name=koProtoPkg"`
	Namespace  []string `ko:"name=namespace"` // message namespace, within package
}

func (ctx *registerCtx) KoName(name string) string {
	if prefix := strings.Join(ctx.Namespace, "_"); prefix == "" {
		return name
	} else {
		return strings.Join([]string{prefix, name}, "_")
	}
}

func (ctx *registerCtx) GoName(name string) string {
	if prefix := strings.Join(ctx.Namespace, "."); prefix == "" {
		return name
	} else {
		return strings.Join([]string{prefix, name}, ".")
	}
}

func (ctx *registerCtx) registerEnum(enumDesc *descriptor.EnumDescriptorProto) error {
	enumName := ctx.KoName(enumDesc.GetName())
	for _, valueDesc := range enumDesc.Value {
		valueName := valueDesc.GetName()
		RegisterEvalPkgMacro(
			ctx.KoProtoPkg,
			fmt.Sprintf("Enum%s_%s", enumName, valueName), // Enum<EnumName>_<ValueName>
			&EvalProtoEnumValueMacro{
				ProtoPkg:  ctx.ProtoPkg,
				EnumName:  enumName,
				ValueName: valueName,
				Number:    valueDesc.GetNumber(),
			},
		)
	}
	return nil
}

func (ctx *registerCtx) registerMessage(msgDesc *descriptor.DescriptorProto) error {
	msgName := ctx.GoName(msgDesc.GetName())
	msgFullName := strings.Join([]string{ctx.ProtoPkg, msgName}, ".")
	msgType := proto.MessageType(msgFullName)
	if msgType == nil {
		return fmt.Errorf("cannot find message type for %q with descriptor %v", msgFullName, Sprint(msgDesc))
	}
	nestedCtx := &registerCtx{
		ProtoPkg:   ctx.ProtoPkg,
		KoProtoPkg: ctx.KoProtoPkg,
		Namespace:  append(append([]string{}, ctx.Namespace...), msgDesc.GetName()),
	}
	for _, nestedMsgDesc := range msgDesc.NestedType {
		if !nestedMsgDesc.GetOptions().GetMapEntry() {
			if err := nestedCtx.registerMessage(nestedMsgDesc); err != nil {
				return err
			}
		}
	}
	for _, nestedEnumDesc := range msgDesc.EnumType {
		if err := nestedCtx.registerEnum(nestedEnumDesc); err != nil {
			return err
		}
	}
	RegisterEvalPkgMacro( // Proto<MsgName>
		ctx.KoProtoPkg,
		fmt.Sprintf("Proto%s", msgName),
		&EvalProtoMessageMacro{
			ProtoPkg:  ctx.ProtoPkg,
			ProtoName: msgName,
			MsgType:   msgType,
		},
	)
	RegisterEvalPkgMacro(
		ctx.KoProtoPkg,
		fmt.Sprintf("Unmarshal%s", msgName), // Unmarshal<MsgName>
		&EvalUnmarshalProtoMacro{
			ProtoPkg:  ctx.ProtoPkg,
			ProtoName: msgName,
			MsgType:   msgType,
		},
	)
	RegisterEvalPkgMacro(
		ctx.KoProtoPkg,
		fmt.Sprintf("Marshal%s", msgName), // Marshal<MsgName>
		&EvalMarshalProtoMacro{
			ProtoPkg:  ctx.ProtoPkg,
			ProtoName: msgName,
			MsgType:   msgType,
		},
	)
	return nil
}

func decodeRegisteredProtoFile(protoFile string) (*descriptor.FileDescriptorProto, error) {
	gzippedFileDescBuf := proto.FileDescriptor(protoFile)
	if len(gzippedFileDescBuf) == 0 {
		return nil, fmt.Errorf("no registration for %q with proto package", protoFile)
	}
	return decodeFileDescriptorBytes(gzippedFileDescBuf)
}

func decodeFileDescriptorBytes(gzippedFileDescBuf []byte) (*descriptor.FileDescriptorProto, error) {
	fileDescBuf, err := UngzipBytes(gzippedFileDescBuf)
	if err != nil {
		return nil, err
	}
	fileDesc := &descriptor.FileDescriptorProto{}
	if err := proto.Unmarshal(fileDescBuf, fileDesc); err != nil {
		return nil, err
	}
	return fileDesc, nil
}

func UngzipBytes(gzipped []byte) ([]byte, error) {
	var w bytes.Buffer
	dec, err := gzip.NewReader(bytes.NewBuffer(gzipped))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(&w, dec); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
