package proto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"

	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func RegisterEvalProtoFile(protoFile string) {
	if err := registerEvalProtoFile(protoFile); err != nil {
		log.Fatalf("registering proto file %q with runtime (%v)", protoFile, err)
	}
}

func registerEvalProtoFile(protoFile string) error {
	gzippedFileDescBuf := proto.FileDescriptor(protoFile)
	if len(gzippedFileDescBuf) == 0 {
		return fmt.Errorf("no registration for %q with proto package", protoFile)
	}
	fileDescBuf, err := UngzipBytes(gzippedFileDescBuf)
	if err != nil {
		return err
	}
	fileDesc := &descriptor.FileDescriptorProto{}
	if err := proto.Unmarshal(fileDescBuf, fileDesc); err != nil {
		return err
	}
	fmt.Println(Sprint(fileDesc))
	return nil
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
