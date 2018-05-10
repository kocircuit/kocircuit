package symbol

import (
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type KeyValueSymbol struct {
	Key   string `ko:"name=key"`
	Value Symbol `ko:"name=value"`
}

func MakeMapSymbol(span *Span, m map[string]Symbol) (Symbol, error) {
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

func (ms *MapSymbol) Disassemble(span *Span) (*pb.Symbol, error) {
	filtered := filterMap(ms.Map)
	dis := &pb.SymbolMap{
		KeyValue: make([]*pb.SymbolKeyValue, 0, len(filtered)),
	}
	for _, key := range sortedMapKeys(filtered) {
		value, err := filtered[key].Disassemble(span)
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
	return Sprint(ms)
}

func (ms *MapSymbol) Equal(span *Span, sym Symbol) bool {
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

func (ms *MapSymbol) Hash(span *Span) ID {
	filtered := filterMap(ms.Map)
	h := make([]ID, 2*len(filtered))
	for i, key := range sortedMapKeys(filtered) {
		h[2*i] = StringID(key)
		h[2*i+1] = filtered[key].Hash(span)
	}
	return Blend(h...)
}

func sortedMapKeys(m map[string]Symbol) []string {
	kk := make([]string, 0, len(m))
	for k := range m {
		kk = append(kk, k)
	}
	sort.Strings(kk)
	return kk
}

func (ms *MapSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return singletonSeries(ms)
}

func (ms *MapSymbol) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to map")
}

func (ms *MapSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
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

func (ms *MapSymbol) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be augmented", ms)
}

func (ms *MapSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be invoked", ms)
}

func (ms *MapSymbol) Type() Type {
	return ms.Type_
}

func (ms *MapSymbol) SortedKeys() []string {
	return sortedMapKeys(ms.Map)
}

func (ms *MapSymbol) Splay() Tree {
	sortedKeys := ms.SortedKeys()
	nameTrees := make([]NameTree, len(sortedKeys))
	for i, key := range sortedKeys {
		nameTrees[i] = NameTree{
			Name:    gate.KoGoName{Ko: key},
			Monadic: false,
			Tree:    ms.Map[key].Splay(),
		}
	}
	return Parallel{
		Label:   Label{Path: "", Name: ""},
		Bracket: "{}",
		Elem:    nameTrees,
	}
}

type MapType struct {
	Value Type `ko:"name=type"`
}

func (mt *MapType) IsType() {}

func (mt *MapType) String() string {
	return Sprint(mt)
}

func (mt *MapType) Splay() Tree {
	return NoQuote{String_: fmt.Sprintf("Map<String:%v>", mt.Value)}
}
