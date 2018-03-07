package model

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func Assignable(span *Span, from, to GoType) bool {
	_, _, err := Assign(span, from, to)
	return err == nil
}

func Assign(span *Span, from, to GoType) (bridge Shaper, cached *AssignCache, err error) {
	ctx := NewAssignCtx(span)
	if bridge, err := ctx.Assign(from, to); err != nil {
		return nil, nil, err
	} else {
		return bridge, ctx.Flush(), nil
	}
}

func MustBeAssignIso(span *Span, from, to GoType) {
	ctx := NewAssignCtx(span)
	var err error
	if _, err = ctx.Assign(from, to); err != nil {
		panic("o")
	}
	if _, err = ctx.Assign(to, from); err != nil {
		panic("o")
	}
}

func MustAssign(span *Span, from, to GoType) {
	if _, _, err := Assign(span, from, to); err != nil {
		panic(err)
	}
}

func NewAssignCtx(span *Span) *AssignCtx {
	return &AssignCtx{
		Origin:   span,
		Known:    SpanCache(span),
		Learning: AssignCacheUnion(), // empty cache
	}
}

// AssignCtx
type AssignCtx struct {
	Origin   *Span        `ko:"name=origin"`
	Known    *AssignCache `ko:"name=known"`
	Learning *AssignCache `ko:"name=learning"`
}

func (ctx *AssignCtx) Flush() *AssignCache {
	return ctx.Learning
}

func (ctx *AssignCtx) Errorf(cause error, format string, arg ...interface{}) error {
	return ctx.Origin.ErrorfSkip(2, cause, format, arg...)
}

func (ctx *AssignCtx) Lookup(from, to GoType) Shaper {
	if result := ctx.Known.Lookup(from, to); result != nil {
		return result
	}
	return ctx.Learning.Lookup(from, to)
}

func (ctx *AssignCtx) Reserve(from, to GoType) *ReserveShaper {
	return ctx.Learning.Reserve(ctx.Origin, from, to)
}

func (ctx *AssignCtx) Unreserve(reserve *ReserveShaper) {
	ctx.Learning.Unreserve(reserve)
}

func (ctx *AssignCtx) Assign(from, to GoType) (bridge Shaper, err error) {
	if bridge, ok := ctx.assignFilterUnknown(from, to); ok {
		return bridge, nil
	}
	// this caching mechanism supports Ko's ability for cyclical go type assignment
	if bridge = ctx.Lookup(from, to); bridge != nil {
		return bridge, nil
	}
	reserve := ctx.Reserve(from, to)
	fromSimplified, fromSimplifier := Simplify(ctx.Origin, from) // XXX: not passing assign cache to Simplify
	toSimplified, toSimplifier := Simplify(ctx.Origin, to)
	if bridge, err = ctx.assignCross(fromSimplified, toSimplified); err != nil {
		ctx.Unreserve(reserve)
		return nil, err
	}
	solution := CompressShapers(
		ctx.Origin,
		fromSimplifier,
		bridge,
		toSimplifier.(ReversibleShaper).Reverse(),
	)
	reserve.Solution = solution // close type recursion loops
	return reserve, nil
}

