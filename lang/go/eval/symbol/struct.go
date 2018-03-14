package symbol

import (
	"sort"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// XXX: Hash and Equal ignore empty fields

func MakeStructSymbol(fields FieldSymbols) *StructSymbol {
	// fields = FilterEmptyFieldSymbols(fields) //XXX: necessary????
	return &StructSymbol{
		Type_: &StructType{Field: FieldSymbolTypes(fields)},
		Field: fields,
	}
}

func FilterEmptyFieldSymbols(fields FieldSymbols) (filtered FieldSymbols) {
	filtered = make(FieldSymbols, 0, len(fields))
	for _, field := range fields {
		if !IsEmptySymbol(field.Value) {
			filtered = append(filtered, field)
		}
	}
	return
}

func FieldSymbolTypes(fields FieldSymbols) []*FieldType {
	types := make([]*FieldType, len(fields))
	for i, field := range fields {
		types[i] = &FieldType{
			Name:  field.Name,
			Type_: field.Value.Type(),
		}
	}
	return types
}

type StructSymbol struct {
	Type_ *StructType  `ko:"name=type"`
	Field FieldSymbols `ko:"name=field"`
}

type FieldSymbol struct {
	Name    string `ko:"name=name"`
	Monadic bool   `ko:"name=monadic"`
	Value   Symbol `ko:"name=value"`
}

func (ss *StructSymbol) IsEmpty() bool {
	return len(ss.Field) == 0
}

func (ss *StructSymbol) String() string {
	return Sprint(ss)
}

func (ss *StructSymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*StructSymbol); ok {
		return FieldSymbolsEqual(ss.Field, other.Field)
	} else {
		return false
	}
}

func FieldSymbolsEqual(x, y FieldSymbols) bool {
	x, y = FilterEmptyFieldSymbols(x), FilterEmptyFieldSymbols(y)
	if len(x) != len(y) {
		return false
	}
	u, v := x.Copy(), y.Copy()
	u.Sort()
	v.Sort()
	for i := range u {
		if u[i].Name != v[i].Name || !u[i].Value.Equal(v[i].Value) {
			return false
		}
	}
	return true
}

func (ss *StructSymbol) Hash() string {
	return FieldSymbolsHash(ss.Field)
}

func FieldSymbolsHash(fields FieldSymbols) string {
	fields = FilterEmptyFieldSymbols(fields)
	h := make([]string, 2*len(fields))
	for i, field := range fields {
		h[2*i] = field.Name
		h[2*i+1] = field.Value.Hash()
	}
	return Mix(h...)
}

func (ss *StructSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "structure %v cannot be augmented", ss)
}

func (ss *StructSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "structure %v cannot be invoked", ss)
}

func (ss *StructSymbol) Type() Type {
	return ss.Type_
}

func (ss *StructSymbol) Splay() Tree {
	nameTrees := make([]NameTree, len(ss.Field))
	for i, field := range ss.Field {
		nameTrees[i] = NameTree{
			Name:    gate.KoGoName{Ko: field.Name},
			Monadic: field.Monadic,
			Tree:    field.Value.Splay(),
		}
	}
	return Parallel{
		Label:   Label{"", ""},
		Bracket: "()",
		Elem:    nameTrees,
	}
}

func (ss *StructSymbol) FindMonadic() *FieldSymbol {
	for _, fs := range ss.Field {
		if fs.Monadic || fs.Name == "" {
			return fs
		}
	}
	return nil
}

func (ss *StructSymbol) FindName(name string) *FieldSymbol {
	for _, fs := range ss.Field {
		if fs.Name == name {
			return fs
		}
	}
	return nil
}

type FieldSymbols []*FieldSymbol

func (fs FieldSymbols) Copy() FieldSymbols {
	c := make(FieldSymbols, len(fs))
	copy(c, fs)
	return c
}

func (fs FieldSymbols) Sort() {
	sort.Sort(fs)
}

func (fs FieldSymbols) Len() int {
	return len(fs)
}

func (fs FieldSymbols) Less(i, j int) bool {
	return fs[i].Name < fs[j].Name
}

func (fs FieldSymbols) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}

type StructType struct {
	Field []*FieldType `ko:"name=field"`
}

type FieldType struct {
	Name  string `ko:"name=name"`
	Type_ Type   `ko:"name=type"`
}

func (*StructType) IsType() {}

func (st *StructType) String() string {
	return Sprint(st)
}

func (st *StructType) Splay() Tree {
	nameTrees := make([]NameTree, len(st.Field))
	for i, field := range st.Field {
		nameTrees[i] = NameTree{
			Name: gate.KoGoName{Ko: field.Name},
			Tree: field.Type_.Splay(),
		}
	}
	return Parallel{
		Label:   Label{"", ""},
		Bracket: "()",
		Elem:    nameTrees,
	}
}
