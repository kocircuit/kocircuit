package symbol

import (
	"fmt"
	"sort"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// MapSymbol captures map[string]Q types.
type MapSymbol struct {
	Type_ *MapType          `ko:"name=type"`
	Map   map[string]Symbol `ko:"name=map"`
}

func (ms *MapSymbol) Disassemble(span *Span) interface{} {
	dis := map[string]interface{}{}
	for k, v := range filterMap(ms.Map) {
		dis[k] = v.Disassemble(span)
	}
	return dis
}

func (ms *MapSymbol) String() string {
	return Sprint(ms)
}

func (ms *MapSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*MapSymbol); ok {
		filteredThis, filtedOther := filterMap(ms.Map), filterMap(other.Map)
		if len(filteredThis) == len(filtedOther) {
			for k := range filteredThis {
				if !filteredThis[k].Equal(filtedOther[k]) {
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

func (ms *MapSymbol) Hash() string {
	filtered := filterMap(ms.Map)
	h := make([]string, 2*len(filtered))
	for i, key := range sortedMapKeys(filtered) {
		h[2*i] = key
		h[2*i+1] = filtered[key].Hash()
	}
	return Mix(h...)
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

func (ms *MapSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
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
		Label:   Label{"", ""},
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
	return NoQuote{fmt.Sprintf("Map<String:%v>", mt.Value)}
}
