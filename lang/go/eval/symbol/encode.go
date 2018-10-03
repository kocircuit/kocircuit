//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package symbol

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
)

func EncodeSymbol(span *model.Span, symbol Symbol) ([]byte, error) {
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

func DecodeSymbol(span *model.Span, asm VarietyAssembler, gzipped []byte) (Symbol, error) {
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

func DecodeArg(span *model.Span, asm VarietyAssembler, argBytes []byte) (*StructSymbol, error) {
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
