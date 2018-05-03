package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGateAt("", "WeaveFigure", new(WeaveFigure))
	RegisterEvalGateAt("", "WeaveFunctional", new(WeaveFunctional))
	RegisterEvalGateAt("", "WeaveFunc", new(WeaveFunc))
	RegisterEvalGateAt("", "WeaveReserve", new(WeaveReserve))
}

type WeaveStepCtx struct {
	Origin *Span  `ko:"name=origin"` // evaluation span (not weave span)
	Pkg    string `ko:"name=pkg"`
	Func   string `ko:"name=func"`
	Step   string `ko:"name=step"`
	Logic  string `ko:"name=logic"`
	Source string `ko:"name=source"`
	Ctx    Symbol `ko:"name=ctx"` // user ctx object
}

func (ctx *WeaveStepCtx) DelegateSpan() *Span {
	return RefineOutline(ctx.Origin, fmt.Sprintf("%s @ %s", ctx.Logic, ctx.Source))
}

func (ctx *WeaveStepCtx) Deconstruct(span *Span) Symbol {
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

type WeaveField struct {
	Name    string `ko:"name=name"`
	Monadic bool   `ko:"name=monadic"`
	Objects Symbol `ko:"name=objects"`
}

func (field *WeaveField) Deconstruct(span *Span) Symbol {
	return MakeStructSymbol(
		FieldSymbols{
			{Name: "name", Value: BasicSymbol{field.Name}},
			{Name: "monadic", Value: BasicSymbol{field.Monadic}},
			{Name: "objects", Value: field.Objects},
		},
	)
}

type WeaveFields []*WeaveField

func (bf WeaveFields) Deconstruct(span *Span) (Symbol, error) {
	elem := make(Symbols, len(bf))
	for i := range bf {
		elem[i] = bf[i].Deconstruct(span)
	}
	return MakeSeriesSymbol(span, elem)
}

type WeaveFigure struct {
	Int64      *int64           `ko:"name=int64"`
	String     *string          `ko:"name=string"`
	Bool       *bool            `ko:"name=bool"`
	Float64    *float64         `ko:"name=float64"`
	Functional *WeaveFunctional `ko:"name=functional"`
}

func (fig *WeaveFigure) Play(*runtime.Context) *WeaveFigure {
	return fig
}

type WeaveFunctional struct {
	Reserve *WeaveReserve `ko:"name=reserve"`
	Func    *WeaveFunc    `ko:"name=func"`
}

func (w *WeaveFunctional) Play(*runtime.Context) *WeaveFunctional {
	return w
}

type WeaveReserve struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

func (w *WeaveReserve) Play(*runtime.Context) *WeaveReserve {
	return w
}

type WeaveFunc struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

func (w *WeaveFunc) Play(*runtime.Context) *WeaveFunc {
	return w
}

func (fig *WeaveFigure) Deconstruct(span *Span) Symbol {
	return DeconstructInterface(span, fig)
}

type WeaveResidue struct {
	Step    string `ko:"name=step"`
	Logic   string `ko:"name=logic"`
	Source  string `ko:"name=source"`
	Returns Symbol `ko:"name=returns"`
	Effect  Symbol `ko:"name=effect"`
}

func (residue *WeaveResidue) Deconstruct(span *Span) Symbol {
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

type WeaveResidues []*WeaveResidue

func (br WeaveResidues) Deconstruct(span *Span) (Symbol, error) {
	elem := make(Symbols, len(br))
	for i := range br {
		elem[i] = br[i].Deconstruct(span)
	}
	return MakeSeriesSymbol(span, elem)
}

type WeaveSummary struct {
	Origin  *Span  `ko:"name=origin"` // evaluation span (not weave span)
	Pkg     string `ko:"name=pkg"`
	Func    string `ko:"name=func"`
	Source  string `ko:"name=source"`
	Ctx     Symbol `ko:"name=ctx"` // user ctx object
	Arg     Symbol `ko:"name=arg"`
	Returns Symbol `ko:"name=returns"`
}

func (summary *WeaveSummary) CombineSpan() *Span {
	return RefineOutline(summary.Origin, fmt.Sprintf("COMBINE @ %s", summary.Source))
}

func (summary *WeaveSummary) Deconstruct(span *Span) Symbol {
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

type WeaveStepResult struct {
	Returns Symbol `ko:"name=returns"`
	Effect  Symbol `ko:"name=effect"`
}
