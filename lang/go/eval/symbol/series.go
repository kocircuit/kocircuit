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
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type SeriesSymbol struct {
	Type_ *SeriesType `ko:"name=type"`
	Elem  Symbols     `ko:"name=elem"`
}

var _ Symbol = &SeriesSymbol{}

// DisassembleToGo converts a Ko value into a Go value
func (ss *SeriesSymbol) DisassembleToGo(span *model.Span) (reflect.Value, error) {
	filtered := FilterEmptySymbols(ss.Elem)
	seriesType := ss.Type_.GoType()
	slice := reflect.MakeSlice(seriesType, 0, len(filtered))
	for _, elem := range filtered {
		value, err := elem.DisassembleToGo(span)
		if err != nil {
			return reflect.Value{}, err
		}
		if !isNil(value) {
			slice = reflect.Append(slice, value)
		}
	}
	if slice.Len() == 0 {
		return reflect.Zero(seriesType), nil
	}
	return slice, nil
}

// DisassembleToPB converts a Ko value into a protobuf
func (ss *SeriesSymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	filtered := FilterEmptySymbols(ss.Elem)
	dis := &pb.SymbolSeries{
		Element: make([]*pb.Symbol, 0, len(filtered)),
	}
	for _, elem := range filtered {
		value, err := elem.DisassembleToPB(span)
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
	return tree.Sprint(ss)
}

func (ss *SeriesSymbol) Equal(span *model.Span, sym Symbol) bool {
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

func (ss *SeriesSymbol) Hash(span *model.Span) model.ID {
	h := make([]model.ID, 0, len(ss.Elem))
	for _, e := range ss.Elem {
		h = append(h, e.Hash(span))
	}
	switch len(h) {
	case 0:
		return EmptySymbol{}.Hash(span)
	case 1:
		return h[0]
	default:
		return model.Blend(h...)
	}
}

func (ss *SeriesSymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return ss
}

func (ss *SeriesSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to series")
}

func (ss *SeriesSymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "cannot select %v into sequence %v", path, ss)
	}
}

func (ss *SeriesSymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "sequence %v cannot be augmented", ss)
}

func (ss *SeriesSymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "sequence %v cannot be invoked", ss)
}

func (ss *SeriesSymbol) Type() Type {
	return ss.Type_
}

func (ss *SeriesSymbol) Splay() tree.Tree {
	indexTrees := make([]tree.IndexTree, len(ss.Elem))
	for i, elemSym := range ss.Elem {
		indexTrees[i] = tree.IndexTree{Index: i, Tree: elemSym.Splay()}
	}
	return tree.Series{
		Label:   tree.Label{Path: "", Name: ""},
		Bracket: "()",
		Elem:    indexTrees,
	}
}

type SeriesType struct {
	Elem Type `ko:"name=elem"`
}

var _ Type = &SeriesType{}

func (*SeriesType) IsType() {}

func (st *SeriesType) String() string {
	return tree.Sprint(st)
}

func (st *SeriesType) Splay() tree.Tree {
	return tree.Series{
		Label:   tree.Label{Path: "", Name: ""},
		Bracket: "()",
		Elem:    []tree.IndexTree{{Index: 0, Tree: st.Elem.Splay()}},
	}
}

// GoType returns the Go equivalent of the type.
func (st *SeriesType) GoType() reflect.Type {
	return reflect.SliceOf(st.Elem.GoType())
}

func MakeStringsSymbol(span *model.Span, ss []string) Symbol {
	switch len(ss) {
	case 0:
		return EmptySymbol{}
	case 1:
		return BasicSymbol{ss[0]}
	default:
		elems := make(Symbols, len(ss))
		for i := range ss {
			elems[i] = BasicSymbol{ss[i]}
		}
		return makeSeriesDontUnify(span, elems, elems[0].Type())
	}
}

func MakeSeriesSymbol(span *model.Span, elem Symbols) (Symbol, error) {
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

func makeSeriesDontUnify(span *model.Span, elem Symbols, elemType Type) *SeriesSymbol {
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