// from and to are not GoAlias{█}, GoPtr{GoPtr{█}}, GoNeverNilPtr{█}, GoArray{1, █}
func (ctx *AssignCtx) assignCross(from, to GoType) (bridge Shaper, err error) {
	if from.TypeID() == to.TypeID() {
		return IdentityShaper(ctx.Origin, from), nil
	}
	switch w := to.(type) {
	case *GoInterface:
		return ctx.assignToInterface(from, w)
	}
	switch u := from.(type) {
	case *GoEmpty:
		switch v := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, v)
		case *GoPtr:
			return &NilShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
			}, nil
		case *GoSlice:
			return &NilShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
			}, nil
		case GoNumber:
			return &UneraseNumberShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
			}, nil
		default:
			return nil, ctx.Errorf(nil, "assigning %s to %s", from, to)
		}
	case *GoPtr:
		switch v := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, v)
		case *GoPtr:
			if bridge, err = ctx.Assign(u.Elem, v); err != nil {
				return nil, err
			}
			// XXX: test this shaper
			return &OptShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: u, To: v},
				IfNotNil: CompressShapers(
					ctx.Origin,
					&DerefShaper{
						Shaping: Shaping{Origin: ctx.Origin, From: u, To: u.Elem}, N: 1,
					},
					bridge,
				),
			}, nil
		case *GoSlice:
			return ctx.assignPtrZoom(u, v)
		default:
			return nil, ctx.Errorf(nil, "assigning optional %v to non-optional %v", from, to)
		}
	case *GoSlice:
		switch v := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, v)
		case *GoPtr:
			return ctx.assignToPtr(from, v)
		case *GoSlice:
			return ctx.assignMatrix(from, to)
		default:
			return ctx.assignMatrix(from, to)
		}
	case *GoArray:
		switch v := to.(type) {
		case *GoArray:
			return ctx.assignMatrix(from, v)
		case *GoSlice:
			return ctx.assignMatrix(from, v)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to non-array %v", from, to)
		}
	default:
		switch v := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, v)
		case *GoPtr:
			return ctx.assignToPtr(from, v)
		case *GoSlice:
			return ctx.assignNonPtrZoom(from, v)
		default:
			return ctx.assignMatrix(from, to)
		}
	}
	ctx.Origin.Panicf(nil, "assign cross: %s to %s", Sprint(from), Sprint(to))
	panic("o")
}

func (ctx *AssignCtx) assignToPtr(from GoType, toPtr *GoPtr) (bridge Shaper, err error) {
	if bridge, err = ctx.Assign(from, toPtr.Elem); err != nil {
		return nil, err
	}
	return CompressShapers(
		ctx.Origin,
		bridge,
		&RefShaper{
			Shaping: Shaping{Origin: ctx.Origin, From: toPtr.Elem, To: toPtr}, N: 1,
		},
	), nil
}

// P -> []...[]Q
func (ctx *AssignCtx) assignNonPtrZoom(from GoType, to *GoSlice) (bridge Shaper, err error) {
	// fmt.Printf("assignNonPtrZoom: %s to %s\n", Sprint(from), Sprint(to))
	unzoomable, zoom := ReduceZoom(ctx.Origin, to)
	if bridge, err = ctx.Assign(from, unzoomable); err != nil {
		return nil, err
	}
	return CompressShapers(ctx.Origin, bridge, zoom.(ReversibleShaper).Reverse()), nil
}

// *P -> []...[]Q
func (ctx *AssignCtx) assignPtrZoom(from *GoPtr, to *GoSlice) (bridge Shaper, err error) {
	// fmt.Printf("assignPtrZoom: %s to %s\n", Sprint(from), Sprint(to))
	unzoomable, zoom := ReduceZoom(ctx.Origin, to)
	if bridge, err = ctx.Assign(from.Elem, unzoomable); err != nil {
		return nil, err
	}
	// XXX: test this shaper
	return &OptShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
		IfNotNil: CompressShapers(
			ctx.Origin,
			&DerefShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: from, To: from.Elem}, N: 1,
			},
			bridge,
			zoom.(ReversibleShaper).Reverse(),
		),
	}, nil
}

func (ctx *AssignCtx) assignFilterUnknown(from, to GoType) (Shaper, bool) {
	_, fromIsUnknown := from.(Unknown)
	_, toIsUnknown := to.(Unknown)
	if !fromIsUnknown && !toIsUnknown {
		return nil, false
	}
	return &UnknownShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}, true
}

