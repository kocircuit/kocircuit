package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func GeneralizeSequence(span *Span, seq ...GoType) (generalized GoType, err error) {
	ctx := newGeneralizeCtx(span)
	generalized = seq[0]
	for _, s := range seq[1:] {
		if generalized, err = generalizeType(ctx, generalized, s); err != nil {
			return nil, err
		}
	}
	return
}

// Generalizing types X and Y returns type XY, such that X and Y assign to XY.
// Generalization is commutative and associative (semantically, not necessarily go structure-wise):
//	G(x,y) = G(y,x)
//	G(x,G(y, z)) = G(G(x,y),z)
// A(x, G(x, y)) always holds
func Generalize(span *Span, x, y GoType) (xy GoType, err error) {
	return generalizeType(newGeneralizeCtx(span), x, y)
}

func GeneralizeStructure(span *Span, x, y GoStructure) (GoStructure, error) {
	if xy, err := generalizeType(newGeneralizeCtx(span), x, y); err != nil {
		return nil, err
	} else {
		return xy.(GoStructure), nil
	}
}

func generalizeType(ctx *generalizeCtx, x, y GoType) (GoType, error) {
	if typer, err := ctx.Generalize(x, y); err != nil {
		return nil, err
	} else {
		return typer.Type(), nil
	}
}

type generalizationID struct{ X, Y string }

type generalizeCtx struct {
	Span *Span
	Seen map[generalizationID]GoRetyper
}

func newGeneralizeCtx(span *Span) *generalizeCtx {
	return &generalizeCtx{Span: span, Seen: map[generalizationID]GoRetyper{}}
}

func (ctx *generalizeCtx) Generalize(x, y GoType) (xy GoRetyper, err error) {
	sx, _ := Simplify(ctx.Span, x)
	sy, _ := Simplify(ctx.Span, y)
	genID := generalizationID{sx.TypeID(), sy.TypeID()}
	if seen := ctx.Seen[genID]; seen != nil {
		return seen, nil
	} else {
		reserve := &GoReserveRetyper{}
		ctx.Seen[genID] = reserve
		if reserve.Solution, err = ctx.GeneralizeMatrix(sx, sy); err != nil {
			return nil, err
		} else {
			return reserve, nil
		}
	}
}

