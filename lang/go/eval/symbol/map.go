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
	"sort"

	"github.com/golang/protobuf/proto"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type KeyValueSymbol struct {
	Key   string `ko:"name=key"`
	Value Symbol `ko:"name=value"`
}

func MakeMapSymbol(span *model.Span, m map[string]Symbol) (Symbol, error) {
	if len(m) == 0 {
		return EmptySymbol{}, nil
	}
	vtypes := make([]Type, 0, len(m))
	for _, s := range m {
		vtypes = append(vtypes, s.Type())
	}
	if unified, err := UnifyTypes(span, vtypes); err != nil {
		return nil, err
	} else {
		return &MapSymbol{
			Type_: &MapType{Value: unified},
			Map:   m,
		}, nil
	}
}

// MapSymbol captures map[string]Q types.
type MapSymbol struct {
	Type_ *MapType          `ko:"name=type"`
	Map   map[string]Symbol `ko:"name=map"`
}

var	_         Symbol = &MapSymbol{}

// DisassembleToGo converts a Ko value into a Go value
func (ms *MapSymbol) DisassembleToGo(span *model.Span) (reflect.Value, error) {
	filtered := filterMap(ms.Map)
	mapType := ms.Type_.GoType()
	m := reflect.MakeMap(mapType)
	for _, key := range sortedMapKeys(filtered) {
		value, err := filtered[key].DisassembleToGo(span)
		if err != nil {
			return reflect.Value{}, err
		}
		if !isNil(value) {
			m.SetMapIndex(reflect.ValueOf(key), value)
		}
	}
	if m.Len() == 0 {
		return reflect.Zero(mapType), nil
	}
	return reflect.ValueOf(m), nil
}

// DisassembleToPB converts a Ko value into a protobuf
func (ms *MapSymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	filtered := filterMap(ms.Map)
	dis := &pb.SymbolMap{
		KeyValue: make([]*pb.SymbolKeyValue, 0, len(filtered)),
	}
	for _, key := range sortedMapKeys(filtered) {
		value, err := filtered[key].DisassembleToPB(span)
		if err != nil {
			return nil, err
		}
		if value != nil {
			dis.KeyValue = append(dis.KeyValue,
				&pb.SymbolKeyValue{
					Key:   proto.String(key),
					Value: value,
				},
			)
		}
	}
	if len(dis.KeyValue) == 0 {
		return nil, nil
	} else {
		return &pb.Symbol{
			Symbol: &pb.Symbol_Map{Map: dis},
		}, nil
	}
}

func (ms *MapSymbol) String() string {
	return tree.Sprint(ms)
}

func (ms *MapSymbol) Equal(span *model.Span, sym Symbol) bool {
	if other, ok := sym.(*MapSymbol); ok {
		filteredThis, filtedOther := filterMap(ms.Map), filterMap(other.Map)
		if len(filteredThis) == len(filtedOther) {
			for k := range filteredThis {
				if !filteredThis[k].Equal(span, filtedOther[k]) {
					return false
				}
			}
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func filterMap(m map[string]Symbol) (filtered map[string]Symbol) {
	filtered = map[string]Symbol{}
	for k, v := range m {
		if !IsEmptySymbol(v) {
			filtered[k] = v
		}
	}
	return filtered
}

func (ms *MapSymbol) Hash(span *model.Span) model.ID {
	filtered := filterMap(ms.Map)
	h := make([]model.ID, 2*len(filtered))
	for i, key := range sortedMapKeys(filtered) {
		h[2*i] = model.StringID(key)
		h[2*i+1] = filtered[key].Hash(span)
	}
	return model.Blend(h...)
}

func sortedMapKeys(m map[string]Symbol) []string {
	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	return kk
}

func (ms *MapSymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return singletonSeries(ms)
}

func (ms *MapSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to map")
}

func (ms *MapSymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return ms, nil, nil
	} else {
		if v, ok := ms.Map[path[0]]; ok {
			return v.Select(span, path[1:])
		} else {
			return EmptySymbol{}, nil, nil
		}
	}
}

func (ms *MapSymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be augmented", ms)
}

func (ms *MapSymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be invoked", ms)
}

func (ms *MapSymbol) Type() Type {
	return ms.Type_
}

func (ms *MapSymbol) SortedKeys() []string {
	return sortedMapKeys(ms.Map)
}

func (ms *MapSymbol) Splay() tree.Tree {
	sortedKeys := ms.SortedKeys()
	nameTrees := make([]tree.NameTree, len(sortedKeys))
	for i, key := range sortedKeys {
		nameTrees[i] = tree.NameTree{
			Name:    gate.KoGoName{Ko: key},
			Monadic: false,
			Tree:    ms.Map[key].Splay(),
		}
	}
	return tree.Parallel{
		Label:   tree.Label{Path: "", Name: ""},
		Bracket: "{}",
		Elem:    nameTrees,
	}
}

// MapType is the type capturing map[string]T
type MapType struct {
	Value Type `ko:"name=type"`
}

var _ Type = &MapType{}

func (mt *MapType) IsType() {}

func (mt *MapType) String() string {
	return tree.Sprint(mt)
}

func (mt *MapType) Splay() tree.Tree {
	return tree.NoQuote{String_: fmt.Sprintf("Map<String:%v>", mt.Value)}
}

// GoType returns the Go equivalent of the type.
func (mt *MapType) GoType() reflect.Type {
	return reflect.MapOf(goString, mt.Value.GoType())
}
