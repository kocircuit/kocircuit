package model

import (
	"bytes"
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type GoSketchCtx struct {
	Seen map[GoType]bool `ko:"name=seen"`
}

func NewGoSketchCtx() *GoSketchCtx {
	return &GoSketchCtx{Seen: map[GoType]bool{}}
}

func (ctx *GoSketchCtx) Reserve(t GoType) bool {
	_, reserved := ctx.Seen[t]
	ctx.Seen[t] = true
	return !reserved
}

func (ctx *GoSketchCtx) Sketch(t ...GoType) interface{} {
	switch len(t) {
	case 0:
		return nil
	case 1:
		return t[0].Sketch(ctx)
	default:
		r := make([]interface{}, len(t))
		for i := range t {
			r[i] = t[i].Sketch(ctx)
		}
		return r
	}
}

type K map[string]interface{}

// GoType describes a go type recursively. The concrete types are:
// GoAlias,
// GoPtr, GoNeverNilPtr,
// GoMap, GoStruct,
// GoSlice, GoArray,
// GoBuiltin,
// GoBoolNumber, GoIntegerNumber, GoFloatNumber, GoStringNumber,
// GoFunc, GoInterface, GoChan
// GoUnknown, GoVariety, GoEmpty
type GoType interface {
	Doc() string
	String() string
	Sketch(*GoSketchCtx) interface{}
	// RenderRef returns the go expression referring to this type.
	RenderRef(GoFileContext) string
	// RenderDef returns the go definition of this type.
	RenderDef(GoFileContext) string
	// RenderZero returns the go expression for the zero value of this type.
	RenderZero(GoFileContext) string
	// TypeID is an identifier for the semantic equivalence of types.
	// TypeIDs can be used for semantic (i.e. comparison) purposes, not for syntactic (i.e rendering) purposes.
	TypeID() string
	// Tag returns key/value tags associated with the type.
	Tag() []*GoTag
}

// GoAlias captures a Go type alias.
type GoAlias struct {
	ID      string     `ko:"name=id"`
	Address *GoAddress `ko:"name=address"`
	For     GoType     `ko:"name=for"`
}

func NewAliasDeferType(address *GoAddress) *GoAlias {
	return &GoAlias{
		ID:      Mix("address", address.TypeID()),
		Address: address,
	}
}

func NewGoAlias(address *GoAddress, forType GoType) *GoAlias {
	return &GoAlias{
		ID:      Mix("address", address.TypeID()),
		Address: address,
		For:     forType,
	}
}

func (alias *GoAlias) TypeID() string { return alias.ID }

func (alias *GoAlias) StructureField() []*GoField {
	return alias.For.(GoStructure).StructureField()
}

func (alias *GoAlias) Doc() string { return alias.Address.Doc() }

func (alias *GoAlias) String() string { return Sprint(alias) }

func (alias *GoAlias) Sketch(ctx *GoSketchCtx) interface{} {
	if !ctx.Reserve(alias) {
		return K{"!ref": alias.Address.String()}
	} else {
		return K{
			"!alias": alias.Address.String(),
			"!for":   alias.For.Sketch(ctx),
		}
	}
}

func (alias *GoAlias) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: Alias TypeRef
func (alias *GoAlias) RenderDef(fileCtx GoFileContext) string {
	if alias2, ok := alias.For.(*GoAlias); ok {
		return fmt.Sprintf("%s %s", alias.Address.RenderExpr(fileCtx), alias2.RenderRef(fileCtx))
	}
	return fmt.Sprintf("%s %s", alias.Address.RenderExpr(fileCtx), alias.For.RenderDef(fileCtx))
}

// RenderRef returns a type reference of the form: pkgAlias.Alias
func (alias *GoAlias) RenderRef(fileCtx GoFileContext) string {
	return alias.Address.RenderExpr(fileCtx)
}

func (alias *GoAlias) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s(%s)", alias.RenderRef(fileCtx), alias.For.RenderZero(fileCtx))
}

// GoMap captures a Go map.
type GoMap struct {
	ID    string `ko:"name=id"`
	Key   GoType `ko:"name=key"`
	Value GoType `ko:"name=value"`
}

func NewGoMap(key, value GoType) *GoMap {
	return &GoMap{
		ID:    Mix("key", key.TypeID(), "value", value.TypeID()),
		Key:   key,
		Value: value,
	}
}

func (goMap *GoMap) TypeID() string { return goMap.ID }

func (goMap *GoMap) Doc() string { return "GoMap" }

func (goMap *GoMap) String() string { return Sprint(goMap) }

func (goMap *GoMap) Sketch(ctx *GoSketchCtx) interface{} {
	return K{"!goMapKey": goMap.Key.Sketch(ctx), "!goMapValue": goMap.Value.Sketch(ctx)}
}

func (goMap *GoMap) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: map[KeyTypeRef]ValueTypeRef
func (goMap *GoMap) RenderDef(fileCtx GoFileContext) string {
	return fmt.Sprintf("map[%s]%s", goMap.Key.RenderRef(fileCtx), goMap.Value.RenderRef(fileCtx))
}

// RenderRef returns a type reference of the form: map[KeyTypeRef]ValueTypeRef
func (goMap *GoMap) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("map[%s]%s", goMap.Key.RenderRef(fileCtx), goMap.Value.RenderRef(fileCtx))
}

// RenderZero returns a zero value of the form: nil
func (goMap *GoMap) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s{}", goMap.RenderRef(fileCtx))
}

// GoChan captures a Go map.
// TODO: Add channel direction.
type GoChan struct {
	ID   string `ko:"name=id"`
	Elem GoType `ko:"name=elem"`
}

