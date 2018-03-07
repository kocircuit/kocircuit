package tree

import (
	"strconv"
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/text"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Tree interface {
	Text() Textual
	TreeHash() string
}

// Sometimes is a tree.
type Sometimes struct {
	Elem Tree `ko:"name=elem"`
}

func (v Sometimes) TreeHash() string {
	return Mix("sometimes", v.Elem.TreeHash())
}

func (v Sometimes) Text() Textual {
	return TextGlue{
		Text: []Textual{TextSlab{String: "*"}, v.Elem.Text()},
	}
}

// Series is a tree.
type Series struct {
	Label Label       `ko:"name=label"`
	Bracket string `ko:"name=bracket"` // "{}", "[]"
	Elem  []IndexTree `ko:"name=elem"`
}

type Label struct {
	Path string `ko:"name=path"`
	Name string `ko:"name=name"`
}

func (l Label) IsEmpty() bool { return l.Path == "" && l.Name == "" }

type IndexTree struct {
	Index int `ko:"name=index"`
	Tree  `ko:"name=tree"`
}

func (series Series) TreeHash() string {
	s := []string{"series"}
	for _, e := range series.Elem {
		s = append(s, e.Tree.TreeHash())
	}
	return Mix(s...)
}

func (series Series) Text() Textual {
	g := TextRubber{
		Header: PtrString(series.Label.Name),
		Open:   PtrString(series.Bracket[:1]), Close: PtrString(series.Bracket[1:]),
		Field: make([]Textual, len(series.Elem)),
	}
	for i := range series.Elem {
		g.Field[i] = series.Elem[i].Text()
	}
	return g
}

// Parallel is a tree.
type Parallel struct {
	Label Label      `ko:"name=label"`
	Bracket string `ko:"name=bracket"` // "{}", "[]"
	Elem  []NameTree `ko:"name=field"`
}

type NameTree struct {
	Name KoGoName `ko:"name=name"`
	Monadic bool `ko:"name=monadic"`
	Tree `ko:"name=tree"`
}

func (parallel Parallel) TreeHash() string {
	s := []string{"parallel"}
	for _, e := range parallel.Elem {
		s = append(s, e.Name.Ko, strconv.FormatBool(e.Monadic), e.Tree.TreeHash())
	}
	return Mix(s...)
}

func (parallel Parallel) Text() Textual {
	g := TextRubber{
		Header: PtrString(parallel.Label.Name),
		Open:   PtrString(parallel.Bracket[:1]), Close: PtrString(parallel.Bracket[1:]),
	}
	for i := range parallel.Elem {
		g.Field = append(
			g.Field,
			TextGlue{
				Text: []Textual{
					TextTile{
						String: NameWithMonadicity(
							parallel.Elem[i].Name.Ko,
							parallel.Elem[i].Monadic,
						),
					},
					TextSlab{String: ": "},
					parallel.Elem[i].Text(),
				},
			},
		)
	}
	return g
}

func NameWithMonadicity(name string, monadic bool) string {
	if monadic {
		return name + "?"
	} else {
		return name
	}
}

// Cycle is a tree.
type Cycle struct{}

func (v Cycle) TreeHash() string {
	return Mix("cycle")
}

func (v Cycle) Text() Textual {
	return TextSlab{String: "backlink"}
}

// Opaque is a tree.
type Opaque struct{}

func (v Opaque) TreeHash() string {
	return Mix("opaque")
}

func (v Opaque) Text() Textual {
	return TextSlab{String: "█"}
}

// Quote is a tree.
type Quote struct {
	String_ string `ko:"name=string"`
}

func (v Quote) TreeHash() string {
	return Mix("quote", v.String_)
}

func (v Quote) Text() Textual {
	return TextSlab{String: fmt.Sprintf("%q", v.String_)}
}

// NoQuote is a tree.
type NoQuote struct {
	String_ string `ko:"name=string"`
}

func (v NoQuote) TreeHash() string {
	return Mix("noquote", v.String_)
}

func (v NoQuote) Text() Textual {
	return TextSlab{String: v.String_}
}

// GoValue is a tree.
type GoValue struct {
	reflect.Value `ko:"name=value"`
}

// From Go doc:
// 	func (v Value) String() string
// 	    String returns the string v's underlying value, as a string. String is a
// 	    special case because of Go's String method convention. Unlike the other
// 	    getters, it does not panic if v's Kind is not String. Instead, it returns a
// 	    string of the form "<T value>" where T is v's type. The fmt package treats
// 	    Values specially. It does not call their String method implicitly but
// 	    instead prints the concrete values they hold.
func (v GoValue) TreeHash() string {
	return MixValue(v.Value)
}

func (v GoValue) Text() Textual {
	switch v.Value.Kind() {
	case reflect.String:
		return TextSlab{String: fmt.Sprintf("%q", v.String())}
	case reflect.Bool:
		return TextSlab{String: fmt.Sprint(v.Bool())}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return TextSlab{String: fmt.Sprint(v.Int())}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return TextSlab{String: fmt.Sprint(v.Uint())}
	case reflect.Float32, reflect.Float64:
		return TextSlab{String: fmt.Sprint(v.Float())}
	case reflect.Uintptr:
		return TextSlab{String: fmt.Sprint(v.Pointer())}
	case reflect.UnsafePointer:
		return TextSlab{String: fmt.Sprint(v.Pointer())}
	case reflect.Complex64, reflect.Complex128:
		return TextSlab{String: fmt.Sprint(v.Complex())}
	case reflect.Invalid:
		return nil
	case reflect.Chan, reflect.Func:
		return TextSlab{String: v.Value.Type().String()}
	}
	panic("o")
}
