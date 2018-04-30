package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

var (
	BasicTrue  = BasicSymbol{true}
	BasicFalse = BasicSymbol{false}
)

func MakeBasicSymbol(span *Span, v interface{}) BasicSymbol {
	return Deconstruct(span, reflect.ValueOf(v)).(BasicSymbol)
}

func BasicByteSymbol(i byte) Symbol {
	return BasicSymbol{i}
}

func BasicInt32Symbol(i int32) Symbol {
	return BasicSymbol{i}
}

func BasicInt64Symbol(i int64) Symbol {
	return BasicSymbol{i}
}

func BasicStringSymbol(s string) Symbol {
	return BasicSymbol{s}
}

func AsBasicString(sym Symbol) (value string, ok bool) {
	basic, ok := sym.(BasicSymbol)
	if !ok {
		return "", false
	}
	s, ok := basic.Value.(string)
	return s, ok
}

func AsBasicBool(sym Symbol) (value bool, ok bool) {
	basic, ok := sym.(BasicSymbol)
	if !ok {
		return false, false
	}
	b, ok := basic.Value.(bool)
	return b, ok
}

type BasicSymbol struct {
	Value interface{} `ko:"name=value"`
}

func (basic BasicSymbol) Disassemble(span *Span) (*pb.Symbol, error) {
	dis := &pb.SymbolBasic{}
	switch u := basic.Value.(type) {
	case bool:
		dis.Basic = &pb.SymbolBasic_Bool{Bool: u}
	case string:
		dis.Basic = &pb.SymbolBasic_String_{String_: u}
	case int8:
		dis.Basic = &pb.SymbolBasic_Int8{Int8: int32(u)}
	case int16:
		dis.Basic = &pb.SymbolBasic_Int16{Int16: int32(u)}
	case int32:
		dis.Basic = &pb.SymbolBasic_Int32{Int32: u}
	case int64:
		dis.Basic = &pb.SymbolBasic_Int64{Int64: u}
	case uint8:
		dis.Basic = &pb.SymbolBasic_Uint8{Uint8: uint32(u)}
	case uint16:
		dis.Basic = &pb.SymbolBasic_Uint16{Uint16: uint32(u)}
	case uint32:
		dis.Basic = &pb.SymbolBasic_Uint32{Uint32: u}
	case uint64:
		dis.Basic = &pb.SymbolBasic_Uint64{Uint64: u}
	case float32:
		dis.Basic = &pb.SymbolBasic_Float32{Float32: u}
	case float64:
		dis.Basic = &pb.SymbolBasic_Float64{Float64: u}
	default:
		panic("o")
	}
	return &pb.Symbol{
		Symbol: &pb.Symbol_Basic{Basic: dis},
	}, nil
}

func (basic BasicSymbol) GoValue() reflect.Value {
	return reflect.ValueOf(basic.Value)
}

func (basic BasicSymbol) String() string {
	return Sprint(basic)
}

func (basic BasicSymbol) Equal(span *Span, sym Symbol) bool {
	if other, ok := sym.(BasicSymbol); ok {
		return basic.Value == other.Value
	} else {
		return false
	}
}

func (basic BasicSymbol) Hash(span *Span) ID {
	return InterfaceID(basic.Value)
}

func (basic BasicSymbol) ConvertTo(span *Span, to BasicType) (BasicSymbol, error) {
	v := reflect.ValueOf(basic.Value)
	if v.Type().ConvertibleTo(to.GoType()) {
		return BasicSymbol{v.Convert(to.GoType()).Interface()}, nil
	} else {
		return BasicSymbol{}, span.Errorf(nil, "cannot convert %v into %v", basic, to)
	}
}

func (basic BasicSymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return singletonSeries(basic)
}

func (basic BasicSymbol) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to basic")
}

func (basic BasicSymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return basic, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "basic value %v cannot be selected into", basic)
	}
}

func (basic BasicSymbol) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "basic value %v cannot be augmented", basic)
}

func (basic BasicSymbol) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "basic value %v cannot be invoked", basic)
}

func (basic BasicSymbol) Type() Type {
	return BasicFromKind(reflect.TypeOf(basic.Value).Kind())
}

func (basic BasicSymbol) Splay() Tree {
	return GoValue{reflect.ValueOf(basic.Value)}
}

func BasicFromKind(kind reflect.Kind) BasicType {
	switch kind {
	case reflect.Bool:
		return BasicBool
	case reflect.String:
		return BasicString
	case reflect.Int8:
		return BasicInt8
	case reflect.Int16:
		return BasicInt16
	case reflect.Int32:
		return BasicInt32
	case reflect.Int64:
		return BasicInt64
	case reflect.Uint8:
		return BasicUint8
	case reflect.Uint16:
		return BasicUint16
	case reflect.Uint32:
		return BasicUint32
	case reflect.Uint64:
		return BasicUint64
	case reflect.Float32:
		return BasicFloat32
	case reflect.Float64:
		return BasicFloat64
	}
	panic("o")
}

type BasicType int

func (BasicType) IsType() {}

func (basic BasicType) Splay() Tree {
	return NoQuote{basic.String()}
}

var (
	goBool    = reflect.TypeOf(bool(true))
	goString  = reflect.TypeOf(string(""))
	goInt8    = reflect.TypeOf(int8(0))
	goInt16   = reflect.TypeOf(int16(0))
	goInt32   = reflect.TypeOf(int32(0))
	goInt64   = reflect.TypeOf(int64(0))
	goUint8   = reflect.TypeOf(uint8(0))
	goUint16  = reflect.TypeOf(uint16(0))
	goUint32  = reflect.TypeOf(uint32(0))
	goUint64  = reflect.TypeOf(uint64(0))
	goFloat32 = reflect.TypeOf(float32(0.0))
	goFloat64 = reflect.TypeOf(float64(0.0))
)

