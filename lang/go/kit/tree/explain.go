package tree

import (
	"reflect"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/go/gate"
)

func Explain(t reflect.Type) Tree {
	ctx := &ExplainCtx{}
	return ctx.Explain(t)
}

type ExplainCtx struct{}

func (ctx *ExplainCtx) Explain(t reflect.Type) Tree {
	if name := TypeName(t); name != "" {
		return NoQuote{String_: name}
	} else {
		return ctx.ExplainNoName(t)
	}
}

func (ctx *ExplainCtx) ExplainNoName(t reflect.Type) Tree {
	switch t.Kind() {
	case reflect.Invalid:
		return NoQuote{String_: "Empty"}
	case reflect.Ptr:
		return Sometimes{Elem: ctx.Explain(t.Elem())}
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return NoQuote{String_: "interface{}"}
		} else {
			return NoQuote{String_: "interface{...}"}
		}
	case reflect.Array:
		return Series{
			Label:   TypeLabel(t),
			Bracket: "{}",
			Elem: []IndexTree{{
				Index: 0,
				Tree:  ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Slice:
		return Series{
			Label:   TypeLabel(t),
			Bracket: "()",
			Elem: []IndexTree{{
				Index: 0,
				Tree:  ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Struct:
		parallel := Parallel{
			Label:   TypeLabel(t),
			Bracket: "()",
		}
		for _, f := range StripFields(t) {
			if f.IsGoExported() {
				field, _ := t.FieldByName(f.GoName())
				parallel.Elem = append(parallel.Elem,
					NameTree{
						Name:    f.Name(),
						Monadic: f.IsMonadic(),
						Tree:    ctx.Explain(field.Type),
					},
				)
			}
		}
		return parallel
	case reflect.Map:
		if t.Key() != reflect.TypeOf("") {
			return Opaque{}
		}
		return Parallel{
			Label:   TypeLabel(t),
			Bracket: "()",
			Elem: []NameTree{{
				Name:    KoGoName{Ko: "<key>", Go: "<key>"},
				Monadic: false,
				Tree:    ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Chan, reflect.Func:
		return Opaque{}
	case reflect.String:
		return NoQuote{String_: "String"}
	case reflect.Bool:
		return NoQuote{String_: "Bool"}
	case reflect.Float32:
		return NoQuote{String_: "Float32"}
	case reflect.Float64:
		return NoQuote{String_: "Float64"}
	case reflect.Uintptr:
		return NoQuote{String_: "Uintptr"}
	case reflect.UnsafePointer:
		return NoQuote{String_: "UnsafePointer"}
	case reflect.Complex64:
		return NoQuote{String_: "Complex64"}
	case reflect.Complex128:
		return NoQuote{String_: "Complex128"}
	case reflect.Int:
		return NoQuote{String_: "Int"}
	case reflect.Int8:
		return NoQuote{String_: "Int8"}
	case reflect.Int16:
		return NoQuote{String_: "Int16"}
	case reflect.Int32:
		return NoQuote{String_: "Int32"}
	case reflect.Int64:
		return NoQuote{String_: "Int64"}
	case reflect.Uint:
		return NoQuote{String_: "Uint"}
	case reflect.Uint8:
		return NoQuote{String_: "Uint8"}
	case reflect.Uint16:
		return NoQuote{String_: "Uint16"}
	case reflect.Uint32:
		return NoQuote{String_: "Uint32"}
	case reflect.Uint64:
		return NoQuote{String_: "Uint64"}
	}
	panic("o")
}

func TypeName(t reflect.Type) string {
	if t == nil || t.Name() == "" || t.PkgPath() == "" {
		return ""
	} else {
		return strings.Join([]string{t.PkgPath(), t.Name()}, ".")
	}
}
