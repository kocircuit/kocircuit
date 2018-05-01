package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type BootStepCtx struct {
	Pkg    string `ko:"name=pkg"`
	Func   string `ko:"name=func"`
	Step   string `ko:"name=step"`
	Logic  string `ko:"name=logic"`
	Source string `ko:"name=source"`
	Ctx    Symbol `ko:"name=ctx"` // user ctx object
}

func (ctx *BootStepCtx) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "pkg", Value: BasicSymbol{ctx.Pkg}},
			{Name: "func", Value: BasicSymbol{ctx.Func}},
			{Name: "step", Value: BasicSymbol{ctx.Step}},
			{Name: "logic", Value: BasicSymbol{ctx.Logic}},
			{Name: "source", Value: BasicSymbol{ctx.Source}},
			{Name: "ctx", Value: ctx.Ctx},
		},
	)
}

type BootField struct {
	Name   string `ko:"name=name"`
	Object Symbol `ko:"name=object"`
}

func (field *BootField) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "name", Value: BasicSymbol{field.Name}},
			{Name: "object", Value: field.Object},
		},
	)
}

type BootFields []*BootField

func (bf BootFields) Deconstruct(span *Span) (Symbol, error) {
	elem := make(Symbols, len(bf))
	for i := range bf {
		elem[i] = bf[i].Deconstruct(span)
	}
	return MakeSeriesSymbol(span, elem)
}

type BootFigure struct {
	Int64      *int64          `ko:"name=int64"`
	String     *string         `ko:"name=string"`
	Bool       *bool           `ko:"name=bool"`
	Float64    *float64        `ko:"name=float64"`
	Functional *BootFunctional `ko:"name=functional"`
}

type BootFunctional struct {
	Macro *string   `ko:"name=macro"`
	Func  *BootFunc `ko:"name=func"`
}

type BootFunc struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

func (fig *BootFigure) Deconstruct(span *Span) Symbol {
	return DeconstructInterface(span, fig)
}

type BootResidue struct {
	Returned Symbol `ko:"name=returned"`
	Effect   Symbol `ko:"name=effect"`
}

func (residue *BootResidue) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "returned", Value: residue.Returned},
			{Name: "effect", Value: residue.Effect},
		},
	)
}

type BootResidues []*BootResidue

func (br BootResidues) Deconstruct(span *Span) (Symbol, error) {
	elem := make(Symbols, len(br))
	for i := range br {
		elem[i] = br[i].Deconstruct(span)
	}
	return MakeSeriesSymbol(span, elem)
}

type BootSummary struct {
	Pkg      string `ko:"name=pkg"`
	Func     string `ko:"name=func"`
	Source   string `ko:"name=source"`
	Ctx      Symbol `ko:"name=ctx"` // user ctx object
	Arg      Symbol `ko:"name=arg"`
	Returned Symbol `ko:"name=returned"`
	Panicked Symbol `ko:"name=panicked"`
}

func (summary *BootSummary) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "pkg", Value: BasicSymbol{summary.Pkg}},
			{Name: "func", Value: BasicSymbol{summary.Func}},
			{Name: "source", Value: BasicSymbol{summary.Source}},
			{Name: "ctx", Value: summary.Ctx},
			{Name: "arg", Value: summary.Arg},
			{Name: "returned", Value: summary.Returned},
			{Name: "panicked", Value: summary.Panicked},
		},
	)
}