func (basic BasicType) GoType() reflect.Type {
	switch basic {
	case BasicBool:
		return goBool
	case BasicString:
		return goString
	case BasicInt8:
		return goInt8
	case BasicInt16:
		return goInt16
	case BasicInt32:
		return goInt32
	case BasicInt64:
		return goInt64
	case BasicUint8:
		return goUint8
	case BasicUint16:
		return goUint16
	case BasicUint32:
		return goUint32
	case BasicUint64:
		return goUint64
	case BasicFloat32:
		return goFloat32
	case BasicFloat64:
		return goFloat64
	}
	panic("o")
}

func (basic BasicType) String() string {
	switch basic {
	case BasicInvalid:
		return "Invalid"
	case BasicBool:
		return "Bool"
	case BasicString:
		return "String"
	case BasicInt8:
		return "Int8"
	case BasicInt16:
		return "Int16"
	case BasicInt32:
		return "Int32"
	case BasicInt64:
		return "Int64"
	case BasicUint8:
		return "Uint8"
	case BasicUint16:
		return "Uint16"
	case BasicUint32:
		return "Uint32"
	case BasicUint64:
		return "Uint64"
	case BasicFloat32:
		return "Float32"
	case BasicFloat64:
		return "Float64"
	}
	panic("o")
}

const (
	BasicInvalid BasicType = iota
	BasicBool
	BasicString
	BasicInt8
	BasicInt16
	BasicInt32
	BasicInt64
	BasicUint8
	BasicUint16
	BasicUint32
	BasicUint64
	BasicFloat32
	BasicFloat64
)

func unifyBasic(x, y BasicType) (BasicType, bool) {
	switch x {
	case BasicBool:
		switch y {
		case BasicBool:
			return BasicBool, true
		}
	case BasicString:
		switch y {
		case BasicString:
			return BasicString, true
		}
	case BasicInt8: // signed integers
		switch y {
		case BasicInt8:
			return BasicInt8, true
		case BasicInt16:
			return BasicInt16, true
		case BasicInt32:
			return BasicInt32, true
		case BasicInt64:
			return BasicInt64, true
		}
	case BasicInt16:
		switch y {
		case BasicInt8:
			return BasicInt16, true
		case BasicInt16:
			return BasicInt16, true
		case BasicInt32:
			return BasicInt32, true
		case BasicInt64:
			return BasicInt64, true
		}
	case BasicInt32:
		switch y {
		case BasicInt8:
			return BasicInt32, true
		case BasicInt16:
			return BasicInt32, true
		case BasicInt32:
			return BasicInt32, true
		case BasicInt64:
			return BasicInt64, true
		}
	case BasicInt64:
		switch y {
		case BasicInt8:
			return BasicInt64, true
		case BasicInt16:
			return BasicInt64, true
		case BasicInt32:
			return BasicInt64, true
		case BasicInt64:
			return BasicInt64, true
		}
	case BasicUint8: // unsigned integers
		switch y {
		case BasicUint8:
			return BasicUint8, true
		case BasicUint16:
			return BasicUint16, true
		case BasicUint32:
			return BasicUint32, true
		case BasicUint64:
			return BasicUint64, true
		}
	case BasicUint16:
		switch y {
		case BasicUint8:
			return BasicUint16, true
		case BasicUint16:
			return BasicUint16, true
		case BasicUint32:
			return BasicUint32, true
		case BasicUint64:
			return BasicUint64, true
		}
	case BasicUint32:
		switch y {
		case BasicUint8:
			return BasicUint32, true
		case BasicUint16:
			return BasicUint32, true
		case BasicUint32:
			return BasicUint32, true
		case BasicUint64:
			return BasicUint64, true
		}
	case BasicUint64:
		switch y {
		case BasicUint8:
			return BasicUint64, true
		case BasicUint16:
			return BasicUint64, true
		case BasicUint32:
			return BasicUint64, true
		case BasicUint64:
			return BasicUint64, true
		}
	case BasicFloat32: // floating-point
		switch y {
		case BasicFloat32:
			return BasicFloat32, true
		case BasicFloat64:
			return BasicFloat64, true
		}
	case BasicFloat64:
		switch y {
		case BasicFloat32:
			return BasicFloat64, true
		case BasicFloat64:
			return BasicFloat64, true
		}
	}
	return BasicInvalid, false
}

func IsBasicKind(s Symbol, kind reflect.Kind) bool {
	if b, ok := s.Type().(BasicType); !ok {
		return false
	} else {
		return b.Kind() == kind
	}
}

func (b BasicType) Kind() reflect.Kind {
	switch b {
	case BasicString:
		return reflect.String
	case BasicBool:
		return reflect.Bool
	case BasicInt8:
		return reflect.Int8
	case BasicInt16:
		return reflect.Int16
	case BasicInt32:
		return reflect.Int32
	case BasicInt64:
		return reflect.Int64
	case BasicUint8:
		return reflect.Uint8
	case BasicUint16:
		return reflect.Uint16
	case BasicUint32:
		return reflect.Uint32
	case BasicUint64:
		return reflect.Uint64
	case BasicFloat32:
		return reflect.Float32
	case BasicFloat64:
		return reflect.Float64
	}
	panic("o")
}
