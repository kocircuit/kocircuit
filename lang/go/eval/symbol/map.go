package symbol

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type MapSymbol struct {
	Value reflect.Value `ko:"name=value"` // go map value
}

func (ms *MapSymbol) GoType() reflect.Type {
	return ms.Value.Type()
}

func (ms *MapSymbol) Disassemble(span *Span) interface{} {
	return ms.Value.Interface()
}

func (ms *MapSymbol) String() string {
	return Sprint(ms)
}

func (ms *MapSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*MapSymbol); ok {
		return ms.Value.Interface() == other.Value.Interface()
	} else {
		return false
	}
}

func (ms *MapSymbol) Hash() string {
	return "█"
}

func (ms *MapSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return singletonSeries(ms)
}

func (ms *MapSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return ms, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "map value %v cannot be selected into", ms)
	}
}

func (ms *MapSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be augmented", ms)
}

func (ms *MapSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "map %v cannot be invoked", ms)
}

func (ms *MapSymbol) Type() Type {
	return &MapType{ms.Value.Type()}
}

func (ms *MapSymbol) Splay() Tree {
	return ms.Type().Splay()
}

type MapType struct {
	Type reflect.Type `ko:"name=type"`
}

func (mt *MapType) IsType() {}

func (mt *MapType) String() string {
	return Sprint(mt)
}

func (mt *MapType) Splay() Tree {
	return NoQuote{fmt.Sprintf("OpaqueMap<%v>", mt.Type)}
}
