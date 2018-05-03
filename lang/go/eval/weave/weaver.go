package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Weaver struct {
	Reserve        Faculty        `ko:"name=reserve"`
	EnterVariety   *VarietySymbol `ko:"name=Enter"`
	LeaveVariety   *VarietySymbol `ko:"name=Leave"`
	LiteralVariety *VarietySymbol `ko:"name=Literal"`
	LinkVariety    *VarietySymbol `ko:"name=Link"`
	SelectVariety  *VarietySymbol `ko:"name=Select"`
	AugmentVariety *VarietySymbol `ko:"name=Augment"`
	InvokeVariety  *VarietySymbol `ko:"name=Invoke"`
	CombineVariety *VarietySymbol `ko:"name=Combine"`
}

func ParseWeaver(span *Span, a Symbol) (_ *Weaver, err error) {
	arg, ok := a.(*StructSymbol)
	if !ok {
		return nil, span.Errorf(nil, "weaver must be a structure, got %v", a)
	}
	t := &Weaver{}
	if t.Reserve, err = ParseWeaveReserve(span, arg.Walk("reserve")); err != nil {
		return nil, span.Errorf(err,
			"weaver reserve must be a sequence of (pkg, name) pairs, got %v", arg.Walk("reserve"))
	}
	if t.EnterVariety, err = ParseWeaverVariety(span, arg, "Enter"); err != nil {
		return nil, err
	}
	if t.LeaveVariety, err = ParseWeaverVariety(span, arg, "Leave"); err != nil {
		return nil, err
	}
	if t.LiteralVariety, err = ParseWeaverVariety(span, arg, "Literal"); err != nil {
		return nil, err
	}
	if t.LinkVariety, err = ParseWeaverVariety(span, arg, "Link"); err != nil {
		return nil, err
	}
	if t.SelectVariety, err = ParseWeaverVariety(span, arg, "Select"); err != nil {
		return nil, err
	}
	if t.AugmentVariety, err = ParseWeaverVariety(span, arg, "Augment"); err != nil {
		return nil, err
	}
	if t.InvokeVariety, err = ParseWeaverVariety(span, arg, "Invoke"); err != nil {
		return nil, err
	}
	if t.CombineVariety, err = ParseWeaverVariety(span, arg, "Combine"); err != nil {
		return nil, err
	}
	return t, nil
}

func ParseWeaverVariety(span *Span, arg *StructSymbol, name string) (*VarietySymbol, error) {
	switch u := arg.Walk(name).(type) {
	case EmptySymbol:
		return nil, nil
	case *VarietySymbol:
		return u, nil
	case nil:
		panic("o")
	default:
		return nil, span.Errorf(nil, "weaver %s must be a variety, got %v", name, u)
	}
}

func (b *Weaver) Enter(ctx *WeaveStepCtx, object Symbol) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.EnterVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Weaver) Leave(ctx *WeaveStepCtx, object Symbol) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LeaveVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Weaver) Link(ctx *WeaveStepCtx, object Symbol, name string, monadic bool) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LinkVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: BasicSymbol{name}},
			{Name: "monadic", Shape: BasicSymbol{monadic}},
		},
	)
}

func (b *Weaver) Select(ctx *WeaveStepCtx, object Symbol, path Path) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.SelectVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
			{Name: "path", Shape: MakeStringsSymbol(delegatedSpan, []string(path))},
		},
	)
}

func (b *Weaver) Augment(ctx *WeaveStepCtx, object Symbol, fields WeaveFields) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	if deFields, err := fields.Deconstruct(ctx.Origin); err != nil {
		return nil, delegatedSpan.Errorf(err, "weave augmentation")
	} else {
		return b.delegate(
			delegatedSpan,
			b.AugmentVariety,
			Fields{
				{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
				{Name: "object", Shape: object},
				{Name: "fields", Shape: deFields},
			},
		)
	}
}

func (b *Weaver) Invoke(ctx *WeaveStepCtx, object Symbol) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.InvokeVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Weaver) Literal(ctx *WeaveStepCtx, figure *WeaveFigure) (*WeaveStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LiteralVariety,
		Fields{
			{Name: "stepCtx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "figure", Shape: figure.Deconstruct(delegatedSpan)},
		},
	)
}

func (b *Weaver) Combine(summary *WeaveSummary, steps WeaveResidues) (*WeaveStepResult, error) {
	combineSpan := summary.CombineSpan()
	if deSteps, err := steps.Deconstruct(combineSpan); err != nil {
		return nil, combineSpan.Errorf(err, "weave combining steps")
	} else {
		return b.delegate(
			combineSpan,
			b.CombineVariety,
			Fields{
				{Name: "summaryCtx", Shape: summary.Deconstruct(combineSpan)},
				{Name: "stepResidues", Shape: deSteps},
			},
		)
	}
}
