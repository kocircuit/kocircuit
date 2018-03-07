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

type ExplainCtx struct {}

func (ctx *ExplainCtx) Explain(t reflect.Type) Tree {
	if name := TypeName(t); name != "" {
		return NoQuote{name}
	} else {
		return ctx.ExplainNoName(t)
	}
}

func (ctx *ExplainCtx) ExplainNoName(t reflect.Type) Tree {
	switch t.Kind() {
	case reflect.Invalid:
		return NoQuote{"Empty"}
	case reflect.Ptr:
		return Sometimes{Elem: ctx.Explain(t.Elem())}
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return NoQuote{"interface{}"}
		} else {
			return NoQuote{"interface{...}"}
		}
	case reflect.Array:
		return Series{
			Label: TypeLabel(t),
			Bracket: "{}",
			Elem: []IndexTree{{
				Index: 0,
				Tree: ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Slice:
		return Series{
			Label: TypeLabel(t),
			Bracket: "()", 
			Elem: []IndexTree{{
				Index: 0,
				Tree: ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Struct:
		parallel := Parallel{
			Label: TypeLabel(t),
			Bracket: "()",
		}
		for _, f := range StripFields(t) {
			if f.IsGoExported() {
				field, _ := t.FieldByName(f.GoName())
				parallel.Elem = append(parallel.Elem,
					NameTree{
						Name: f.Name(),
						Monadic: f.IsMonadic(),
						Tree: ctx.Explain(field.Type),
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
			Label: TypeLabel(t), 
			Bracket: "()",
			Elem: []NameTree{{
				Name: KoGoName{Ko: "<key>", Go: "<key>"},
				Monadic: false,
				Tree: ctx.Explain(t.Elem()),
			}},
		}
	case reflect.Chan, reflect.Func:
		return Opaque{}
	case reflect.String:
		return NoQuote{"String"}
	case reflect.Bool:
		return NoQuote{"Bool"}
	case reflect.Float32:
		return NoQuote{"Float32"}
	case reflect.Float64:
		return NoQuote{"Float64"}
	case reflect.Uintptr:
		return NoQuote{"Uintptr"}
	case reflect.UnsafePointer:
		return NoQuote{"UnsafePointer"}
	case reflect.Complex64:
		return NoQuote{"Complex64"}
	case reflect.Complex128:
		return NoQuote{"Complex128"}
	case reflect.Int:
		return NoQuote{"Int"}
	case reflect.Int8:
		return NoQuote{"Int8"}
	case reflect.Int16:
		return NoQuote{"Int16"}
	case reflect.Int32:
		return NoQuote{"Int32"}
	case reflect.Int64:
		return NoQuote{"Int64"}
	case reflect.Uint:
		return NoQuote{"Uint"}
	case reflect.Uint8:
		return NoQuote{"Uint8"}
	case reflect.Uint16:
		return NoQuote{"Uint16"}
	case reflect.Uint32:
		return NoQuote{"Uint32"}
	case reflect.Uint64:
		return NoQuote{"Uint64"}
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
