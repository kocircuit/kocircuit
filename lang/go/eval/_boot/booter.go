package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Booter struct {
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

func ParseBooter(span *Span, a Symbol) (_ *Booter, err error) {
	arg, ok := a.(*StructSymbol)
	if !ok {
		return nil, span.Errorf(nil, "booter must be a structure, got %v", a)
	}
	t := &Booter{}
	if t.Reserve, err = ParseBootReserve(span, arg.Walk("reserve")); err != nil {
		return nil, span.Errorf(err,
			"booter reserve must be a sequence of (pkg, name) pairs, got %v", arg.Walk("reserve"))
	}
	if t.EnterVariety, err = ParseBooterVariety(span, arg, "Enter"); err != nil {
		return nil, err
	}
	if t.LeaveVariety, err = ParseBooterVariety(span, arg, "Leave"); err != nil {
		return nil, err
	}
	if t.LiteralVariety, err = ParseBooterVariety(span, arg, "Literal"); err != nil {
		return nil, err
	}
	if t.LinkVariety, err = ParseBooterVariety(span, arg, "Link"); err != nil {
		return nil, err
	}
	if t.SelectVariety, err = ParseBooterVariety(span, arg, "Select"); err != nil {
		return nil, err
	}
	if t.AugmentVariety, err = ParseBooterVariety(span, arg, "Augment"); err != nil {
		return nil, err
	}
	if t.InvokeVariety, err = ParseBooterVariety(span, arg, "Invoke"); err != nil {
		return nil, err
	}
	if t.CombineVariety, err = ParseBooterVariety(span, arg, "Combine"); err != nil {
		return nil, err
	}
	return t, nil
}

func ParseBooterVariety(span *Span, arg *StructSymbol, name string) (*VarietySymbol, error) {
	switch u := arg.Walk(name).(type) {
	case EmptySymbol:
		return nil, nil
	case *VarietySymbol:
		return u, nil
	case nil:
		panic("o")
	default:
		return nil, span.Errorf(nil, "booter %s must be a variety, got %v", name, u)
	}
}

func (b *Booter) Enter(ctx *BootStepCtx, object Symbol) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.EnterVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Leave(ctx *BootStepCtx, object Symbol) (*BootStepResult, error) {
	println("LEAVE", object.String())
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LeaveVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Link(ctx *BootStepCtx, object Symbol, name string, monadic bool) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LinkVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: BasicSymbol{name}},
			{Name: "monadic", Shape: BasicSymbol{monadic}},
		},
	)
}

func (b *Booter) Select(ctx *BootStepCtx, object Symbol, path Path) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.SelectVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
			{Name: "path", Shape: MakeStringsSymbol(delegatedSpan, []string(path))},
		},
	)
}

func (b *Booter) Augment(ctx *BootStepCtx, object Symbol, fields BootFields) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	if deFields, err := fields.Deconstruct(ctx.Origin); err != nil {
		return nil, delegatedSpan.Errorf(err, "boot augmentation")
	} else {
		return b.delegate(
			delegatedSpan,
			b.AugmentVariety,
			Fields{
				{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
				{Name: "object", Shape: object},
				{Name: "fields", Shape: deFields},
			},
		)
	}
}

func (b *Booter) Invoke(ctx *BootStepCtx, object Symbol) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.InvokeVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Literal(ctx *BootStepCtx, figure *BootFigure) (*BootStepResult, error) {
	delegatedSpan := ctx.DelegateSpan()
	return b.delegate(
		delegatedSpan,
		b.LiteralVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(delegatedSpan)},
			{Name: "figure", Shape: figure.Deconstruct(delegatedSpan)},
		},
	)
}

func (b *Booter) Combine(summary *BootSummary, steps BootResidues) (*BootStepResult, error) {
	combineSpan := summary.CombineSpan()
	if deSteps, err := steps.Deconstruct(combineSpan); err != nil {
		return nil, combineSpan.Errorf(err, "boot combining steps")
	} else {
		return b.delegate(
			combineSpan,
			b.CombineVariety,
			Fields{
				{Name: "summary", Shape: summary.Deconstruct(combineSpan)},
				{Name: "steps", Shape: deSteps},
			},
		)
	}
}
