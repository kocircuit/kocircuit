package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type OptionalType struct {
	Elem Type `ko:"name=elem"`
}

func (*OptionalType) IsType() {}

func (ot *OptionalType) Splay() Tree {
	return Sometimes{ot.Elem.Splay()}
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
	case BasicType, *StructType, VarietyType:
		return &OptionalType{Elem: t}
	}
	panic("o")
}
