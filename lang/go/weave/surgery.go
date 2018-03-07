package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Bool", new(GoBoolMacro))
	RegisterGoMacro("String", new(GoStringMacro))
	RegisterGoMacro("Byte", new(GoInt8Macro))
	RegisterGoMacro("Int8", new(GoInt8Macro))
	RegisterGoMacro("Int16", new(GoInt16Macro))
	RegisterGoMacro("Int32", new(GoInt32Macro))
	RegisterGoMacro("Int64", new(GoInt64Macro))
	RegisterGoMacro("Uint8", new(GoUint8Macro))
	RegisterGoMacro("Uint16", new(GoUint16Macro))
	RegisterGoMacro("Uint32", new(GoUint32Macro))
	RegisterGoMacro("Uint64", new(GoUint64Macro))
	RegisterGoMacro("Float32", new(GoFloat32Macro))
	RegisterGoMacro("Float64", new(GoFloat64Macro))
}

type GoBoolMacro struct{}

func (m GoBoolMacro) MacroID() string            { return m.Help() }
func (m GoBoolMacro) Label() string              { return goBoolMacro.Label() }
func (m GoBoolMacro) MacroSheathString() *string { return PtrString("Bool") }
func (m GoBoolMacro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoBoolMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goBoolMacro.Invoke(span, arg)
}

type GoStringMacro struct{}

func (m GoStringMacro) MacroID() string            { return m.Help() }
func (m GoStringMacro) Label() string              { return goStringMacro.Label() }
func (m GoStringMacro) MacroSheathString() *string { return PtrString("String") }
func (m GoStringMacro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoStringMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goStringMacro.Invoke(span, arg)
}

type GoInt8Macro struct{}

func (m GoInt8Macro) MacroID() string            { return m.Help() }
func (m GoInt8Macro) Label() string              { return goInt8Macro.Label() }
func (m GoInt8Macro) MacroSheathString() *string { return PtrString("Int8") }
func (m GoInt8Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoInt8Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goInt8Macro.Invoke(span, arg)
}

type GoInt16Macro struct{}

func (m GoInt16Macro) MacroID() string            { return m.Help() }
func (m GoInt16Macro) Label() string              { return goInt16Macro.Label() }
func (m GoInt16Macro) MacroSheathString() *string { return PtrString("Int16") }
func (m GoInt16Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoInt16Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goInt16Macro.Invoke(span, arg)
}

type GoInt32Macro struct{}

func (m GoInt32Macro) MacroID() string            { return m.Help() }
func (m GoInt32Macro) Label() string              { return goInt32Macro.Label() }
func (m GoInt32Macro) MacroSheathString() *string { return PtrString("Int32") }
func (m GoInt32Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoInt32Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goInt32Macro.Invoke(span, arg)
}

type GoInt64Macro struct{}

func (m GoInt64Macro) MacroID() string            { return m.Help() }
func (m GoInt64Macro) Label() string              { return goInt64Macro.Label() }
func (m GoInt64Macro) MacroSheathString() *string { return PtrString("Int64") }
func (m GoInt64Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoInt64Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goInt64Macro.Invoke(span, arg)
}

type GoUint8Macro struct{}

func (m GoUint8Macro) MacroID() string            { return m.Help() }
func (m GoUint8Macro) Label() string              { return goUint8Macro.Label() }
func (m GoUint8Macro) MacroSheathString() *string { return PtrString("Uint8") }
func (m GoUint8Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoUint8Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goUint8Macro.Invoke(span, arg)
}

type GoUint16Macro struct{}

func (m GoUint16Macro) MacroID() string            { return m.Help() }
func (m GoUint16Macro) Label() string              { return goUint16Macro.Label() }
func (m GoUint16Macro) MacroSheathString() *string { return PtrString("Uint16") }
func (m GoUint16Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoUint16Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goUint16Macro.Invoke(span, arg)
}

type GoUint32Macro struct{}