func (ctx *generalizeCtx) GeneralizeMatrix(x, y GoType) (xy GoRetyper, err error) {
	switch u := x.(type) {
	case Unknown:
		switch w := y.(type) {
		case Unknown:
			return &GoTypeRetyper{NewGoUnknown(ctx.Span)}, nil
		default:
			return &GoTypeRetyper{w}, nil
		}
	case *GoEmpty:
		switch w := y.(type) {
		case Unknown:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoEmpty:
			return &GoTypeRetyper{u}, nil
		case *GoInterface:
			if w.IsRestricted() {
				return nil, ctx.Span.Errorf(nil, "restricted interfaces are not generalizable %s", Sprint(w))
			} else {
				return &GoTypeRetyper{w}, nil
			}
		case *GoPtr, *GoSlice:
			return &GoTypeRetyper{w}, nil
		case *GoArray:
			return &GoTypeRetyper{NewGoSlice(w.Elem)}, nil
		case GoNumber:
			return &GoTypeRetyper{w}, nil
		// case *GoStruct: // empty <-> () equivalence?
		// 	XXX
		default:
			return &GoTypeRetyper{NewGoPtr(w)}, nil
		}
	case *GoInterface:
		if u.IsRestricted() {
			return nil, ctx.Span.Errorf(nil, "restricted interfaces are not generalizable %s", Sprint(u))
		} else {
			switch w := y.(type) {
			case Unknown, *GoEmpty:
				return ctx.GeneralizeMatrix(w, u) // matrix symmetry
			default:
				return &GoTypeRetyper{u}, nil
			}
		}
	case *GoPtr:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoPtr:
			if xy, err = ctx.Generalize(u.Elem, w.Elem); err != nil {
				return nil, err
			}
			return &GoPtrRetyper{Elem: xy}, nil
		case *GoSlice:
			if xy, err = ctx.Generalize(u.Elem, w.Elem); err != nil {
				return nil, err
			}
			return &GoSliceRetyper{Elem: xy}, nil
		default: // always generalize pointers?
			if xy, err = ctx.Generalize(u.Elem, w); err != nil {
				return nil, err
			}
			return &GoPtrRetyper{Elem: xy}, nil
		}
	case *GoArray:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoArray:
			return ctx.GeneralizeArray(u, w)
		case *GoSlice:
			return ctx.GeneralizeArraySlice(u, w)
		default:
			if xy, err = ctx.Generalize(u.Elem, w); err != nil {
				return nil, err
			}
			return &GoSliceRetyper{Elem: xy}, nil
		}
	case *GoSlice:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface, *GoPtr, *GoArray:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoSlice:
			return ctx.GeneralizeSlice(u, w)
		default:
			if xy, err = ctx.Generalize(u.Elem, w); err != nil {
				return nil, err
			}
			return &GoSliceRetyper{Elem: xy}, nil
		}
	case *GoVariety:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface, *GoPtr, *GoSlice:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoVariety:
			return ctx.GeneralizeVariety(u, w)
		}
	case *GoStruct:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface, *GoPtr, *GoSlice:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case *GoStruct:
			return ctx.GeneralizeStructure(u, w)
		}
	case *GoBuiltin:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface, *GoPtr, *GoSlice:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case GoNumber:
			return ctx.GeneralizeBuiltin(u, w.Builtin())
		case *GoBuiltin:
			return ctx.GeneralizeBuiltin(u, w)
		}
	case GoNumber:
		switch w := y.(type) {
		case Unknown, *GoEmpty, *GoInterface, *GoPtr, *GoSlice, *GoBuiltin:
			return ctx.GeneralizeMatrix(w, u) // matrix symmetry
		case GoNumber:
			return ctx.GeneralizeBuiltin(u.Builtin(), w.Builtin())
		}
	}
	return nil, ctx.Span.Errorf(nil, "%s and %s cannot be generalized", Sprint(x), Sprint(y))
}

func (ctx *generalizeCtx) GeneralizeVariety(x, y *GoVariety) (_ GoRetyper, err error) {
	switch {
	case Assignable(ctx.Span, x, y):
		return &GoTypeRetyper{y}, nil
	case Assignable(ctx.Span, y, x):
		return &GoTypeRetyper{x}, nil
	default:
		if w, err := GoVarietyUnion(ctx.Span, x, y); err != nil {
			return nil, err
		} else {
			return &GoTypeRetyper{w}, nil
		}
	}
}

func (ctx *generalizeCtx) GeneralizeArray(x, y *GoArray) (_ GoRetyper, err error) {
	var elem GoRetyper
	if elem, err = ctx.Generalize(x.Elem, y.Elem); err != nil {
		return nil, err
	} else {
		if x.Len != y.Len {
			return &GoSliceRetyper{Elem: elem}, nil
		} else {
			return &GoArrayRetyper{Len: x.Len, Elem: elem}, nil
		}
	}
}

func (ctx *generalizeCtx) GeneralizeSlice(x, y *GoSlice) (_ GoRetyper, err error) {
	var elem GoRetyper
	if elem, err = ctx.Generalize(x.Elem, y.Elem); err != nil {
		return nil, err
	} else {
		return &GoSliceRetyper{Elem: elem}, nil
	}
}

func (ctx *generalizeCtx) GeneralizeArraySlice(x *GoArray, y *GoSlice) (_ GoRetyper, err error) {
	var elem GoRetyper
	if elem, err = ctx.Generalize(x.Elem, y.Elem); err != nil {
		return nil, err
	} else {
		return &GoSliceRetyper{Elem: elem}, nil
	}
}

