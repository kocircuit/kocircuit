package symbol

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type OpaqueSymbol struct {
	Value reflect.Value `ko:"name=value"`
}

func (opaque *OpaqueSymbol) Interface() interface{} {
	return opaque.Value.Interface()
}

func (opaque *OpaqueSymbol) Disassemble(span *Span) interface{} {
	if opaque.Value.Kind() == reflect.Interface {
		return MustDeconstruct(span, opaque.Value.Elem()).Disassemble(span)
	} else {
		return opaque.Interface()
	}
}

func (opaque *OpaqueSymbol) String() string {
	return Sprint(opaque)
}

func (opaque *OpaqueSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*OpaqueSymbol); ok {
		return opaque.Value.Interface() == other.Value.Interface()
	} else {
		return false
	}
}

func (opaque *OpaqueSymbol) Hash() string {
	return "█"
}

func (opaque *OpaqueSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return singletonSeries(opaque)
}

func (opaque *OpaqueSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return opaque, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "opaque value %v cannot be selected into", opaque)
	}
}

func (opaque *OpaqueSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "opaque value %v cannot be augmented", opaque)
}

func (opaque *OpaqueSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "opaque value %v cannot be invoked", opaque)
}

func (opaque *OpaqueSymbol) GoType() reflect.Type {
	return opaque.Value.Type()
}

func (opaque *OpaqueSymbol) Type() Type {
	return &OpaqueType{Type: opaque.Value.Type()}
}

func (opaque *OpaqueSymbol) Splay() Tree {
	return opaque.Type().Splay()
}

type OpaqueType struct {
	Type reflect.Type `ko:"name=type"`
}

func (opaque *OpaqueType) IsType() {}

func (opaque *OpaqueType) String() string {
	return Sprint(opaque)
}

func (opaque *OpaqueType) Splay() Tree {
	return NoQuote{fmt.Sprintf("Opaque<%v>", opaque.Type)}
}
