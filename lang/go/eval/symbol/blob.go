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
	"fmt"
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func MakeBlobSymbol(b []byte) *BlobSymbol {
	return &BlobSymbol{Value: reflect.ValueOf(b)}
}

type BlobSymbol struct {
	Value reflect.Value `ko:"name=value"` // []byte
}

var _ Symbol = &BlobSymbol{}

// DisassembleToGo converts a Ko value into a Go value
func (blob *BlobSymbol) DisassembleToGo(span *model.Span) (reflect.Value, error) {
	return reflect.ValueOf(blob.Bytes()), nil
}

// DisassembleToPB converts a Ko value into a protobuf
func (blob *BlobSymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	return &pb.Symbol{
		Symbol: &pb.Symbol_Blob{
			Blob: &pb.SymbolBlob{Bytes: blob.Bytes()},
		},
	}, nil
}

func (blob *BlobSymbol) Bytes() []byte {
	return blob.Value.Bytes()
}

func (blob *BlobSymbol) String() string {
	return fmt.Sprintf("Blob<%d>", blob.Value.Len())
}

func (blob *BlobSymbol) Equal(span *model.Span, sym Symbol) bool {
	if other, ok := sym.(*BlobSymbol); ok {
		left, right := blob.Bytes(), other.Bytes()
		if len(left) != len(right) {
			return false
		} else {
			for i := range left {
				if left[i] != right[i] {
					return false
				}
			}
			return true
		}
	} else {
		return false
	}
}

func (blob *BlobSymbol) Splay() tree.Tree {
	return tree.NoQuote{String_: blob.String()}
}

func (blob *BlobSymbol) Hash(span *model.Span) model.ID {
	return model.BytesID(blob.Bytes())
}

func (blob *BlobSymbol) Type() Type {
	return BlobType{}
}

func (blob *BlobSymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	goBytes := blob.Bytes()
	if len(goBytes) == 0 {
		return EmptySeries
	} else {
		elemSymbols := make(Symbols, len(goBytes))
		for i, b := range goBytes {
			elemSymbols[i] = BasicByteSymbol(b)
		}
		return makeSeriesDontUnify(span, elemSymbols, BasicInt8)
	}
}

func (blob *BlobSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to blob")
}

func (blob *BlobSymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return blob, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "blob value %v cannot be selected into", blob)
	}
}

func (blob *BlobSymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "blob value %v cannot be augmented", blob)
}

func (blob *BlobSymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "blob value %v cannot be invoked", blob)
}

// BlobType captures the []byte type.
type BlobType struct{}

var (
	_          Type = BlobType{}
	goBlobType      = reflect.TypeOf((*[]byte)(nil)).Elem()
)

func (blob BlobType) IsType() {}

func (blob BlobType) String() string {
	return "Blob"
}

func (blob BlobType) Splay() tree.Tree {
	return tree.NoQuote{String_: "Blob"}
}

// GoType returns the Go equivalent of the type.
func (BlobType) GoType() reflect.Type {
	return goBlobType
}
