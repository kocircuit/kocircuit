package symbol

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type NamedSymbol struct {
	Value reflect.Value `ko:"name=value"`
}

func (named *NamedSymbol) String() string {
	return Sprint(named.Value.Interface())
}

func (named *NamedSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*NamedSymbol); ok {
		return reflect.DeepEqual(named.Value.Interface(), other.Value.Interface())
	} else {
		return false
	}
}

func (named *NamedSymbol) Splay() Tree {
	return Splay(named.Value.Interface())
}

func (named *NamedSymbol) Hash() string {
	return named.Splay().TreeHash()
}

func (named *NamedSymbol) GoType() reflect.Type {
	return named.Value.Type()
}

func (named *NamedSymbol) Type() Type {
	return NamedType{
		Type: named.GoType(),
	}
}

func (named *NamedSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return named, nil, nil
	} else {
		if step, err := named.Walk(span, path[0]); err != nil {
			return nil, nil, err
		} else {
			return step.Select(span, path[1:])
		}
	}
}

func (named *NamedSymbol) Walk(span *Span, field string) (Symbol, error) {
	v := named.Value
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return EmptySymbol{}, nil
	}
	if v.Kind() != reflect.Struct {
		return nil, span.Errorf(nil, "cannot select %s into %v", field, named)
	}
	if fieldIndex, ok := gate.StripFields(v.Type()).FieldByKoName(field); ok {
		return Deconstruct(span, v.Field(fieldIndex))
	} else {
		return EmptySymbol{}, nil
	}
}

func (named *NamedSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "named value %v cannot be augmented", named)
}

func (named *NamedSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "named value %v cannot be invoked", named)
}

type NamedType struct {
	Type reflect.Type `ko:"name=type"`
}

func (named NamedType) IsType() {}

func (named NamedType) String() string {
	return Sprint(named)
}

func (named NamedType) Splay() Tree {
	return NoQuote{fmt.Sprintf("Named<%s.%s>", named.Type.PkgPath(), named.Type.Name())}
}