func (ctx *generalizeCtx) GeneralizeStructure(x, y GoStructure) (_ GoRetyper, err error) {
	xy := &GoStructRetyper{Span: ctx.Span}
	field := map[string][]*GoField{}
	xfield, yfield := x.StructureField(), y.StructureField()
	for _, f := range xfield {
		field[f.KoName()] = append(field[f.KoName()], f)
	}
	for _, f := range yfield {
		field[f.KoName()] = append(field[f.KoName()], f)
	}
	for _, f := range field {
		switch len(f) {
		case 1:
			emptyField := BuildGoField(ctx.Span, f[0].KoName(), NewGoEmpty(ctx.Span), f[0].IsMonadic())
			if fieldRetyper, err := ctx.GeneralizeField(f[0], emptyField); err != nil {
				return nil, ctx.Span.Errorf(err, "cannot generalize %s against empty", Sprint(f[0]))
			} else {
				xy.Field = append(xy.Field, fieldRetyper)
			}
		case 2:
			if fieldRetyper, err := ctx.GeneralizeField(f[0], f[1]); err != nil {
				return nil, ctx.Span.Errorf(err, "cannot generalize %s and %s", Sprint(f[0]), Sprint(f[1]))
			} else {
				xy.Field = append(xy.Field, fieldRetyper)
			}
		default:
			panic("o")
		}
	}
	SortGoRetyperField(xy.Field)
	return xy, nil
}

func (ctx *generalizeCtx) GeneralizeField(f, g *GoField) (*GoStructRetyperField, error) {
	if fg, err := ctx.Generalize(f.Type, g.Type); err != nil {
		return nil, err
	} else {
		switch { // monadic transfer
		case f.IsMonadic() && g.IsMonadic():
			return &GoStructRetyperField{Name: f.KoName(), Retyper: fg, Monadic: true}, nil
		case !f.IsMonadic() && !g.IsMonadic():
			return &GoStructRetyperField{Name: f.KoName(), Retyper: fg, Monadic: false}, nil
		default:
			panic("o")
		}
	}
}

