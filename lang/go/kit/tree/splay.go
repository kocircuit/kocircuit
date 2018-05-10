package tree

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/gate"
)

func Splay(v interface{}) Tree {
	ctx := &SplayCtx{}
	return ctx.SplayValue(reflect.ValueOf(v))
}

type SplayCtx struct {
	Seen MultiSet `ko:"name=seen"`
}

func (ctx *SplayCtx) Cache(x interface{}) *SplayCtx {
	return &SplayCtx{Seen: ctx.Seen.Add(x)}
}

type Splayer interface {
	Splay() Tree
}

type Relabeler interface {
	Relabel() Label
}

func (ctx *SplayCtx) SplayValue(v reflect.Value) Tree {
	var label Label
	if v.Kind() != reflect.Invalid && v.CanInterface() {
		if g, ok := v.Interface().(Splayer); ok {
			return g.Splay()
		}
		if g, ok := v.Interface().(Relabeler); ok {
			label = g.Relabel()
		}
	}
	switch v.Kind() {
	case reflect.Invalid:
		return NoQuote{String_: "empty"}
	case reflect.Ptr: // cache pointers
		if ctx.Seen.Count(v.Interface()) > 0 {
			return Cycle{}
		} else {
			return Sometimes{Elem: ctx.Cache(v.Interface()).SplayValue(v.Elem())}
		}
	case reflect.Interface: // cache interfaces
		if ctx.Seen.Count(v.InterfaceData()) > 0 {
			return Cycle{}
		} else {
			return ctx.Cache(v.InterfaceData()).SplayValue(v.Elem())
		}
	case reflect.Slice, reflect.Array:
		if label.IsEmpty() {
			label = TypeLabel(v.Type())
		}
		series := Series{Label: label, Bracket: "()"}
		for i := 0; i < v.Len(); i++ {
			series.Elem = append(
				series.Elem,
				IndexTree{Index: i, Tree: ctx.SplayValue(v.Index(i))},
			)
		}
		if len(series.Elem) == 0 {
			return NoQuote{String_: "empty"}
		}
		return series
	case reflect.Struct:
		if label.IsEmpty() {
			label = TypeLabel(v.Type())
		}
		parallel := Parallel{Label: label, Bracket: "()"}
		for _, f := range StripFields(v.Type()) {
			if f.IsGoExported() {
				parallel.Elem = append(parallel.Elem,
					NameTree{
						Name:    f.Name(),
						Monadic: f.IsMonadic(),
						Tree:    ctx.SplayValue(v.FieldByName(f.GoName())),
					},
				)
			}
		}
		return parallel
	case reflect.Map:
		if v.Type().Key() != reflect.TypeOf("") {
			return Opaque{}
		}
		if label.IsEmpty() {
			label = TypeLabel(v.Type())
		}
		parallel := Parallel{Label: label, Bracket: "()"}
		for _, k := range v.MapKeys() {
			key := KoGoName{Ko: k.String(), Go: k.String()}
			parallel.Elem = append(parallel.Elem,
				NameTree{
					Name:    key,
					Monadic: false,
					Tree:    ctx.SplayValue(v.MapIndex(k)),
				},
			)
		}
		if len(parallel.Elem) == 0 {
			return NoQuote{String_: "empty"}
		}
		return parallel
	case reflect.String:
		return GoValue{v}
	case reflect.Bool:
		return GoValue{v}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return GoValue{v}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return GoValue{v}
	case reflect.Float32, reflect.Float64:
		return GoValue{v}
	case reflect.Uintptr:
		return GoValue{v}
	case reflect.UnsafePointer:
		return GoValue{v}
	case reflect.Complex64, reflect.Complex128:
		return GoValue{v}
	case reflect.Chan, reflect.Func:
		return GoValue{v}
	}
	panic("o")
}

func StructFieldIsExported(f reflect.StructField) bool {
	return f.PkgPath == ""
}

func TypeLabel(t reflect.Type) Label {
	return Label{Path: t.PkgPath(), Name: t.Name()}
}

func IsBuiltinKind(kind reflect.Kind) bool {
	switch kind {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.Uintptr:
		return true
	case reflect.UnsafePointer:
		return true
	case reflect.Complex64, reflect.Complex128:
		return true
	}
	return false
}
