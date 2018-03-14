package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type SeriesSymbol struct {
	Type_ *SeriesType `ko:"name=type"`
	Elem  Symbols     `ko:"name=elem"`
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

func (ss *SeriesSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*SeriesSymbol); ok {
		if len(ss.Elem) != len(other.Elem) {
			return false
		}
		for i := range ss.Elem {
			if !ss.Elem[i].Equal(other.Elem[i]) {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

func (ss *SeriesSymbol) Hash() string {
	h := make([]string, len(ss.Elem))
	for i, e := range ss.Elem {
		h[i] = e.Hash()
	}
	return Mix(h...)
}

func singletonSeries(e Symbol) *SeriesSymbol {
	return &SeriesSymbol{
		Type_: &SeriesType{Elem: e.Type()},
		Elem:  Symbols{e},
	}
}

func (ss *SeriesSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return ss
}

func (ss *SeriesSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "cannot select %v into sequence %v", path, ss)
	}
}

func (ss *SeriesSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
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

func MakeSeriesSymbol(span *Span, elem Symbols) (*SeriesSymbol, error) {
	elem = FilterEmptySymbols(elem)
	if unified, err := UnifyTypes(span, elem.Types()); err != nil {
		return nil, err
	} else {
		return &SeriesSymbol{
			Type_: &SeriesType{Elem: unified},
			Elem:  elem,
		}, nil
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
