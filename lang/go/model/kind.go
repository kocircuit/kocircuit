package model

func IsPtr(typ GoType) bool {
	_, ok := typ.(*GoPtr)
	return ok
}

func IsArray(typ GoType) bool {
	_, ok := typ.(*GoArray)
	return ok
}

func IsSlice(typ GoType) bool {
	_, ok := typ.(*GoSlice)
	return ok
}

type Kind int

const (
	KindInvalid = Kind(iota)
	KindStruct
	KindVariety
	KindPure
	KindUnknown
)

func (k Kind) String() string {
	switch k {
	case KindInvalid:
		panic("invalid")
	case KindStruct:
		return "structure"
	case KindVariety:
		return "variety"
	case KindPure:
		return "pure"
	case KindUnknown:
		return "unknown"
	}
	panic("o")
}

func KindOf(t GoType) Kind {
	switch t.(type) {
	case *GoStruct:
		return KindStruct
	case Unknown:
		return KindUnknown
	case GoVarietal: // must come after Unknown, beause Unknown is GoVarietal
		return KindVariety
	}
	return KindPure
}
