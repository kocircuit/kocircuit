package boot

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type BootStepCtx struct {
	Origin *Span  `ko:"name=origin"` // evaluation span (not boot span)
	Pkg    string `ko:"name=pkg"`
	Func   string `ko:"name=func"`
	Step   string `ko:"name=step"`
	Logic  string `ko:"name=logic"`
	Source string `ko:"name=source"`
	Ctx    Symbol `ko:"name=ctx"` // user ctx object
}

func (ctx *BootStepCtx) DelegateSpan() *Span {
	return RefineOutline(ctx.Origin, fmt.Sprintf("%s @ %s", ctx.Logic, ctx.Source))
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
	Name    string `ko:"name=name"`
	Monadic bool   `ko:"name=monadic"`
	Objects Symbol `ko:"name=objects"`
}

func (field *BootField) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "name", Value: BasicSymbol{field.Name}},
			{Name: "monadic", Value: BasicSymbol{field.Monadic}},
			{Name: "objects", Value: field.Objects},
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
	Reserve *BootReserve `ko:"name=reserve"`
	Func    *BootFunc    `ko:"name=func"`
}

type BootReserve struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

type BootFunc struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

func (fig *BootFigure) Deconstruct(span *Span) Symbol {
	return DeconstructInterface(span, fig)
}

type BootResidue struct {
	Step    string `ko:"name=step"`
	Logic   string `ko:"name=logic"`
	Source  string `ko:"name=source"`
	Returns Symbol `ko:"name=returns"`
	Effect  Symbol `ko:"name=effect"`
}

func (residue *BootResidue) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "step", Value: BasicSymbol{residue.Step}},
			{Name: "logic", Value: BasicSymbol{residue.Logic}},
			{Name: "source", Value: BasicSymbol{residue.Source}},
			{Name: "returns", Value: residue.Returns},
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
	Origin  *Span  `ko:"name=origin"` // evaluation span (not boot span)
	Pkg     string `ko:"name=pkg"`
	Func    string `ko:"name=func"`
	Source  string `ko:"name=source"`
	Ctx     Symbol `ko:"name=ctx"` // user ctx object
	Arg     Symbol `ko:"name=arg"`
	Returns Symbol `ko:"name=returns"`
}

func (summary *BootSummary) CombineSpan() *Span {
	return RefineOutline(summary.Origin, fmt.Sprintf("COMBINE @ %s", summary.Source))
}

func (summary *BootSummary) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "pkg", Value: BasicSymbol{summary.Pkg}},
			{Name: "func", Value: BasicSymbol{summary.Func}},
			{Name: "source", Value: BasicSymbol{summary.Source}},
			{Name: "ctx", Value: summary.Ctx},
			{Name: "arg", Value: summary.Arg},
			{Name: "returns", Value: summary.Returns},
		},
	)
}

type BootStepResult struct {
	Returns Symbol `ko:"name=returns"`
	Effect  Symbol `ko:"name=effect"`
}
