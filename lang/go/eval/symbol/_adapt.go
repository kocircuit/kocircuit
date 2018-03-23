package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func Adapt(span *Span, s reflect.Value, t reflect.Type) (reflect.Value, error) {
	ctx := &typingCtx{Span: span}
	return ctx.Adapt(s, t)
}

func (ctx *typingCtx) Adapt(s reflect.Value, t reflect.Type) (reflect.Value, error) {
	for s.Type() == reflect.Ptr {
		s = s.Elem()
	}
	switch t.Kind() {
	case reflect.Invalid:
		panic("o")
	case reflect.String:
		return ctx.adaptConvertible(s, t)
	case reflect.Bool:
		return ctx.adaptConvertible(s, t)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ctx.adaptConvertibleBits(s, t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ctx.adaptConvertibleBits(s, t)
	case reflect.Float32, reflect.Float64:
		return ctx.adaptConvertibleBits(s, t)
	case reflect.Complex64, reflect.Complex128:
		return ctx.adaptConvertibleBits(s, t)
	case reflect.Uintptr:
		return ctx.adaptConvertible(s, t)
	case reflect.Array:
		return ctx.adaptConvertible(s, t)
	case reflect.UnsafePointer:
		return ctx.adaptConvertible(s, t)
	case reflect.Chan:
		return ctx.adaptConvertibleOrNil(s, t)
	case reflect.Func:
		return ctx.adaptConvertibleOrNil(s, t)
	case reflect.Interface:
		return ctx.adaptConvertibleOrNil(s, t)
	case reflect.Map:
		return ctx.adaptConvertibleOrNil(s, t)
	case reflect.Ptr:
		if s.IsValid() {
			if adaptedElem, err := ctx.Adapt(s, t.Elem()); err != nil {
				return reflect.Value{}, err
			} else {
				w := reflect.New(t).Elem()
				w.Set(EnsureAddressible(adaptedElem).Addr())
				return w, nil
			}
		} else {
			return reflect.Zero(t), nil
		}
	case reflect.Slice:
		if s.IsValid() {
			if s.Kind() == reflect.Slice { // []P->[]Q
				if w, err := ctx.adaptConvertible(s, t); err == nil { // directly convertible,
					return w, nil
				} else {
					return ctx.adaptSliceSlice(s, t) // XXX: or element by element
				}
			} else { // P->[]Q
				if elem, err := ctx.Adapt(s, t.Elem()); err != nil {
					return reflect.Value{}, err
				} else {
					XXX //XXX: build singleton slice
				}
			}
		} else {
			return reflect.Zero(t), nil
		}
	case reflect.Struct:
		if s.IsValid() {
			if s.Kind() == reflect.Struct {
				return ctx.adaptStructStruct(s, t) // XXX: struct->struct
			} else if s.Kind() == reflect.Map && s.Type().Key().Kind() == reflect.String {
				return ctx.adaptMapStruct(s, t) // XXX: map[string]T->struct
			}
		}
	}
	return reflect.Value{}, ctx.Errorf(nil, "cannot adapt %s to %v", Sprint(s.Interface()), t)
}

func EnsureAddressible(v reflect.Value) (addressible reflect.Value) {
	if v.CanAddr() {
		return v
	} else {
		addressible = reflect.New(v.Type()).Elem()
		addressible.Set(v)
		return addressible
	}
}

func (ctx *typingCtx) adaptConvertibleOrNil(s reflect.Value, t reflect.Type) (reflect.Value, error) {
	if s.IsValid() {
		return ctx.adaptConvertible(s, t)
	} else {
		return reflect.Zero(t), nil
	}
}

func (ctx *typingCtx) adaptConvertible(s reflect.Value, t reflect.Type) (reflect.Value, error) {
	if !s.IsValid() {
		return reflect.Value{}, ctx.Errorf(nil, "empty is not adaptible to %v", t)
	}
	if s.Type().ConvertibleTo(t) {
		return s.Convert(t), nil
	} else {
		return reflect.Value{}, ctx.Errorf(nil, "%s not adaptible to %v", Sprint(s.Interface()), t)
	}
}

func (ctx *typingCtx) adaptConvertibleBits(s reflect.Value, t reflect.Type) (reflect.Value, error) {
	if !s.IsValid() {
		return reflect.Value{}, ctx.Errorf(nil, "empty is not adaptible to %v", t)
	}
	if s.Type().ConvertibleTo(t) && s.Type().Bits() <= t.Bits() {
		return s.Convert(t), nil
	} else {
		return reflect.Value{}, ctx.Errorf(nil, "%s not adaptible to %v", Sprint(s.Interface()), t)
	}
}
