package symbol

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
)

func EncodeSymbol(span *Span, symbol Symbol) ([]byte, error) {
	pbDisassembled, err := symbol.Disassemble(span)
	if err != nil {
		return nil, span.Errorf(err, "disassemble symbol")
	}
	buf, err := proto.Marshal(pbDisassembled)
	if err != nil {
		return nil, span.Errorf(err, "marshal proto")
	}
	var w bytes.Buffer
	gw := gzip.NewWriter(&w)
	if _, err := gw.Write(buf); err != nil {
		return nil, span.Errorf(err, "write gzipped bytes")
	}
	if err = gw.Close(); err != nil {
		return nil, span.Errorf(err, "flush gzipped bytes")
	}
	return w.Bytes(), nil
}

func DecodeSymbol(span *Span, asm VarietyAssembler, gzipped []byte) (Symbol, error) {
	r := bytes.NewBuffer(gzipped)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, span.Errorf(err, "un-gzipping symbol")
	}
	buf, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, span.Errorf(err, "reading un-gzipped symbol")
	}
	if err := gr.Close(); err != nil {
		return nil, span.Errorf(err, "closing un-gzipper")
	}
	pbSymbol := &pb.Symbol{}
	if err := proto.Unmarshal(buf, pbSymbol); err != nil {
		return nil, span.Errorf(err, "proto unmarshal")
	}
	return AssembleWithError(span, asm, pbSymbol)
}

func DecodeArg(span *Span, asm VarietyAssembler, argBytes []byte) (*StructSymbol, error) {
	if argBytes == nil {
		return MakeStructSymbol(nil), nil
	}
	sym, err := DecodeSymbol(span, asm, argBytes)
	if err != nil {
		return nil, span.Errorf(err, "decoding arg")
	}
	switch u := sym.(type) {
	case nil:
		return MakeStructSymbol(nil), nil
	case EmptySymbol:
		return MakeStructSymbol(nil), nil
	case *StructSymbol:
		return u, nil
	default:
		return nil, span.Errorf(nil, "arg must be structure or empty, got %v", u)
	}
}