func (ctx *generalizeCtx) GeneralizeBuiltin(x, y *GoBuiltin) (xy GoRetyper, err error) {
	switch x.TypeID() {
	case GoBool.TypeID():
		switch y.TypeID() {
		case GoBool.TypeID():
			return &GoTypeRetyper{GoBool}, nil
		}
	case GoString.TypeID():
		switch y.TypeID() {
		case GoString.TypeID():
			return &GoTypeRetyper{GoString}, nil
		}
	case GoUintptr.TypeID():
		switch y.TypeID() {
		case GoUintptr.TypeID():
			return &GoTypeRetyper{GoUintptr}, nil
		}
	case GoUnsafePointer.TypeID():
		switch y.TypeID() {
		case GoUnsafePointer.TypeID():
			return &GoTypeRetyper{GoUnsafePointer}, nil
		}
	case GoInt8.TypeID():
		switch y.TypeID() {
		case GoInt8.TypeID(): // equal
			return &GoTypeRetyper{GoInt8}, nil
		case GoInt16.TypeID(): // larger
			return &GoTypeRetyper{GoInt16}, nil
		case GoInt32.TypeID(): // larger
			return &GoTypeRetyper{GoInt32}, nil
		case GoInt64.TypeID(): // larger
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt.TypeID(): // larger
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoInt16.TypeID():
		switch y.TypeID() {
		case GoInt8.TypeID(): // smaller
			return &GoTypeRetyper{GoInt16}, nil
		case GoInt16.TypeID(): // equal
			return &GoTypeRetyper{GoInt16}, nil
		case GoInt32.TypeID(): // larger
			return &GoTypeRetyper{GoInt32}, nil
		case GoInt64.TypeID(): // larger
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt.TypeID(): // larger
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoInt32.TypeID():
		switch y.TypeID() {
		case GoInt8.TypeID(): // smaller
			return &GoTypeRetyper{GoInt32}, nil
		case GoInt16.TypeID(): // smaller
			return &GoTypeRetyper{GoInt32}, nil
		case GoInt32.TypeID(): // equal
			return &GoTypeRetyper{GoInt32}, nil
		case GoInt64.TypeID(): // larger
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt.TypeID(): // larger
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoInt64.TypeID():
		switch y.TypeID() {
		case GoInt8.TypeID(): // smaller
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt16.TypeID(): // smaller
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt32.TypeID(): // smaller
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt64.TypeID(): // equal
			return &GoTypeRetyper{GoInt64}, nil
		case GoInt.TypeID(): // larger
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoInt.TypeID(): // Int is considered more general than Int64, categorically.
		switch y.TypeID() {
		case GoInt8.TypeID(), GoInt16.TypeID(), GoInt32.TypeID(), GoInt64.TypeID(): // smaller
			return &GoTypeRetyper{GoInt}, nil
		case GoInt.TypeID(): // equal
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoUint8.TypeID():
		switch y.TypeID() {
		case GoUint8.TypeID(): // equal
			return &GoTypeRetyper{GoUint8}, nil
		case GoUint16.TypeID(): // larger
			return &GoTypeRetyper{GoUint16}, nil
		case GoUint32.TypeID(): // larger
			return &GoTypeRetyper{GoUint32}, nil
		case GoUint64.TypeID(): // larger
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint.TypeID(): // larger
			return &GoTypeRetyper{GoUint}, nil
		}
	case GoUint16.TypeID():
		switch y.TypeID() {
		case GoUint8.TypeID(): // smaller
			return &GoTypeRetyper{GoUint16}, nil
		case GoUint16.TypeID(): // equal
			return &GoTypeRetyper{GoUint16}, nil
		case GoUint32.TypeID(): // larger
			return &GoTypeRetyper{GoUint32}, nil
		case GoUint64.TypeID(): // larger
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint.TypeID(): // larger
			return &GoTypeRetyper{GoInt}, nil
		}
	case GoUint32.TypeID():
		switch y.TypeID() {
		case GoUint8.TypeID(): // smaller
			return &GoTypeRetyper{GoUint32}, nil
		case GoUint16.TypeID(): // smaller
			return &GoTypeRetyper{GoUint32}, nil
		case GoUint32.TypeID(): // equal
			return &GoTypeRetyper{GoUint32}, nil
		case GoUint64.TypeID(): // larger
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint.TypeID(): // larger
			return &GoTypeRetyper{GoUint}, nil
		}
	case GoUint64.TypeID():
		switch y.TypeID() {
		case GoUint8.TypeID(): // smaller
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint16.TypeID(): // smaller
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint32.TypeID(): // smaller
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint64.TypeID(): // equal
			return &GoTypeRetyper{GoUint64}, nil
		case GoUint.TypeID(): // larger
			return &GoTypeRetyper{GoUint}, nil
		}
	case GoUint.TypeID():
		switch y.TypeID() {
		case GoUint8.TypeID(), GoUint16.TypeID(), GoUint32.TypeID(), GoUint64.TypeID(): // smaller
			return &GoTypeRetyper{GoUint}, nil
		case GoUint.TypeID(): // equal
			return &GoTypeRetyper{GoUint}, nil
		}
	case GoFloat32.TypeID():
		switch y.TypeID() {
		case GoFloat32.TypeID(): // equal
			return &GoTypeRetyper{GoFloat32}, nil
		case GoFloat64.TypeID(): // larger
			return &GoTypeRetyper{GoFloat64}, nil
		}
	case GoFloat64.TypeID():
		switch y.TypeID() {
		case GoFloat32.TypeID(): // smaller
			return &GoTypeRetyper{GoFloat64}, nil
		case GoFloat64.TypeID(): // equal
			return &GoTypeRetyper{GoFloat64}, nil
		}
	case GoComplex64.TypeID():
		switch y.TypeID() {
		case GoComplex64.TypeID(): // equal
			return &GoTypeRetyper{GoComplex64}, nil
		case GoComplex128.TypeID(): // larger
			return &GoTypeRetyper{GoComplex128}, nil
		}
	case GoComplex128.TypeID():
		switch y.TypeID() {
		case GoComplex64.TypeID(): // smaller
			return &GoTypeRetyper{GoComplex128}, nil
		case GoComplex128.TypeID(): // equal
			return &GoTypeRetyper{GoComplex128}, nil
		}
	}
	return nil, ctx.Span.Errorf(nil, "builtins %s and %s cannot be generalized", Sprint(x), Sprint(y))
}
