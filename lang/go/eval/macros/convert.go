package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("String", new(EvalStringMacro))
	RegisterEvalMacro("Bool", new(EvalBoolMacro))
	RegisterEvalMacro("Int8", new(EvalInt8Macro))
	RegisterEvalMacro("Int16", new(EvalInt16Macro))
	RegisterEvalMacro("Int32", new(EvalInt32Macro))
	RegisterEvalMacro("Int64", new(EvalInt64Macro))
	RegisterEvalMacro("Uint8", new(EvalUint8Macro))
	RegisterEvalMacro("Uint16", new(EvalUint16Macro))
	RegisterEvalMacro("Uint32", new(EvalUint32Macro))
	RegisterEvalMacro("Uint64", new(EvalUint64Macro))
	RegisterEvalMacro("Float32", new(EvalFloat32Macro))
	RegisterEvalMacro("Float64", new(EvalFloat64Macro))
}

type EvalStringMacro struct{}

func (m EvalStringMacro) MacroID() string            { return m.Help() }
func (m EvalStringMacro) Label() string              { return "string" }
func (m EvalStringMacro) MacroSheathString() *string { return PtrString("String") }
func (m EvalStringMacro) Help() string               { return "String" }
func (m EvalStringMacro) Doc() string {
	return `String converts its argument to a string.`
}

func (EvalStringMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicString)
}

type EvalBoolMacro struct{}

func (m EvalBoolMacro) MacroID() string            { return m.Help() }
func (m EvalBoolMacro) Label() string              { return "bool" }
func (m EvalBoolMacro) MacroSheathString() *string { return PtrString("Bool") }
func (m EvalBoolMacro) Help() string               { return "Bool" }
func (m EvalBoolMacro) Doc() string {
	return `Bool converts its argument to a boolean.`
}

func (EvalBoolMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicBool)
}

// signed integers

type EvalInt8Macro struct{}

func (m EvalInt8Macro) MacroID() string            { return m.Help() }
func (m EvalInt8Macro) Label() string              { return "int8" }
func (m EvalInt8Macro) MacroSheathString() *string { return PtrString("Int8") }
func (m EvalInt8Macro) Help() string               { return "Int8" }
func (m EvalInt8Macro) Doc() string {
	return `Int8 converts its argument an 8-bit signed integer.`
}

func (EvalInt8Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicInt8)
}

type EvalInt16Macro struct{}

func (m EvalInt16Macro) MacroID() string            { return m.Help() }
func (m EvalInt16Macro) Label() string              { return "int16" }
func (m EvalInt16Macro) MacroSheathString() *string { return PtrString("Int16") }
func (m EvalInt16Macro) Help() string               { return "Int16" }
func (m EvalInt16Macro) Doc() string {
	return `Int16 converts its argument a 16-bit signed integer.`
}

func (EvalInt16Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicInt16)
}

type EvalInt32Macro struct{}

func (m EvalInt32Macro) MacroID() string            { return m.Help() }
func (m EvalInt32Macro) Label() string              { return "int32" }
func (m EvalInt32Macro) MacroSheathString() *string { return PtrString("Int32") }
func (m EvalInt32Macro) Help() string               { return "Int32" }
func (m EvalInt32Macro) Doc() string {
	return `Int32 converts its argument a 32-bit signed integer.`
}

func (EvalInt32Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicInt32)
}

type EvalInt64Macro struct{}

func (m EvalInt64Macro) MacroID() string            { return m.Help() }
func (m EvalInt64Macro) Label() string              { return "int64" }
func (m EvalInt64Macro) MacroSheathString() *string { return PtrString("Int64") }
func (m EvalInt64Macro) Help() string               { return "Int64" }
func (m EvalInt64Macro) Doc() string {
	return `Int64 converts its argument a 64-bit signed integer.`
}

func (EvalInt64Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicInt64)
}

// unsigned integers

type EvalUint8Macro struct{}

func (m EvalUint8Macro) MacroID() string            { return m.Help() }
func (m EvalUint8Macro) Label() string              { return "uint8" }
func (m EvalUint8Macro) MacroSheathString() *string { return PtrString("Uint8") }
func (m EvalUint8Macro) Help() string               { return "Uint8" }
func (m EvalUint8Macro) Doc() string {
	return `Uint8 converts its argument to an 8-bit unsigned integer.`
}

func (EvalUint8Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicUint8)
}

type EvalUint16Macro struct{}

func (m EvalUint16Macro) MacroID() string            { return m.Help() }
func (m EvalUint16Macro) Label() string              { return "uint16" }
func (m EvalUint16Macro) MacroSheathString() *string { return PtrString("Uint16") }
func (m EvalUint16Macro) Help() string               { return "Uint16" }
func (m EvalUint16Macro) Doc() string {
	return `Uint16 converts its argument to a 16-bit unsigned integer.`
}

func (EvalUint16Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicUint16)
}

type EvalUint32Macro struct{}

func (m EvalUint32Macro) MacroID() string            { return m.Help() }
func (m EvalUint32Macro) Label() string              { return "uint32" }
func (m EvalUint32Macro) MacroSheathString() *string { return PtrString("Uint32") }
func (m EvalUint32Macro) Help() string               { return "Uint32" }
func (m EvalUint32Macro) Doc() string {
	return `Uint32 converts its argument to a 32-bit unsigned integer.`
}

func (EvalUint32Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicUint32)
}

type EvalUint64Macro struct{}

func (m EvalUint64Macro) MacroID() string            { return m.Help() }
func (m EvalUint64Macro) Label() string              { return "uint64" }
func (m EvalUint64Macro) MacroSheathString() *string { return PtrString("Uint64") }
func (m EvalUint64Macro) Help() string               { return "Uint64" }
func (m EvalUint64Macro) Doc() string {
	return `Uint64 converts its argument to a 64-bit unsigned integer.`
}

func (EvalUint64Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicUint64)
}

// floating

type EvalFloat32Macro struct{}

func (m EvalFloat32Macro) MacroID() string            { return m.Help() }
func (m EvalFloat32Macro) Label() string              { return "float32" }
func (m EvalFloat32Macro) MacroSheathString() *string { return PtrString("Float32") }
func (m EvalFloat32Macro) Help() string               { return "Float32" }
func (m EvalFloat32Macro) Doc() string {
	return `Float32 converts its argument to a 32-bit floating-point number.`
}

func (EvalFloat32Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicFloat32)
}

type EvalFloat64Macro struct{}

func (m EvalFloat64Macro) MacroID() string            { return m.Help() }
func (m EvalFloat64Macro) Label() string              { return "float64" }
func (m EvalFloat64Macro) MacroSheathString() *string { return PtrString("Float64") }
func (m EvalFloat64Macro) Help() string               { return "Float64" }
func (m EvalFloat64Macro) Doc() string {
	return `Float64 converts its argument to a 64-bit floating-point number.`
}

func (EvalFloat64Macro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return convert(span, arg, BasicFloat64)
}

// shared

func convert(span *Span, arg Arg, to BasicType) (returns Return, effect Effect, err error) {
	if x := arg.(*StructSymbol).LinkField("value", true); x != nil {
		if basic, ok := x.(BasicSymbol); ok {
			if converted, err := basic.ConvertTo(span, to); err != nil {
				return nil, nil, err
			} else {
				return converted, nil, nil
			}
		} else {
			return nil, nil, span.Errorf(nil, "non-basic type %v is not convertible to %v", x, to)
		}
	} else {
		return nil, nil, span.Errorf(nil, "%v needs an argument", to)
	}
}
