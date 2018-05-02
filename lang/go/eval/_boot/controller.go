package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Boot struct {
	Idiom  Repo
	Repo   Repo
	Booter *Booter
	Func   *Func
	Ctx    Symbol
	Arg    Symbol
}

func (b *Boot) Play(origin *Span) (returned, effect Symbol, err error) {
	bootController := &BootController{
		Origin: origin,
		Booter: booter,
		Func:   b.Func,
		Ctx:    b.Ctx,
		Arg:    b.Arg,
	}
	bootController.Program = Program{
		Idiom: b.Idiom,
		Repo:  b.Repo,
		System: System{
			Faculty:  b.Booter.Reserve,
			Boundary: bootController,
			Combiner: bootController,
		},
	}
	return bootController.Boot()
}

type BootController struct {
	Origin  *Span   `ko:"name=origin"` // evaluation span (not boot span)
	Booter  *Booter `ko:"name=booter"`
	Func    *Func   `ko:"name=func"`
	Ctx     Symbol  `ko:"name=ctx"`
	Arg     Symbol  `ko:"name=arg"`
	Program Program `ko:"name=program"`
}

// boot forwards eval panics from the caller evaluator environment.
func (b *BootController) Boot() (returned, effect Symbol, err error) {
	// evaluation strategy is sequential
	if shape, effect, err := b.Program.EvalSeq(NewSpan(), b.Func, b.Arg); err != nil {
		return nil, nil, err
	} else {
		return b.Unwrap(shape.(BootObject)), b.UnwrapEffect(effect.(BootEffect)), nil
	}
}

func (b *BootController) BootStepCtx(bootSpan *Span) *BootStepCtx {
	bootStep := NearestStep(bootSpan)
	return &BootStepCtx{
		Origin: b.Origin,
		Pkg:    b.Func.Pkg,
		Func:   b.Func.Name,
		Step:   bootStep.Label,
		Logic:  bootStep.Logic.String(),
		Source: bootStep.RegionString(),
		Ctx:    b.Ctx,
	}
}

func (b *BootController) BootSummary(returned Symbol) *BootSummary {
	return &BootSummary{
		Origin:   b.Origin,
		Pkg:      b.Func.Pkg,
		Func:     b.Func.Name,
		Source:   b.Func.RegionString(),
		Ctx:      b.Ctx,
		Arg:      b.Arg,
		Returned: returned,
	}
}
