package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Type implementations:
//	BasicType, EmptyType,
// *SeriesType, *StructType, *NamedType, *OpaqueType, VarietyType
// *OptionalType
type Type interface {
	Splayer
	IsType()
}

type Types []Type
