package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func Integrate(span *Span, s Symbol, t reflect.Type) (reflect.Value, error) {
	ctx := &typingCtx{Span: span}
	return ctx.Integrate(s, t)
}

func (ctx *typingCtx) Integrate(s Symbol, t reflect.Type) (reflect.Value, error) {
	if r, err := ctx.integrateNamed(s, t); err == nil { // try
		return r, nil
	}
	// if s is named, deconstruct its go value
	if named, ok := s.(*NamedSymbol); ok {
		if dec, err := ctx.DeconstructKind(named.Value); err != nil {
			return reflect.Value{}, err
		} else {
			s = dec
		}
	}
	return ctx.IntegrateKind(s, t)
}

func (ctx *typingCtx) integrateNamed(s Symbol, t reflect.Type) (reflect.Value, error) {
	tName := TypeName(t)
	if tName == "" {
		return reflect.Value{}, ctx.Errorf(nil, "to-type is not named")
	}
	if t.Kind() == reflect.Interface {
		// interface named types are handled in Integrate
		return reflect.Value{}, ctx.Errorf(nil, "to-type is an interface")
	}
	sNamed, ok := s.(*NamedSymbol)
	if !ok {
		return reflect.Value{}, ctx.Errorf(nil, "from-symbol is not named")
	}
	sGoType := sNamed.GoType()
	if sGoType == t {
		return sNamed.Value, nil
	} else {
		return reflect.Value{}, ctx.Errorf(nil,
			"cannot integrate named type %s to named type %s",
			TypeName(sGoType), TypeName(t),
		)
	}
}

func (ctx *typingCtx) IntegrateKind(s Symbol, t reflect.Type) (reflect.Value, error) {
	switch t.Kind() {
	case reflect.Invalid:
		panic("o")
	case reflect.String:
		if g, err := ctx.IntegrateBasic(s, t); err == nil {
			return g, nil
		} else if blob, ok := s.(*BlobSymbol); ok { // blob -> string
			return blob.Value.Convert(t), nil
		}
	case reflect.Bool:
		return ctx.IntegrateBasic(s, t)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ctx.IntegrateBasicBits(s, t)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ctx.IntegrateBasicBits(s, t)
	case reflect.Float32, reflect.Float64:
		return ctx.IntegrateBasicBits(s, t)
	case reflect.Uintptr: // defer to IntegrateFrom Named/Opaque
	case reflect.Complex64: // defer to IntegrateFrom Named/Opaque
	case reflect.Complex128: // defer to IntegrateFrom Named/Opaque
	case reflect.Array: // defer to IntegrateFrom Named/Opaque
	case reflect.Chan: // defer to IntegrateFrom Named/Opaque
	case reflect.UnsafePointer: // defer to IntegrateFrom Named/Opaque
	case reflect.Func: // defer to IntegrateFrom Named/Opaque
	case reflect.Interface: // defer to IntegrateFrom Named/Opaque
	case reflect.Map:
		if IsEmptySymbol(s) {
			return reflect.Zero(t), nil
		} else if extracted, err := ExtractMap(ctx.Span, s, t); err == nil {
			return ctx.IntegrateKind(extracted, t) // try again
		} else if ms, ok := s.(*MapSymbol); ok {
			return ctx.IntegrateMapMap(ms, t)
		} else {
			// defer to IntegrateFrom Named/Opaque
		}
	case reflect.Ptr:
		if IsEmptySymbol(s) {
			return reflect.Zero(t), nil
		} else {
			if elem, err := ctx.Integrate(s, t.Elem()); err != nil {
				return reflect.Value{}, err
			} else {
				var u reflect.Value
				if elem.CanAddr() {
					u = elem
				} else {
					u = reflect.New(elem.Type()).Elem()
					u.Set(elem)
				}
				w := reflect.New(t).Elem()
				w.Set(u.Addr())
				return w, nil
			}
		}
	case reflect.Slice:
		if IsEmptySymbol(s) {
			return reflect.Zero(t), nil
		} else if t == byteSliceType {
			if blob, isBlob := s.(*BlobSymbol); isBlob { // blob -> []byte
				return blob.Value, nil
			} else if s, isString := AsBasicString(s); isString { // string -> []byte
				return reflect.ValueOf(s).Convert(t), nil
			}
		}
		return ctx.IntegrateSlice(s, t)
	case reflect.Struct: // catch missing fields
		if IsEmptySymbol(s) {
			return ctx.IntegrateStruct(MakeStructSymbol(nil), t)
		} else if ss, ok := s.(*StructSymbol); ok {
			return ctx.IntegrateStruct(ss, t)
		} else if ms, ok := s.(*MapSymbol); ok {
			return ctx.IntegrateMapStruct(ms, t)
		}
	}
	//
	switch u := s.(type) {
	case *OpaqueSymbol:
		return ctx.IntegrateFromOpaque(u, t)
	case *NamedSymbol:
		return ctx.IntegrateFromNamed(u, t)
	}
	return reflect.Value{}, ctx.Errorf(nil, "cannot integrate %s into %v", Sprint(s), t)
}

var TypeOfInterface = reflect.TypeOf((*interface{})(nil)).Elem()

func (ctx *typingCtx) IntegrateFromOpaque(u *OpaqueSymbol, t reflect.Type) (reflect.Value, error) {
	if u.GoType().AssignableTo(t) {
		if u.Value.CanAddr() {
			return u.Value, nil
		} else {
			w := reflect.New(t).Elem()
			w.Set(u.Value)
			return w, nil
		}
	} else {
		return reflect.Value{},
			ctx.Errorf(nil, "cannot integrate opaque type %v into go type %v", u.Type(), t)
	}
}

