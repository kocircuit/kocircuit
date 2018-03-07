package model

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// GoBuiltin captures a Go builtin type.
type GoBuiltin struct {
	ID   string       `ko:"name=id"`
	Kind reflect.Kind `ko:"name=kind"`
}

func NewGoBuiltin(kind reflect.Kind) *GoBuiltin {
	return &GoBuiltin{ID: Mix("builtin", kind.String()), Kind: kind}
}

func (builtin *GoBuiltin) TypeID() string { return builtin.ID }

func (builtin *GoBuiltin) Doc() string { return "" }

func (builtin *GoBuiltin) String() string { return Sprint(builtin) }

func (builtin *GoBuiltin) Sketch(ctx *GoSketchCtx) interface{} {
	return builtin.Kind.String()
}

func (builtin *GoBuiltin) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: builtin
func (builtin *GoBuiltin) RenderDef(_ GoFileContext) string { return builtin.Kind.String() }

// RenderRef returns a type reference of the form: Builtin
func (builtin *GoBuiltin) RenderRef(_ GoFileContext) string { return builtin.Kind.String() }

// RenderZero returns a zero value for the respective Go builtin type.
func (builtin *GoBuiltin) RenderZero(_ GoFileContext) string {
	switch builtin.Kind {
	case reflect.Bool:
		return `false`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return `0`
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return `0`
	case reflect.Uintptr:
		return `0`
	case reflect.Float32, reflect.Float64:
		return `0.0`
	case reflect.Complex64, reflect.Complex128:
		return `0+0i`
	case reflect.String:
		return `""`
	case reflect.UnsafePointer:
		return `nil`
	}
	panic("unknown builtin kind")
}

func (builtin *GoBuiltin) ZeroType() reflect.Type {
	switch builtin.Kind {
	case reflect.Bool:
		return reflect.TypeOf(bool(false))
	case reflect.String:
		return reflect.TypeOf(string(""))
	case reflect.Int:
		return reflect.TypeOf(int(0))
	case reflect.Int8:
		return reflect.TypeOf(int8(0))
	case reflect.Int16:
		return reflect.TypeOf(int16(0))
	case reflect.Int32:
		return reflect.TypeOf(int32(0))
	case reflect.Int64:
		return reflect.TypeOf(int64(0))
	case reflect.Uint:
		return reflect.TypeOf(uint(0))
	case reflect.Uint8:
		return reflect.TypeOf(uint8(0))
	case reflect.Uint16:
		return reflect.TypeOf(uint16(0))
	case reflect.Uint32:
		return reflect.TypeOf(uint32(0))
	case reflect.Uint64:
		return reflect.TypeOf(uint64(0))
	case reflect.Uintptr:
		return reflect.TypeOf(uintptr(0))
	case reflect.Float32:
		return reflect.TypeOf(float32(0.0))
	case reflect.Float64:
		return reflect.TypeOf(float64(0.0))
	case reflect.Complex64:
		return reflect.TypeOf(complex64(0 + 0i))
	case reflect.Complex128:
		return reflect.TypeOf(complex128(0 + 0i))
	}
	panic("o")
}

func (builtin *GoBuiltin) NewValue() reflect.Value {
	return reflect.New(builtin.ZeroType()).Elem()
}

var (
	GoBool          = NewGoBuiltin(reflect.Bool)
	GoString        = NewGoBuiltin(reflect.String)
	GoInt           = NewGoBuiltin(reflect.Int)
	GoInt8          = NewGoBuiltin(reflect.Int8)
	GoInt16         = NewGoBuiltin(reflect.Int16)
	GoInt32         = NewGoBuiltin(reflect.Int32)
	GoInt64         = NewGoBuiltin(reflect.Int64)
	GoUint          = NewGoBuiltin(reflect.Uint)
	GoUint8         = NewGoBuiltin(reflect.Uint8)
	GoUint16        = NewGoBuiltin(reflect.Uint16)
	GoUint32        = NewGoBuiltin(reflect.Uint32)
	GoUint64        = NewGoBuiltin(reflect.Uint64)
	GoUintptr       = NewGoBuiltin(reflect.Uintptr)
	GoFloat32       = NewGoBuiltin(reflect.Float32)
	GoFloat64       = NewGoBuiltin(reflect.Float64)
	GoComplex64     = NewGoBuiltin(reflect.Complex64)
	GoComplex128    = NewGoBuiltin(reflect.Complex128)
	GoUnsafePointer = NewGoBuiltin(reflect.UnsafePointer)
)
