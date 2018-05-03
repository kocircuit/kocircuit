package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Weave struct {
	Idiom  Repo
	Repo   Repo
	Weaver *Weaver
	Func   *Func
	Ctx    Symbol
	Arg    Symbol
}

func (b *Weave) Play(origin *Span) (returns, effect Symbol, err error) {
	weaveController := &WeaveController{
		Origin: origin,
		Weaver: b.Weaver,
		Func:   b.Func,
		Ctx:    b.Ctx,
		Arg:    b.Arg,
	}
	weaveController.Program = Program{
		Idiom: b.Idiom,
		Repo:  b.Repo,
		System: System{
			Faculty:  b.Weaver.Operator,
			Boundary: weaveController,
			Combiner: weaveController,
		},
	}
	return weaveController.Weave()
}

type WeaveController struct {
	Origin  *Span   `ko:"name=origin"` // evaluation span (not weave span)
	Weaver  *Weaver `ko:"name=weaver"`
	Func    *Func   `ko:"name=func"`
	Ctx     Symbol  `ko:"name=ctx"`
	Arg     Symbol  `ko:"name=arg"`
	Program Program `ko:"name=program"`
}

// weave forwards eval panics from the caller evaluator environment.
func (b *WeaveController) Weave() (returns, effect Symbol, err error) {
	// evaluation strategy is sequential
	if shape, effect, err := b.Program.EvalSeq(NewSpan(), b.Func, b.Wrap(b.Arg)); err != nil {
		return nil, nil, err
	} else {
		return MakeStructSymbol(
			FieldSymbols{
				{Name: "returns", Value: b.Unwrap(shape.(WeaveObject))},
				{Name: "effect", Value: b.UnwrapEffect(effect.(WeaveEffect))},
			},
		), nil, nil
	}
}

func (b *WeaveController) WeaveStepCtx(weaveSpan *Span) *WeaveStepCtx {
	weaveStep := NearestStep(weaveSpan)
	return &WeaveStepCtx{
		Origin: b.Origin,
		Pkg:    b.Func.Pkg,
		Func:   b.Func.Name,
		Step:   weaveStep.Label,
		Logic:  weaveStep.Logic.String(),
		Source: weaveStep.RegionString(),
		Ctx:    b.Ctx,
	}
}

func (b *WeaveController) WeaveSummary(returns Symbol) *WeaveSummary {
	return &WeaveSummary{
		Origin:  b.Origin,
		Pkg:     b.Func.Pkg,
		Func:    b.Func.Name,
		Source:  b.Func.RegionString(),
		Ctx:     b.Ctx,
		Arg:     b.Arg,
		Returns: returns,
	}
}
