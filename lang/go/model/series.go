package model

import (
	"fmt"
	"strconv"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type GoSequence interface {
	GoType
	Sequence
}

type Sequence interface {
	SequenceElem() GoType
}

// GoSlice captures a Go slice.
type GoSlice struct {
	ID   string `ko:"name=id"`
	Elem GoType `ko:"name=elem"`
}

func NewGoSlice(elem GoType) *GoSlice {
	return &GoSlice{ID: Mix("slice", elem.TypeID()), Elem: elem}
}

func (goSlice *GoSlice) TypeID() string { return goSlice.ID }

func (goSlice *GoSlice) SequenceElem() GoType { return goSlice.Elem }

func (goSlice *GoSlice) Doc() string { return goSlice.Elem.Doc() }

func (goSlice *GoSlice) String() string { return Sprint(goSlice) }

func (goSlice *GoSlice) Sketch(ctx *GoSketchCtx) interface{} {
	return []interface{}{goSlice.Elem.Sketch(ctx)}
}

func (goSlice *GoSlice) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: []ElemTypeRef
func (goSlice *GoSlice) RenderDef(fileCtx GoFileContext) string {
	return fmt.Sprintf("[]%s", goSlice.Elem.RenderRef(fileCtx))
}

// RenderRef returns a type reference of the form: []ElemTypeRef
func (goSlice *GoSlice) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("[]%s", goSlice.Elem.RenderRef(fileCtx))
}

// RenderZero returns a zero value of the form: nil
func (goSlice *GoSlice) RenderZero(_ GoFileContext) string { return "nil" }

// GoArray captures a Go array.
type GoArray struct {
	ID   string `ko:"name=id"`
	Len  int    `ko:"name=len"`
	Elem GoType `ko:"name=elem"`
}

func NewGoArray(l int, elem GoType) *GoArray {
	return &GoArray{
		ID:   Mix("array", strconv.Itoa(l), elem.TypeID()),
		Len:  l,
		Elem: elem,
	}
}

func (goArray *GoArray) TypeID() string { return goArray.ID }

func (goArray *GoArray) SequenceElem() GoType { return goArray.Elem }

func (goArray *GoArray) Doc() string { return goArray.Elem.Doc() }

func (goArray *GoArray) String() string { return Sprint(goArray) }

func (goArray *GoArray) Sketch(ctx *GoSketchCtx) interface{} {
	return []interface{}{goArray.Elem.Sketch(ctx)}
	// return K{fmt.Sprintf("array[%d]", goArray.Len): goArray.Elem.Sketch(ctx)}
}

func (goArray *GoArray) Tag() []*GoTag { return nil }

// RenderDef returns a type definition of the form: [Len]ElemTypeRef
func (goArray *GoArray) RenderDef(fileCtx GoFileContext) string {
	return fmt.Sprintf("[%d]%s", goArray.Len, goArray.Elem.RenderRef(fileCtx))
}

// RenderRef returns a type reference of the form: [Len]ElemTypeRef
func (goArray *GoArray) RenderRef(fileCtx GoFileContext) string {
	return fmt.Sprintf("[%d]%s", goArray.Len, goArray.Elem.RenderRef(fileCtx))
}

// RenderZero returns a zero value of the form: [Len]ElemTypeRef{}
func (goArray *GoArray) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("[%d]%s{}", goArray.Len, goArray.Elem.RenderRef(fileCtx))
}
