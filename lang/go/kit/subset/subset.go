package subset

import (
	"fmt"
	"reflect"
	"strings"
)

func IsEqual(u, v interface{}) bool {
	return IsSubset(u, v) && IsSubset(v, u)
}

func IsSubset(u, v interface{}) bool {
	return VerifyIsSubset(u, v) == nil
}

func VerifyIsSubset(u, v interface{}) error {
	return verifyIsSubset(nil, reflect.ValueOf(u), reflect.ValueOf(v))
}

func verifyIsSubset(ctx []string, u, v reflect.Value) error {
	if u.Kind() == reflect.Invalid { // nil is a subset of everything
		return nil
	}
	if u.Kind() != v.Kind() {
		return fmt.Errorf("incompatible go kinds (%v and %v) at %s", u.Kind(), v.Kind(), strings.Join(ctx, "."))
	}
	if IsBasicValue(u) {
		if u.Interface() == v.Interface() {
			return nil
		}
		return fmt.Errorf("unequal builtin values (%v and %v) at %s", u.Interface(), v.Interface(), strings.Join(ctx, "."))
	}
	switch u.Kind() {
	case reflect.Array:
		if u.Len() != v.Len() {
			return fmt.Errorf("array lengths %v and %v do not match", u, v)
		}
		for i := 0; i < u.Len(); i++ {
			if err := verifyIsSubset(appendString(ctx, fmt.Sprintf("[%d]", i)), u.Index(i), v.Index(i)); err != nil {
				return err
			}
		}
		return nil
	case reflect.Chan:
		panic("no subset for chan")
	case reflect.Func:
		panic("no subset for func")
	case reflect.Ptr, reflect.Interface:
		if u.IsNil() {
			return nil
		}
		if v.IsNil() {
			return fmt.Errorf("non-nil %v is not subset to nil at %s", u.Interface(), strings.Join(ctx, "."))
		}
		return verifyIsSubset(appendString(ctx, "*"), u.Elem(), v.Elem())
	case reflect.Slice:
		if u.Len() > v.Len() {
			return fmt.Errorf("slice length %d is greater than %d at %s", u.Len(), v.Len(), strings.Join(ctx, "."))
		}
		for i := 0; i < u.Len(); i++ {
			if err := verifyIsSubset(appendString(ctx, fmt.Sprintf("[%d]", i)), u.Index(i), v.Index(i)); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		for _, k := range u.MapKeys() {
			if err := verifyIsSubset(appendString(ctx, fmt.Sprintf("map[%v]", k.Interface())), u.MapIndex(k), v.MapIndex(k)); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		for i := 0; i < u.NumField(); i++ {
			if err := verifyIsSubset(appendString(ctx, u.Type().Field(i).Name), u.Field(i), v.FieldByName(u.Type().Field(i).Name)); err != nil {
				return err
			}
		}
		return nil
	}
	panic("o")
}

func appendString(to []string, a ...string) []string {
	return append(append([]string{}, to...), a...)
}

func IsBasicValue(v reflect.Value) bool {
	return IsBasicType(v.Type())
}

func IsBasicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String:
		fallthrough
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fallthrough
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fallthrough
	case reflect.UnsafePointer:
		fallthrough
	case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}
