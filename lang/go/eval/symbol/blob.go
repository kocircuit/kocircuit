package symbol

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func MakeBlobSymbol(b []byte) *BlobSymbol {
	return &BlobSymbol{Value: reflect.ValueOf(b)}
}

type BlobSymbol struct {
	Value reflect.Value `ko:"name=value"` // []byte
}

func (blob *BlobSymbol) Disassemble(span *Span) (*pb.Symbol, error) {
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

func (blob *BlobSymbol) Equal(span *Span, sym Symbol) bool {
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

func (blob *BlobSymbol) Splay() Tree {
	return NoQuote{blob.String()}
}

func (blob *BlobSymbol) Hash(span *Span) ID {
	return BytesID(blob.Bytes())
}

func (blob *BlobSymbol) Type() Type {
	return BlobType{}
}

func (blob *BlobSymbol) LiftToSeries(span *Span) *SeriesSymbol {
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

func (blob *BlobSymbol) SelectArg(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "selecting argument into blob")
}

func (blob *BlobSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return blob, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "blob value %v cannot be selected into", blob)
	}
}

func (blob *BlobSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "blob value %v cannot be augmented", blob)
}

func (blob *BlobSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "blob value %v cannot be invoked", blob)
}

type BlobType struct{}

func (blob BlobType) IsType() {}

func (blob BlobType) String() string {
	return "Blob"
}

func (blob BlobType) Splay() Tree {
	return NoQuote{"Blob"}
}
