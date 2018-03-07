package model

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
)

type GoNumber interface {
	GoType
	Value() interface{}
	Builtin() *GoBuiltin
	NumberExpr() GoExpr
}

func NewNumberValue(num GoNumber) reflect.Value {
	v := num.Builtin().NewValue()
	v.Set(reflect.ValueOf(num.Value()))
	return v
}

var (
	GoTrue  = NewGoBoolNumber(true)
	GoFalse = NewGoBoolNumber(false)
)

type GoBoolNumber struct {
	ID     string `ko:"name=id"`
	Value_ bool   `ko:"name=value"`
}

func NewGoBoolNumber(v bool) *GoBoolNumber {
	return &GoBoolNumber{ID: MixScalar(v), Value_: v}
}

func (lit *GoBoolNumber) TypeID() string { return lit.ID }

func (lit *GoBoolNumber) NumberExpr() GoExpr { return &GoVerbatimExpr{lit.String()} }

func (lit *GoBoolNumber) Doc() string { return fmt.Sprintf("bool-number(%v)", lit.Value_) }

func (lit *GoBoolNumber) String() string { return fmt.Sprintf("%v", lit.Value_) }

func (lit *GoBoolNumber) Sketch(ctx *GoSketchCtx) interface{} { return lit.Value_ }

func (lit *GoBoolNumber) Value() interface{} { return lit.Value_ }

func (lit *GoBoolNumber) Builtin() *GoBuiltin { return GoBool }

func (lit *GoBoolNumber) Tag() []*GoTag { return nil }

func (lit *GoBoolNumber) RenderDef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderDef(fileCtx)
}

func (lit *GoBoolNumber) RenderRef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderRef(fileCtx)
}

func (lit *GoBoolNumber) RenderZero(fileCtx GoFileContext) string {
	return lit.NumberExpr().RenderExpr(fileCtx)
}

type GoIntegerNumber struct {
	ID     string `ko:"name=id"`
	Value_ int64  `ko:"name=value"`
}

func NewGoIntegerNumber(v int64) *GoIntegerNumber {
	return &GoIntegerNumber{ID: MixScalar(v), Value_: v}
}

func (lit *GoIntegerNumber) TypeID() string { return lit.ID }

func (lit *GoIntegerNumber) Negative() *GoIntegerNumber {
	return NewGoIntegerNumber(-lit.Value_)
}

func (lit *GoIntegerNumber) NumberExpr() GoExpr {
	return &GoVerbatimExpr{lit.String()}
}

func (lit *GoIntegerNumber) Doc() string { return fmt.Sprintf("integer-number(%v)", lit.Value_) }

func (lit *GoIntegerNumber) String() string { return fmt.Sprintf("int64(%d)", lit.Value_) }

func (lit *GoIntegerNumber) Sketch(ctx *GoSketchCtx) interface{} { return lit.Value_ }

func (lit *GoIntegerNumber) Value() interface{} { return lit.Value_ }

func (lit *GoIntegerNumber) Builtin() *GoBuiltin { return GoInt64 }

func (lit *GoIntegerNumber) Tag() []*GoTag { return nil }

func (lit *GoIntegerNumber) RenderDef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderDef(fileCtx)
}

func (lit *GoIntegerNumber) RenderRef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderRef(fileCtx)
}

func (lit *GoIntegerNumber) RenderZero(fileCtx GoFileContext) string {
	return lit.NumberExpr().RenderExpr(fileCtx)
}

type GoFloatNumber struct {
	ID     string  `ko:"name=id"`
	Value_ float64 `ko:"name=value"`
}

func NewGoFloatNumber(v float64) *GoFloatNumber {
	return &GoFloatNumber{ID: MixScalar(v), Value_: v}
}

func (lit *GoFloatNumber) TypeID() string { return lit.ID }

func (lit *GoFloatNumber) NumberExpr() GoExpr {
	return &GoVerbatimExpr{lit.String()}
}

func (lit *GoFloatNumber) Doc() string { return fmt.Sprintf("float-number(%v)", lit.Value_) }

func (lit *GoFloatNumber) String() string { return fmt.Sprintf("float64(%g)", lit.Value_) }

func (lit *GoFloatNumber) Sketch(ctx *GoSketchCtx) interface{} { return lit.Value_ }

func (lit *GoFloatNumber) Value() interface{} { return lit.Value_ }

func (lit *GoFloatNumber) Builtin() *GoBuiltin { return GoFloat64 }

func (lit *GoFloatNumber) Tag() []*GoTag { return nil }

func (lit *GoFloatNumber) RenderDef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderDef(fileCtx)
}

func (lit *GoFloatNumber) RenderRef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderRef(fileCtx)
}

func (lit *GoFloatNumber) RenderZero(fileCtx GoFileContext) string {
	return lit.NumberExpr().RenderExpr(fileCtx)
}

type GoStringNumber struct {
	ID     string `ko:"name=id"`
	Value_ string `ko:"name=value"`
}

func NewGoStringNumber(v string) *GoStringNumber {
	return &GoStringNumber{ID: MixScalar(v), Value_: v}
}

func (lit *GoStringNumber) TypeID() string { return lit.ID }

func (lit *GoStringNumber) NumberExpr() GoExpr {
	return &GoVerbatimExpr{lit.String()}
}

func (lit *GoStringNumber) Doc() string { return fmt.Sprintf("string-number(%v)", lit.Value_) }

func (lit *GoStringNumber) String() string { return fmt.Sprintf("%q", lit.Value_) }

func (lit *GoStringNumber) Sketch(ctx *GoSketchCtx) interface{} { return lit.Value_ }

func (lit *GoStringNumber) Value() interface{} { return lit.Value_ }

func (lit *GoStringNumber) Builtin() *GoBuiltin { return GoString }

func (lit *GoStringNumber) Tag() []*GoTag { return nil }

func (lit *GoStringNumber) RenderDef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderDef(fileCtx)
}

func (lit *GoStringNumber) RenderRef(fileCtx GoFileContext) string {
	return lit.Builtin().RenderRef(fileCtx)
}

func (lit *GoStringNumber) RenderZero(fileCtx GoFileContext) string {
	return lit.NumberExpr().RenderExpr(fileCtx)
}

func NumberEqual(from GoNumber, to GoNumber) bool {
	switch u := from.(type) {
	case *GoBoolNumber:
		switch v := to.(type) {
		case *GoBoolNumber:
			return u.Value_ == v.Value_
		}
	case *GoIntegerNumber:
		switch v := to.(type) {
		case *GoIntegerNumber:
			return u.Value_ == v.Value_
		}
	case *GoFloatNumber:
		switch v := to.(type) {
		case *GoFloatNumber:
			return u.Value_ == v.Value_
		}
	case *GoStringNumber:
		switch v := to.(type) {
		case *GoStringNumber:
			return u.Value_ == v.Value_
		}
	}
	return false
}
