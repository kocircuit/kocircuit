package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Type implementations:
// *OptionalType
//	BasicType, EmptyType,
// *SeriesType, *StructType
// NamedType, *OpaqueType
// VarietyType
// BlobType
type Type interface {
	String() string
	Splayer
	IsType()
}

type Types []Type