func NewGoChan(elem GoType) *GoChan {
	return &GoChan{ID: Mix("chan", elem.TypeID()), Elem: elem}
}

func (goChan *GoChan) TypeID() string { return goChan.ID }

func (goChan *GoChan) Doc() string { return fmt.Sprintf("GoChan(%s)", goChan.Elem.Doc()) }

func (goChan *GoChan) String() string { return Sprint(goChan) }

func (goChan *GoChan) Sketch(ctx *GoSketchCtx) interface{} {
	return K{"!goChan": goChan.Elem.Sketch(ctx)}
}

func (goChan *GoChan) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: chan ElemTypeRef
func (goChan *GoChan) RenderDef(fileCtx GoFileContext) string {
	return fmt.Sprintf("chan %s", goChan.Elem.RenderRef(fileCtx))
}

// RenderRef returns a type reference of the form: chan ElemTypeRef
func (goChan *GoChan) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("chan %s", goChan.Elem.RenderRef(fileCtx))
}

// RenderZero returns a zero value of the form: nil
func (goChan *GoChan) RenderZero(_ GoFileContext) string { return "nil" }

// GoFunc captures a Go function.
type GoFunc struct {
	ID       string   `ko:"name=id"`
	Arg      []GoType `ko:"name=arg"`
	Returns  []GoType `ko:"name=returns"`
	Variadic bool     `ko:"name=variadic"`
}

func NewGoFunc(arg, returns []GoType, variadic bool) *GoFunc {
	mix := []string{"arg"}
	for _, a := range arg {
		mix = append(mix, a.TypeID())
	}
	mix = append(mix, "returns")
	for _, r := range returns {
		mix = append(mix, r.TypeID())
	}
	return &GoFunc{
		ID:       Mix(mix...),
		Arg:      arg,
		Returns:  returns,
		Variadic: variadic,
	}
}

func (goFunc *GoFunc) TypeID() string { return goFunc.ID }

func (goFunc *GoFunc) Doc() string { return "GoFunc" }

func (goFunc *GoFunc) String() string { return Sprint(goFunc) }

func (goFunc *GoFunc) Sketch(ctx *GoSketchCtx) interface{} {
	return K{
		"!goFuncArg":      ctx.Sketch(goFunc.Arg...),
		"!goFuncVariadic": goFunc.Variadic,
		"!goFuncReturns":  ctx.Sketch(goFunc.Returns...),
	}
}

func (goFunc *GoFunc) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: func (A1, ...) (R1, ...)
func (goFunc *GoFunc) RenderDef(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "func(")
	for i, a := range goFunc.Arg {
		if i != 0 {
			fmt.Fprint(&w, ", ")
		}
		if goFunc.Variadic && i+1 == len(goFunc.Arg) {
			fmt.Fprint(&w, "...")
		}
		fmt.Fprint(&w, a.RenderRef(fileCtx))
	}
	fmt.Fprintf(&w, ") (")
	for i, r := range goFunc.Returns {
		if i != 0 {
			fmt.Fprint(&w, ", ")
		}
		fmt.Fprint(&w, r.RenderRef(fileCtx))
	}
	fmt.Fprintf(&w, ")")
	return w.String()
}

// RenderRef returns a type reference of the form: func (A1, ...) (R1, ...)
func (goFunc *GoFunc) RenderRef(fileCtx GoFileContext) string {
	return goFunc.RenderDef(fileCtx)
}

// RenderZero returns a zero value of the form: nil
func (goFunc *GoFunc) RenderZero(_ GoFileContext) string { return "nil" }

// GoInterface captures a Go interface.
type GoInterface struct {
	ID   string       `ko:"name=id"`
	Type reflect.Type `ko:"name=type"`
}

func NewGoInterface(typ reflect.Type) *GoInterface {
	return &GoInterface{ID: Mix("interface", typ.PkgPath(), typ.Name()), Type: typ}
}

func (goInterface *GoInterface) TypeID() string { return goInterface.ID }

func (goInterface *GoInterface) IsRestricted() bool {
	var v interface{}
	return !goInterface.Type.AssignableTo(reflect.TypeOf(&v).Elem())
}

func (goInterface *GoInterface) Doc() string { return "GoInterface" }

func (goInterface *GoInterface) String() string { return Sprint(goInterface) }

func (goInterface *GoInterface) Sketch(ctx *GoSketchCtx) interface{} {
	return K{"!goInterface": goInterface.Type.String()}
}

func (goInterface *GoInterface) Tag() []*GoTag { return nil }

func (goInterface *GoInterface) address() *GoAddress {
	return &GoAddress{
		GroupPath: GoGroupPath{
			Group: GoHereditaryPkgGroup,
			Path:  goInterface.Type.PkgPath(),
		},
		Name: goInterface.Type.Name(),
	}
}

// RenderDef returns a type definition of the form: interface{...}
func (goInterface *GoInterface) RenderDef(_ GoFileContext) string {
	panic("o")
}

type unrestrictedInterfaceStruct struct {
	I interface{}
}

var UnrestrictedInterface = reflect.TypeOf(unrestrictedInterfaceStruct{}).Field(0).Type

// RenderRef returns a type reference of the form: interface{...}
func (goInterface *GoInterface) RenderRef(fileCtx GoFileContext) string {
	if goInterface.Type == UnrestrictedInterface {
		return "interface{}"
	} else {
		return goInterface.address().RenderExpr(fileCtx)
	}
}

// RenderZero returns a zero value of the form: nil
func (goInterface *GoInterface) RenderZero(_ GoFileContext) string { return "nil" }
