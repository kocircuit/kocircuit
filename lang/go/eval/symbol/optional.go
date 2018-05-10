package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type OptionalType struct {
	Elem Type `ko:"name=elem"`
}

func (*OptionalType) IsType() {}

func (ot *OptionalType) String() string {
	return Sprint(ot)
}

func (ot *OptionalType) Splay() Tree {
	return Sometimes{Elem: ot.Elem.Splay()}
}

// Optionally makes a type optional, unless it is already optional or series.
func Optionally(t Type) Type {
	switch t.(type) {
	case EmptyType:
		return t
	case *OptionalType:
		return t
	case *SeriesType:
		return t
	case BasicType, *StructType, VarietyType, NamedType:
		return &OptionalType{Elem: t}
	}
	panic("o")
}
