package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type SeriesSymbol struct {
	Type_ *SeriesType `ko:"name=type"`
	Elem  Symbols     `ko:"name=elem"`
}

func (ss *SeriesSymbol) Disassemble(span *Span) (*pb.Symbol, error) {
	filtered := FilterEmptySymbols(ss.Elem)
	dis := &pb.SymbolSeries{
		Element: make([]*pb.Symbol, 0, len(filtered)),
	}
	for _, elem := range filtered {
		value, err := elem.Disassemble(span)
		if err != nil {
			return nil, err
		}
		if value != nil {
			dis.Element = append(dis.Element, value)
		}
	}
	if len(dis.Element) == 0 {
		return nil, nil
	} else {
		return &pb.Symbol{
			Symbol: &pb.Symbol_Series{Series: dis},
		}, nil
	}
}

func (ss *SeriesSymbol) IsEmpty() bool {
	return ss.Len() == 0
}

func (ss *SeriesSymbol) Len() int {
	return len(ss.Elem)
}

func (ss *SeriesSymbol) String() string {
	return Sprint(ss)
}

func (ss *SeriesSymbol) Equal(span *Span, sym Symbol) bool {
	other := sym.LiftToSeries(span)
	if len(ss.Elem) != len(other.Elem) {
		return false
	}
	for i := range ss.Elem {
		if !ss.Elem[i].Equal(span, other.Elem[i]) {
			return false
		}
	}
	return true
}

func (ss *SeriesSymbol) Hash(span *Span) ID {
	h := make([]ID, 0, len(ss.Elem))
	for _, e := range ss.Elem {
		h = append(h, e.Hash(span))
	}
	switch len(h) {
	case 0:
		return EmptySymbol{}.Hash(span)
	case 1:
		return h[0]
	default:
		return Blend(h...)
	}
}

func (ss *SeriesSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return ss
}

func (ss *SeriesSymbol) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to series")
}

func (ss *SeriesSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "cannot select %v into sequence %v", path, ss)
	}
}

func (ss *SeriesSymbol) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "sequence %v cannot be augmented", ss)
}

func (ss *SeriesSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "sequence %v cannot be invoked", ss)
}

func (ss *SeriesSymbol) Type() Type {
	return ss.Type_
}

func (ss *SeriesSymbol) Splay() Tree {
	indexTrees := make([]IndexTree, len(ss.Elem))
	for i, elemSym := range ss.Elem {
		indexTrees[i] = IndexTree{Index: i, Tree: elemSym.Splay()}
	}
	return Series{
		Label:   Label{Path: "", Name: ""},
		Bracket: "()",
		Elem:    indexTrees,
	}
}

type SeriesType struct {
	Elem Type `ko:"name=elem"`
}

func (*SeriesType) IsType() {}

func (st *SeriesType) String() string {
	return Sprint(st)
}

func (st *SeriesType) Splay() Tree {
	return Series{
		Label:   Label{Path: "", Name: ""},
		Bracket: "()",
		Elem:    []IndexTree{{Index: 0, Tree: st.Elem.Splay()}},
	}
}

func MustMakeSeriesSymbol(span *Span, elem Symbols) Symbol {
	if ss, err := MakeSeriesSymbol(span, elem); err != nil {
		return ss
	} else {
		panic("must")
	}
}

func MakeSeriesSymbol(span *Span, elem Symbols) (Symbol, error) {
	elem = FilterEmptySymbols(elem)
	if len(elem) == 0 {
		return EmptySymbol{}, nil
	}
	if unified, err := UnifyTypes(span, elem.Types()); err != nil {
		return nil, err
	} else {
		return &SeriesSymbol{
			Type_: &SeriesType{Elem: unified},
			Elem:  elem,
		}, nil
	}
}

var EmptySeries = &SeriesSymbol{
	Type_: &SeriesType{Elem: EmptyType{}},
}

func makeSeriesDontUnify(span *Span, elem Symbols, elemType Type) *SeriesSymbol {
	return &SeriesSymbol{
		Type_: &SeriesType{Elem: elemType},
		Elem:  elem,
	}
}

func singletonSeries(e Symbol) *SeriesSymbol {
	return &SeriesSymbol{
		Type_: &SeriesType{Elem: e.Type()},
		Elem:  Symbols{e},
	}
}

func FilterEmptySymbols(symbols Symbols) (filtered Symbols) {
	filtered = make(Symbols, 0, len(symbols))
	for _, sym := range symbols {
		if !IsEmptySymbol(sym) {
			filtered = append(filtered, sym)
		}
	}
	return
}