func (m GoUint32Macro) MacroID() string            { return m.Help() }
func (m GoUint32Macro) Label() string              { return goUint32Macro.Label() }
func (m GoUint32Macro) MacroSheathString() *string { return PtrString("Uint32") }
func (m GoUint32Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoUint32Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goUint32Macro.Invoke(span, arg)
}

type GoUint64Macro struct{}

func (m GoUint64Macro) MacroID() string            { return m.Help() }
func (m GoUint64Macro) Label() string              { return goUint64Macro.Label() }
func (m GoUint64Macro) MacroSheathString() *string { return PtrString("Uint64") }
func (m GoUint64Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoUint64Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goUint64Macro.Invoke(span, arg)
}

type GoFloat32Macro struct{}

func (m GoFloat32Macro) MacroID() string            { return m.Help() }
func (m GoFloat32Macro) Label() string              { return goFloat32Macro.Label() }
func (m GoFloat32Macro) MacroSheathString() *string { return PtrString("Float32") }
func (m GoFloat32Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoFloat32Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goFloat32Macro.Invoke(span, arg)
}

type GoFloat64Macro struct{}

func (m GoFloat64Macro) MacroID() string            { return m.Help() }
func (m GoFloat64Macro) Label() string              { return goFloat64Macro.Label() }
func (m GoFloat64Macro) MacroSheathString() *string { return PtrString("Float64") }
func (m GoFloat64Macro) Help() string               { return GoInterfaceTypeAddress(m).String() }
func (GoFloat64Macro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return goFloat64Macro.Invoke(span, arg)
}

var (
	goBoolMacro    = &goTypeMacro{Label_: "Bool", Type: GoBool}
	goStringMacro  = &goTypeMacro{Label_: "String", Type: GoString}
	goInt8Macro    = &goTypeMacro{Label_: "Int8", Type: GoInt8}
	goInt16Macro   = &goTypeMacro{Label_: "Int16", Type: GoInt16}
	goInt32Macro   = &goTypeMacro{Label_: "Int32", Type: GoInt32}
	goInt64Macro   = &goTypeMacro{Label_: "Int64", Type: GoInt64}
	goUint8Macro   = &goTypeMacro{Label_: "Uint8", Type: GoUint8}
	goUint16Macro  = &goTypeMacro{Label_: "Uint16", Type: GoUint16}
	goUint32Macro  = &goTypeMacro{Label_: "Uint32", Type: GoUint32}
	goUint64Macro  = &goTypeMacro{Label_: "Uint64", Type: GoUint64}
	goFloat32Macro = &goTypeMacro{Label_: "Float32", Type: GoFloat32}
	goFloat64Macro = &goTypeMacro{Label_: "Float64", Type: GoFloat64}
)

type goTypeMacro struct {
	Label_ string `ko:"name=label"`
	Type   GoType `ko:"name=type"`
}

func (m *goTypeMacro) Label() string { return m.Label_ }

func (m *goTypeMacro) Help() string { return m.Label_ }

func (m *goTypeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if f, err := solveTypeMacro(span, arg, m.Type); err != nil {
		return nil, nil, err
	} else {
		return m.Type,
			SlotFormMacroEffect(arg.(GoStructure), f),
			nil
	}
}

func solveTypeMacro(span *Span, arg Arg, returns GoType) (f *convertSlotForm, err error) {
	if monadic := StructureMonadicField(arg.(GoStructure)); monadic == nil {
		return nil, span.Errorf(nil, "type macro expects a monadic argument, got %s", Sprint(arg))
	} else {
		f = &convertSlotForm{Returns: returns}
		if f.Extractor, f.Number, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg.(GoStructure)); err != nil {
			return nil, span.Errorf(err, "type macro expects a monadic argument")
		}
		assn := NewAssignCtx(span)
		if f.Converter, err = assn.Assign(f.Number, f.Returns); err != nil {
			return nil, err
		}
		return
	}
}