func (ctx *typingCtx) IntegrateFromNamed(u *NamedSymbol, t reflect.Type) (reflect.Value, error) {
	goType := u.GoType()
	if goType.AssignableTo(t) { // T -> interface
		return u.Value, nil
	} else if reflect.PtrTo(goType).AssignableTo(t) { // *T -> interface
		if u.Value.CanAddr() {
			return u.Value.Addr(), nil
		} else {
			w := reflect.New(goType)
			w.Elem().Set(u.Value)
			return w, nil
		}
	} else {
		return reflect.Value{},
			ctx.Errorf(nil, "cannot integrate named type %v into go type %v", u.GoType(), t)
	}
}

func (ctx *typingCtx) IntegrateBasic(s Symbol, t reflect.Type) (reflect.Value, error) {
	if basic, ok := s.(BasicSymbol); !ok {
		return reflect.Value{}, ctx.Errorf(nil, "value %v is not basic, cannot integrate to %v", s, t)
	} else {
		stype := reflect.TypeOf(basic.Value)
		if stype.ConvertibleTo(t) {
			return reflect.ValueOf(basic.Value).Convert(t), nil
		} else {
			return reflect.Value{}, ctx.Errorf(nil, "value %v (of type %v) is not convertible to %v", s, s.Type(), t)
		}
	}
}

func (ctx *typingCtx) IntegrateBasicBits(s Symbol, t reflect.Type) (reflect.Value, error) {
	if basic, ok := s.(BasicSymbol); !ok {
		return reflect.Value{}, ctx.Errorf(nil, "value %v is not basic, cannot integrate to %v", s, t)
	} else {
		stype := reflect.TypeOf(basic.Value)
		if stype.ConvertibleTo(t) && stype.Bits() <= t.Bits() {
			return reflect.ValueOf(basic.Value).Convert(t), nil
		} else {
			return reflect.Value{}, ctx.Errorf(nil, "value %v (of type %v) is not convertible to %v", s, s.Type(), t)
		}
	}
}

func (ctx *typingCtx) IntegrateSlice(s Symbol, t reflect.Type) (reflect.Value, error) {
	ss := s.LiftToSeries(ctx.Span)
	elems := make([]reflect.Value, len(ss.Elem))
	ctx2 := ctx.Refine("()")
	for i, symElem := range ss.Elem {
		if u, err := ctx2.Integrate(symElem, t.Elem()); err != nil {
			return reflect.Value{}, err
		} else {
			elems[i] = u
		}
	}
	w := reflect.MakeSlice(t, len(elems), len(elems))
	for i, elem := range elems {
		w.Index(i).Set(elem)
	}
	return w, nil
}

func (ctx *typingCtx) IntegrateStruct(ss *StructSymbol, t reflect.Type) (reflect.Value, error) {
	w := reflect.New(t).Elem()
	for i := 0; i < t.NumField(); i++ {
		toField := t.Field(i)
		if from := FindIntegrationField(ss, toField); from == nil {
			if !gate.StructFieldIsOptional(toField) {
				return reflect.Value{},
					ctx.Errorf(nil, "go field %s in %v is required, not found in %v", toField.Name, t, ss)
			}
		} else {
			if u, err := ctx.Refine(toField.Name).Integrate(from.Value, toField.Type); err != nil {
				return reflect.Value{}, err
			} else {
				w.Field(i).Set(u)
			}
		}
	}
	return w, nil
}

func FindIntegrationField(from *StructSymbol, to reflect.StructField) *FieldSymbol {
	name, hasKoName := gate.StructFieldKoProtoGoName(to)
	if !hasKoName {
		return nil
	}
	if gate.IsStructFieldMonadic(to) {
		if monadicField := from.FindMonadic(); monadicField != nil {
			return monadicField
		}
	}
	return from.FindName(name)
}

func (ctx *typingCtx) IntegrateMapMap(ms *MapSymbol, t reflect.Type) (reflect.Value, error) {
	if t.Key() != typeOfString {
		return reflect.Value{}, ctx.Errorf(nil, "map %v cannot integrate into go map %v", ms, t)
	}
	w := reflect.MakeMap(t)
	for k, vsym := range ms.Map {
		if wsym, err := ctx.Refine(k).Integrate(vsym, t.Elem()); err != nil {
			return reflect.Value{},
				ctx.Errorf(err,
					"map value %v (type %v) cannot integrate into go map value %v",
					vsym, vsym.Type(), t.Elem(),
				)
		} else {
			w.SetMapIndex(reflect.ValueOf(k), wsym)
		}
	}
	return w, nil
}

func (ctx *typingCtx) IntegrateMapStruct(ms *MapSymbol, t reflect.Type) (reflect.Value, error) {
	w := reflect.New(t).Elem()
	for i := 0; i < t.NumField(); i++ {
		toField := t.Field(i)
		if fromValue := FindIntegrationKey(ms, toField); fromValue == nil {
			if !gate.StructFieldIsOptional(toField) {
				return reflect.Value{},
					ctx.Errorf(nil, "go field %s in %v is required, not found in %v", toField.Name, t, ms)
			}
		} else {
			if u, err := ctx.Refine(toField.Name).Integrate(fromValue, toField.Type); err != nil {
				return reflect.Value{}, err
			} else {
				w.Field(i).Set(u)
			}
		}
	}
	return w, nil
}

func FindIntegrationKey(from *MapSymbol, to reflect.StructField) Symbol {
	name, hasKoName := gate.StructFieldKoProtoGoName(to)
	if !hasKoName {
		return nil
	}
	if gate.IsStructFieldMonadic(to) {
		if monadic, ok := from.Map[""]; ok {
			return monadic
		}
	}
	return from.Map[name]
}
