package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func IsEmptySymbol(sym Symbol) bool {
	_, isEmpty := sym.(EmptySymbol)
	return isEmpty
}

func IsEmptyType(t Type) bool {
	_, isEmpty := t.(EmptyType)
	return isEmpty
}

type EmptySymbol struct{}

func (empty EmptySymbol) Disassemble(span *Span) (*pb.Symbol, error) {
	return nil, nil
}

func (empty EmptySymbol) String() string {
	return Sprint(empty)
}

func (empty EmptySymbol) Equal(span *Span, sym Symbol) bool {
	_, ok := sym.(EmptySymbol)
	return ok
}

func (empty EmptySymbol) Hash() string {
	return "" //XXX
}

func (empty EmptySymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return EmptySeries
}

func (empty EmptySymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	return empty, nil, nil
}

func (empty EmptySymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "empty value cannot be augmented")
}

func (empty EmptySymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "empty value cannot be invoked")
}

func (empty EmptySymbol) Type() Type {
	return EmptyType{}
}

func (empty EmptySymbol) Splay() Tree {
	return NoQuote{"empty"}
}

type EmptyType struct{}

func (EmptyType) IsType() {}

func (EmptyType) String() string {
	return Sprint(EmptyType{})
}

func (EmptyType) Splay() Tree {
	return NoQuote{"Empty"}
}