func (ctx *AssignCtx) assignMatrix(from, to GoType) (bridge Shaper, err error) {
	switch u := from.(type) {
	case *GoBuiltin:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoBuiltin:
			return ctx.assignBuiltinBuiltin(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoStruct:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoStruct:
			return ctx.assignStructStruct(u, w)
		case *GoMap:
			return ctx.assignStructMap(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoMap:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoMap:
			return ctx.assignMapMap(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoSlice:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoSlice:
			return ctx.assignSliceSlice(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoArray:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoArray:
			return ctx.assignArrayArray(u, w)
		case *GoSlice:
			return ctx.assignArraySlice(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoChan:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoChan:
			// return ctx.assignChanChan( u, w)
			return nil, ctx.Errorf(nil, "chan shaping not supported, assigning %v to %v", from, to)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoFunc:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoFunc:
			// return ctx.assignFuncFunc( u, w)
			return nil, ctx.Errorf(nil, "func shaping not supported, assigning %v to %v", from, to)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoInterface:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case GoNumber:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoBuiltin:
			return ctx.assignNumberBuiltin(u, w)
		case GoNumber:
			return ctx.assignNumberNumber(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoVariety:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoVariety:
			return ctx.assignVarietyVariety(u, w)
		default:
			return nil, ctx.Errorf(nil, "assigning %v to %v", from, to)
		}
	case *GoEmpty:
		switch w := to.(type) {
		case *GoEmpty:
			return ctx.assignToEmpty(u, w)
		case *GoPtr:
			return &NilShaper{
				Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
			}, nil
		default:
			return nil, ctx.Errorf(nil, "assigning %s to %s", Sprint(from), Sprint(to))
		}
	}
	panic("u")
}

func (ctx *AssignCtx) assignToEmpty(from GoType, to *GoEmpty) (Shaper, error) {
	return &IrreversibleEraseShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}, nil
}

func (ctx *AssignCtx) assignToInterface(from GoType, to *GoInterface) (Shaper, error) {
	// no need to check that inteface is restricted, go compiler will catch errors
	// if to.IsRestricted() {
	// 	return nil, ctx.Errorf(nil,"cannot assign to restricted interfaces %s", Sprint(to))
	// }
	return &ReShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}, nil
}

func (ctx *AssignCtx) assignBuiltinBuiltin(from *GoBuiltin, to *GoBuiltin) (Shaper, error) {
	convertShaper := &ConvertTypeShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}
	switch from.Kind {
	case reflect.Bool:
		switch to.Kind {
		case reflect.Bool:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.Int8:
		switch to.Kind {
		case reflect.Int8:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Int16:
		switch to.Kind {
		case reflect.Int16:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Int32, reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Int32:
		switch to.Kind {
		case reflect.Int32:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Int64:
		switch to.Kind {
		case reflect.Int64:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Int:
			return convertShaper, nil
		}
	case reflect.Int:
		switch to.Kind {
		case reflect.Int64:
			return convertShaper, nil
		case reflect.Int:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.Uint8:
		switch to.Kind {
		case reflect.Uint8:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			return convertShaper, nil
		case reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Uint16:
		switch to.Kind {
		case reflect.Uint16:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Uint32, reflect.Uint64, reflect.Uint:
			return convertShaper, nil
		case reflect.Int32, reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Uint32:
		switch to.Kind {
		case reflect.Uint32:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Uint64, reflect.Uint:
			return convertShaper, nil
		case reflect.Int64, reflect.Int:
			return convertShaper, nil
		}
	case reflect.Uint64:
		switch to.Kind {
		case reflect.Uint64:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Uint:
			return convertShaper, nil
		}
	case reflect.Uint:
		switch to.Kind {
		case reflect.Uint64:
			return convertShaper, nil
		case reflect.Uint:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.Uintptr:
		switch to.Kind {
		case reflect.Uintptr:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Uint, reflect.Uint64:
			return convertShaper, nil
		}
	case reflect.Float32:
		switch to.Kind {
		case reflect.Float32:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Float64:
			return convertShaper, nil
		}
	case reflect.Float64:
		switch to.Kind {
		case reflect.Float64:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.Complex64:
		switch to.Kind {
		case reflect.Complex64:
			return IdentityShaper(ctx.Origin, from), nil
		case reflect.Complex128:
			return convertShaper, nil
		}
	case reflect.Complex128:
		switch to.Kind {
		case reflect.Complex128:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.String:
		switch to.Kind {
		case reflect.String:
			return IdentityShaper(ctx.Origin, from), nil
		}
	case reflect.UnsafePointer:
		switch to.Kind {
		case reflect.UnsafePointer:
			return IdentityShaper(ctx.Origin, from), nil
		}
	}
	return nil, ctx.Errorf(nil, "assigning go builtin %v to %v", from.Kind, to.Kind)
}

func RenameMonadicForFunc(span *Span, passed GoStructure, fu *Func) *GoStruct {
	return RenameMonadic(span, passed, fu.Monadic)
}

func RenameMonadic(span *Span, passed GoStructure, name string) *GoStruct {
	passedField := passed.StructureField()
	field := make([]*GoField, len(passedField))
	for i, p := range passedField {
		if p.IsMonadic() {
			field[i] = &GoField{
				Comment: p.Comment,
				Name:    GoNameFor(name), // new ko name
				Type:    p.Type,
				Tag:     KoTags(name, true),
			}
		} else {
			field[i] = p
		}
	}
	return NewGoStruct(field...)
}

func (ctx *AssignCtx) assignStructStruct(from *GoStruct, to *GoStruct) (Shaper, error) {
	liftFrom, liftTo := GoFieldsOnPath(from.Field, nil), GoFieldsOnPath(to.Field, nil)
	if insertFieldShapers, err := ctx.assignInsert(liftFrom, liftTo); err != nil {
		return nil, err
	} else {
		return &StructStructShaper{
			Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
			Field:   insertFieldShapers.FieldShapers(),
		}, nil
	}
}

type InsertFieldShaper struct {
	FromTo *FromToGoField `ko:"name=fromTo"`
	Shaper Shaper         `ko:"name=shaper"`
}

func (ifs *InsertFieldShaper) FieldShaper() *FieldShaper {
	return &FieldShaper{
		From:   ifs.FromTo.From.Field.Name,
		To:     ifs.FromTo.To.Field.Name,
		Shaper: ifs.Shaper,
	}
}

func (ifs *InsertFieldShaper) VariationFieldShaper() *VariationFieldShaper {
	return &VariationFieldShaper{
		From:   ifs.FromTo.From,
		To:     ifs.FromTo.To,
		Shaper: ifs.Shaper,
	}
}

type InsertFieldShapers []*InsertFieldShaper

func (ifs InsertFieldShapers) FieldShapers() (fieldShaper []*FieldShaper) {
	fieldShaper = make([]*FieldShaper, len(ifs))
	for i, ifs := range ifs {
		fieldShaper[i] = ifs.FieldShaper()
	}
	return
}

func (ifs InsertFieldShapers) VariationFieldShapers() (fieldShaper []*VariationFieldShaper) {
	fieldShaper = make([]*VariationFieldShaper, len(ifs))
	for i, ifs := range ifs {
		fieldShaper[i] = ifs.VariationFieldShaper()
	}
	return
}

func (ctx *AssignCtx) assignInsert(from, to []*GoPathField) (fieldShapers InsertFieldShapers, err error) {
	alignedFields := AlignGoFields(ctx.Origin, from, to)
	for _, fromTo := range alignedFields {
		if fromToShaper, err := ctx.Assign(fromTo.From.Field.Type, fromTo.To.Field.Type); err != nil {
			return nil, ctx.Errorf(err, "assigning from field %s to field %s", Sprint(fromTo.From), Sprint(fromTo.To))
		} else {
			if _, toIsEmpty := fromTo.To.Field.Type.(*GoEmpty); !toIsEmpty { // if to-field is not an empty, render the assignment
				fieldShapers = append(
					fieldShapers,
					&InsertFieldShaper{FromTo: fromTo, Shaper: fromToShaper},
				)
			}
		}
	}
	return
}

func (ctx *AssignCtx) assignStructMap(from *GoStruct, to *GoMap) (Shaper, error) {
	if to.Key.TypeID() != GoString.TypeID() {
		return nil, ctx.Errorf(nil, "struct-receiving map key not a string (%v)", to.Key)
	}
	sms := &StructMapShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}
	for _, f := range from.Field {
		if f2g, err := ctx.Assign(f.Type, to.Value); err != nil {
			return nil, ctx.Errorf(err, "assigning from field %s of %v to map value %v", f.KoName(), from, to.Value)
		} else {
			sms.Field = append(sms.Field,
				&FieldShaper{From: f.Name, To: f.KoName(), Shaper: f2g},
			)
		}
	}
	return sms, nil
}

func (ctx *AssignCtx) assignMapMap(from *GoMap, to *GoMap) (Shaper, error) {
	key, err := ctx.Assign(from.Key, to.Key)
	if err != nil {
		return nil, ctx.Errorf(nil, "assigning key of %v to %v", from, to)
	}
	value, err := ctx.Assign(from.Value, to.Value)
	if err != nil {
		return nil, ctx.Errorf(nil, "assigning value of %v to %v", from, to)
	}
	return &MapMapShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
		Key:     key,
		Value:   value,
	}, nil
}

func (ctx *AssignCtx) assignSliceSlice(from *GoSlice, to *GoSlice) (Shaper, error) {
	elem, err := ctx.Assign(from.Elem, to.Elem)
	if err != nil {
		return nil, ctx.Errorf(err, "assigning element of %v to %v", from, to)
	}
	return &BatchShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
		Elem:    elem,
	}, nil
}

func (ctx *AssignCtx) assignArrayArray(from *GoArray, to *GoArray) (Shaper, error) {
	elem, err := ctx.Assign(from.Elem, to.Elem)
	if err != nil {
		return nil, ctx.Errorf(err, "assigning element of %v to %v", from, to)
	}
	return &BatchShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
		Elem:    elem,
	}, nil
}

func (ctx *AssignCtx) assignArraySlice(from *GoArray, to *GoSlice) (Shaper, error) {
	elem, err := ctx.Assign(from.Elem, to.Elem)
	if err != nil {
		return nil, ctx.Errorf(err, "assigning element of %v to %v", from, to)
	}
	return &BatchShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
		Elem:    elem,
	}, nil
}

func (ctx *AssignCtx) assignNumberBuiltin(from GoNumber, to *GoBuiltin) (shaper Shaper, err error) {
	// check that go convertsion and assignment is possible without loss of value
	f0 := from.Value()
	fv, tv := NewNumberValue(from), to.NewValue()
	if !fv.Type().ConvertibleTo(tv.Type()) || !tv.Type().ConvertibleTo(fv.Type()) {
		return nil, ctx.Errorf(nil, "literal %s not convertible to builtin %s", Sprint(from), Sprint(to))
	}
	tv.Set(fv.Convert(tv.Type()))
	f1 := tv.Convert(fv.Type()).Interface()
	if f0 != f1 {
		return nil, ctx.Errorf(nil, "literal %s does not fit into builtin %s", Sprint(from), Sprint(to))
	}
	return &NumberShaper{
		Shaping: Shaping{Origin: ctx.Origin, From: from, To: to},
	}, nil
}

func (ctx *AssignCtx) assignNumberNumber(from GoNumber, to GoNumber) (shaper Shaper, err error) {
	if !NumberEqual(from, to) {
		return nil, ctx.Errorf(nil, "assigning number %s to %s", Sprint(from), Sprint(to))
	}
	return IdentityShaper(ctx.Origin, from), nil
}
