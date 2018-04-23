package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Symbol implementations:
//	BasicSymbol, EmptySymbol,
// *SeriesSymbol, *StructSymbol
// *NamedSymbol, *OpaqueSymbol, *MapSymbol
// *VarietySymbol,
// *BlobSymbol
type Symbol interface {
	Shape   // String, Select, Augment, Invoke
	Splayer // Splay
	Type() Type
	Hash(*Span) string
	Equal(*Span, Symbol) bool
	LiftToSeries(*Span) *SeriesSymbol
	Disassemble(*Span) (*pb.Symbol, error)
}

var (
	symbolPtr    *Symbol
	typeOfSymbol = reflect.TypeOf((*Symbol)(nil)).Elem()
)

type Symbols []Symbol

func (syms Symbols) Types() Types {
	types := make(Types, len(syms))
	for i, sym := range syms {
		types[i] = sym.Type()
	}
	return types
}
