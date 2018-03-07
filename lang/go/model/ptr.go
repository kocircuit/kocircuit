package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// GoModifier captures GoPtr, GoNeverNilPtr. More generally,
// GoModifier is a type that modifies (the semantics of) a singe other type.
type GoModifier interface {
	Modified() GoType
}

// GoPtr captures a Go pointer.
type GoPtr struct {
	ID   string `ko:"name=id"`
	Elem GoType `ko:"name=elem"`
}

func NewGoPtr(elem GoType) *GoPtr {
	return &GoPtr{ID: Mix("ptr", elem.TypeID()), Elem: elem}
}

func (ptr *GoPtr) TypeID() string { return ptr.ID }

func (ptr *GoPtr) Modified() GoType { return ptr.Elem }

func (ptr *GoPtr) Doc() string { return ptr.Elem.Doc() }

func (ptr *GoPtr) String() string { return Sprint(ptr) }

func (ptr *GoPtr) Sketch(ctx *GoSketchCtx) interface{} {
	e := ptr.Elem.Sketch(ctx)
	return &e
}

func (ptr *GoPtr) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: *TypeDef
func (ptr *GoPtr) RenderDef(fileCtx GoFileContext) string {
	switch a := ptr.Elem.(type) {
	case *GoAlias:
		return a.RenderDef(fileCtx)
	default:
		return fmt.Sprintf("*%s", ptr.Elem.RenderRef(fileCtx))
	}
}

// RenderRef returns a type reference of the form: *TypeRef
func (ptr *GoPtr) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("*%s", ptr.Elem.RenderRef(fileCtx))
}

// RenderZero returns a zero value of the form: nil
func (ptr *GoPtr) RenderZero(_ GoFileContext) string { return "nil" }

// GoNeverNilPtr captures a never-nil Go pointer to a circuit valve.
type GoNeverNilPtr struct {
	ID   string `ko:"name=id"`
	Elem GoType `ko:"name=elem"` // argument structure
}

func NewGoNeverNilPtr(elem GoType) *GoNeverNilPtr {
	return &GoNeverNilPtr{ID: Mix("nnptr", elem.TypeID()), Elem: elem}
}

func (neverNilPtr *GoNeverNilPtr) TypeID() string { return neverNilPtr.ID }

func (neverNilPtr *GoNeverNilPtr) Modified() GoType { return neverNilPtr.Elem }

func (neverNilPtr *GoNeverNilPtr) StructureField() []*GoField {
	return neverNilPtr.Elem.(GoStructure).StructureField()
}

func (neverNilPtr *GoNeverNilPtr) Doc() string { return neverNilPtr.Elem.Doc() }

func (neverNilPtr *GoNeverNilPtr) String() string { return Sprint(neverNilPtr) }

func (neverNilPtr *GoNeverNilPtr) Sketch(ctx *GoSketchCtx) interface{} {
	return neverNilPtr.Elem.Sketch(ctx)
}

func (neverNilPtr *GoNeverNilPtr) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: *TypeDef
func (neverNilPtr *GoNeverNilPtr) RenderDef(fileCtx GoFileContext) string {
	switch a := neverNilPtr.Elem.(type) {
	case *GoAlias:
		return a.RenderDef(fileCtx)
	default:
		return fmt.Sprintf("*%s", neverNilPtr.Elem.RenderRef(fileCtx))
	}
}

// RenderRef returns a type reference of the form: *TypeRef
func (neverNilPtr *GoNeverNilPtr) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("*%s", neverNilPtr.Elem.RenderRef(fileCtx))
}

func (neverNilPtr *GoNeverNilPtr) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("&%s{}", neverNilPtr.Elem.RenderRef(fileCtx))
}
