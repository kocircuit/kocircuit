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

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func DeconstructInterface(span *model.Span, any interface{}) Symbol {
	return Deconstruct(span, reflect.ValueOf(any))
}

func Deconstruct(span *model.Span, v reflect.Value) Symbol {
	ctx := &typingCtx{Span: span}
	if symbol, err := ctx.Deconstruct(v); err != nil {
		panic(err)
	} else {
		return symbol
	}
}

func DeconstructKind(span *model.Span, v reflect.Value) Symbol {
	ctx := &typingCtx{Span: span}
	if symbol, err := ctx.DeconstructKind(v); err != nil {
		panic(err)
	} else {
		return symbol
	}
}

func (ctx *typingCtx) Deconstruct(v reflect.Value) (Symbol, error) {
	if v.IsValid() {
		if typeName := tree.TypeName(v.Type()); typeName != "" && v.Kind() != reflect.Interface {
			return &NamedSymbol{Value: v}, nil
		}
	}
	return ctx.DeconstructKind(v)
}

func (ctx *typingCtx) DeconstructKind(v reflect.Value) (Symbol, error) {
	if v.IsValid() && v.Type() == typeOfSymbol {
		if v.IsNil() {
			return EmptySymbol{}, nil
		} else {
			return v.Interface().(Symbol), nil
		}
	}
	switch v.Kind() {
	case reflect.Invalid:
		return EmptySymbol{}, nil
	case reflect.String: // string
		return BasicSymbol{string(v.String())}, nil
	case reflect.Bool: // bool
		return BasicSymbol{bool(v.Bool())}, nil
	case reflect.Int: // signed integers
		return BasicSymbol{int64(v.Int())}, nil
	case reflect.Int8:
		return BasicSymbol{int8(v.Int())}, nil
	case reflect.Int16:
		return BasicSymbol{int16(v.Int())}, nil
	case reflect.Int32:
		return BasicSymbol{int32(v.Int())}, nil
	case reflect.Int64:
		return BasicSymbol{int64(v.Int())}, nil
	case reflect.Uint: // unsigned integers
		return BasicSymbol{uint64(v.Uint())}, nil
	case reflect.Uint8:
		return BasicSymbol{uint8(v.Uint())}, nil
	case reflect.Uint16:
		return BasicSymbol{uint16(v.Uint())}, nil
	case reflect.Uint32:
		return BasicSymbol{uint32(v.Uint())}, nil
	case reflect.Uint64:
		return BasicSymbol{uint64(v.Uint())}, nil
	case reflect.Float32: // floating point
		return BasicSymbol{float32(v.Float())}, nil
	case reflect.Float64:
		return BasicSymbol{float64(v.Float())}, nil
	case reflect.Uintptr:
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Complex64:
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Complex128:
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Array: // non-protocol type, go-specific
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Chan: // go-specific
		return &OpaqueSymbol{Value: v}, nil
	case reflect.UnsafePointer: // go-specific
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Func: // go-specific
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Map:
		if v.Type().Key() == typeOfString {
			return ctx.DeconstructMap(v) //XXX
		} else {
			return &OpaqueSymbol{Value: v}, nil
		}
	case reflect.Interface:
		return &OpaqueSymbol{Value: v}, nil
	case reflect.Ptr:
		if v.IsNil() {
			return EmptySymbol{}, nil
		} else {
			return ctx.Deconstruct(v.Elem())
		}
	case reflect.Slice:
		if v.Type() == byteSliceType {
			return &BlobSymbol{Value: v}, nil
		} else {
			return ctx.DeconstructSlice(v)
		}
	case reflect.Struct:
		return ctx.DeconstructStruct(v)
	}
	panic("o")
}

var byteSliceType = reflect.TypeOf([]byte{1})
var typeOfString = reflect.TypeOf(string(""))

func (ctx *typingCtx) DeconstructSlice(v reflect.Value) (Symbol, error) {
	ds := make(Symbols, 0, v.Len())
	dt := make(Types, 0, v.Len())
	valueType := v.Type().Elem()
	ctx2 := ctx.Refine("()")
	for i := 0; i < v.Len(); i++ {
		if elem, err := ctx2.Deconstruct(v.Index(i).Convert(valueType)); err != nil {
			return nil, err
		} else {
			if !IsEmptySymbol(elem) {
				ds = append(ds, elem)
				dt = append(dt, elem.Type())
			}
		}
	}
	if ss, err := MakeSeriesSymbol(ctx.Span, ds); err != nil {
		panic(err)
	} else {
		return ss, nil
	}
}

func (ctx *typingCtx) DeconstructStruct(v reflect.Value) (Symbol, error) {
	fields := make(FieldSymbols, 0, v.NumField())
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		name, hasKoName := gate.StructFieldKoProtoGoName(t.Field(i))
		if !hasKoName {
			continue // skip
		}
		monadic := gate.IsStructFieldMonadic(t.Field(i))
		if value, err := ctx.Refine(name).Deconstruct(v.Field(i)); err != nil {
			return nil, err
		} else {
			fields = append(fields,
				&FieldSymbol{Name: name, Monadic: monadic, Value: value},
			)
		}
	}
	return MakeStructSymbol(fields), nil
}

// v must be map[string]T
func (ctx *typingCtx) DeconstructMap(v reflect.Value) (Symbol, error) {
	mapKeys := v.MapKeys()
	dv := map[string]Symbol{}
	valueType := v.Type().Elem()
	for _, key := range mapKeys {
		if dvalue, err := ctx.Deconstruct(v.MapIndex(key).Convert(valueType)); !IsEmptySymbol(dvalue) {
			if err != nil {
				panic("o")
			}
			dv[key.String()] = dvalue
		}
	}
	return MakeMapSymbol(ctx.Span, dv)
}
