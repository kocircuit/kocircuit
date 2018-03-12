package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func Integrate(span *Span, s Symbol, t reflect.Type) (reflect.Value, error) {
	ctx := &typingCtx{Span: span}
	return ctx.Integrate(s, t)
}

func (ctx *typingCtx) Integrate(s Symbol, t reflect.Type) (reflect.Value, error) {
	if typeName := TypeName(t); typeName != "" {
		return ctx.IntegrateNamed(s, t)
	}
	switch t.Kind() {
	case reflect.Invalid:
		panic("o")
	case reflect.String, reflect.Bool:
		if IsBasicKind(s, t.Kind()) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Int: // int(go) can be assigned from int64(ko)
		if IsBasicKind(s, reflect.Int64) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if IsBasicKind(s, t.Kind()) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Uint: // uint(go) can be assigned from uint64(ko)
		if IsBasicKind(s, reflect.Uint64) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if IsBasicKind(s, t.Kind()) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Float32, reflect.Float64:
		if IsBasicKind(s, t.Kind()) {
			return reflect.ValueOf(s.(BasicSymbol).Value).Convert(t), nil
		}
	case reflect.Uintptr:
		return reflect.Value{}, ctx.Errorf(nil, "go uintptr type cannot be integrated")
	case reflect.Complex64:
		return reflect.Value{}, ctx.Errorf(nil, "go complex64 type cannot be integrated")
	case reflect.Complex128:
		return reflect.Value{}, ctx.Errorf(nil, "go complex128 type cannot be integrated")
	case reflect.Array:
		return reflect.Value{}, ctx.Errorf(nil, "go array type cannot be integrated")
	case reflect.Chan:
		return reflect.Value{}, ctx.Errorf(nil, "go chan type cannot be integrated")
	case reflect.UnsafePointer:
		return reflect.Value{}, ctx.Errorf(nil, "go unsafe pointer type cannot be integrated")
	case reflect.Func:
		return reflect.Value{}, ctx.Errorf(nil, "go func type cannot be integrated")
	case reflect.Map:
		return reflect.Value{}, ctx.Errorf(nil, "go map type cannot be integrated")
	case reflect.Interface:
		return ctx.IntegrateInterface(s, t)
	case reflect.Ptr:
		if IsEmptySymbol(s) {
			return reflect.Zero(t), nil
		} else {
			if elem, err := ctx.Integrate(s, t.Elem()); err != nil {
				return reflect.Value{}, err
			} else {
				var u reflect.Value
				if elem.CanAddr() {
					u = elem
				} else {
					u = reflect.New(elem.Type()).Elem()
					u.Set(elem)
				}
				w := reflect.New(t).Elem()
				w.Set(u.Addr())
				return w, nil
			}
		}
	case reflect.Slice:
		if IsEmptySymbol(s) {
			return reflect.Zero(t), nil
		} else {
			return ctx.IntegrateSlice(s, t)
		}
	case reflect.Struct: // catch missing fields
		return ctx.IntegrateStruct(s, t)
	}
	return reflect.Value{}, ctx.Errorf(nil, "cannot integrate %s into %v", Sprint(s), t)
}

func (ctx *typingCtx) IntegrateInterface(s Symbol, t reflect.Type) (reflect.Value, error) {
	switch s.(type) {
	case *OpaqueSymbol:
		if opaque.Type_.Type.AssignableTo(t) {
			w := reflect.New(t).Elem()
			w.Set(opaque.Value)
			return w, nil
		} else {
			return reflect.Value{},
				ctx.Errorf(nil, "cannot integrate opaque type %v into go type %v", opaque.Type_, t)
		}
	case *NamedType:
		XXX
	default:
		return reflect.Value{}, ctx.Errorf(nil, "cannot integrate %v into go interface %v", s, t)
	}
}

func (ctx *typingCtx) IntegrateSlice(s Symbol, t reflect.Type) (reflect.Value, error) {
	ss := LiftToSeries(ctx.Span, s)
	elems := make([]reflect.Value, len(ss.Elem))
	ctx2 := ctx.Refine("()")
	for i, symElem := range ss.Elem {
		if u, err := ctx2.Integrate(symElem, t.Elem()); err != nil {
			return reflect.Value{}, err
		} else {
			elems[i] = u
		}
	}
	w := reflect.MakeSlice(t, len(elems), len(elems))
	for i, elem := range elems {
		w.Index(i).Set(elem)
	}
	return w, nil
}

func (ctx *typingCtx) IntegrateStruct(s Symbol, t reflect.Type) (reflect.Value, error) {
	if ss, ok := s.(*StructSymbol); !ok {
		return reflect.Value{}, ctx.Errorf(nil, "cannot integrate non-struct %s into struct %v", Sprint(s), t)
	} else {
		w := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			toField := t.Field(i)
			if from := FindIntegrationField(ss, toField); from == nil {
				switch toField.Type.Kind() {
				case reflect.Ptr, reflect.Slice: // to field is optional
				default:
					return reflect.Value{}, ctx.Errorf(nil, "go field %s in %v is required", toField.Name, t)
				}
			} else {
				if u, err := ctx.Refine(toField.Name).Integrate(from.Value, toField.Type); err != nil {
					return reflect.Value{}, err
				} else {
					w.Field(i).Set(u)
				}
			}
		}
		return w, nil
	}
}

func FindIntegrationField(from *StructSymbol, to reflect.StructField) *FieldSymbol {
	name, hasKoName := gate.StructFieldKoName(to)
	if !hasKoName {
		return nil
	}
	if gate.IsStructFieldMonadic(to) {
		if monadicField := from.FindMonadic(); monadicField != nil {
			return monadicField
		}
	}
	return from.FindName(name)
}

func (ctx *typingCtx) IntegrateNamed(s Symbol, t reflect.Type) (reflect.Value, error) {
	if named, ok := s.(*NamedSymbol); ok {
		stype := named.GoType()
		if stype == t {
			return named.Value, nil
		} else {
			return reflect.Value{},
				ctx.Errorf(nil, "cannot integrate named type %s to named type %s", TypeName(stype), TypeName(t))
		}
	} else {
		return reflect.Value{}, ctx.Errorf(nil, "cannot integrate %v to named type %s", s, TypeName(t))
	}
}
