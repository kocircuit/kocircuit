package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"strings"

	// . "github.com/kocircuit/kocircuit/lang/go/kit/tree"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func RegisterEvalProtoFile(protoFile string) {
	if err := registerEvalProtoFile(protoFile); err != nil {
		log.Fatalf("registering proto file %q with runtime (%v)", protoFile, err)
	}
}

func registerEvalProtoFile(protoFile string) error {
	fileDesc, err := decodeFileDescriptor(protoFile)
	if err != nil {
		return err
	}
	protoPkg := fileDesc.GetPackage()
	// register proto macros for messages (make, read and write)
	for _, msgDesc := range fileDesc.MessageType {
		msgName := msgDesc.GetName()
		msgFullName := strings.Join([]string{protoPkg, msgName}, ".")
		msgType := proto.MessageType(msgFullName)
		if msgType == nil {
			return fmt.Errorf("cannot find message type for %q", msgFullName)
		}
		// RegisterEvalMacro(XXX) //XXX
		//XXX: read and write macros
	}
	// XXX: enums
	return nil
}

func decodeFileDescriptor(protoFile string) (*descriptor.FileDescriptorProto, error) {
	gzippedFileDescBuf := proto.FileDescriptor(protoFile)
	if len(gzippedFileDescBuf) == 0 {
		return nil, fmt.Errorf("no registration for %q with proto package", protoFile)
	}
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
