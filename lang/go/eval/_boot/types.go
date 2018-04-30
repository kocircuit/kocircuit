package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type BootStepCtx struct {
	Pkg    string `ko:"name=pkg"`
	Func   string `ko:"name=func"`
	Step   string `ko:"name=step"`
	Source string `ko:"name=source"`
	Ctx    Symbol `ko:"name=ctx"` // user ctx object
}

type BootField struct {
	Name   string `ko:"name=name"`
	Object Symbol `ko:"name=object"`
}

type BootFigure struct {
	Int64   *int64         `ko:"name=int64"`
	String  *string        `ko:"name=string"`
	Bool    *bool          `ko:"name=bool"`
	Float64 *float64       `ko:"name=float64"`
	Func    *VarietySymbol `ko:"name=func"`
}

type BootResidue struct {
	Returned Symbol `ko:"name=returned"`
	Effect   Symbol `ko:"name=effect"`
}

type BootSummary struct {
	Pkg    string `ko:"name=pkg"`
	Func   string `ko:"name=func"`
	Source string `ko:"name=source"`
	Ctx    Symbol `ko:"name=ctx"` // user ctx object
	//
	Arg      Symbol `ko:"name=arg"`
	Returned Symbol `ko:"name=returned"`
	Panicked Symbol `ko:"name=panicked"`
}
