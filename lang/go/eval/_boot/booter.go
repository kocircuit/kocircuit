package boot

import (
	"fmt"
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Booter struct {
	Origin         *Span          `ko:"name=origin"`
	EnterVariety   *VarietySymbol `ko:"name=Enter"`
	LeaveVariety   *VarietySymbol `ko:"name=Leave"`
	LiteralVariety *VarietySymbol `ko:"name=Literal"`
	LinkVariety    *VarietySymbol `ko:"name=Link"`
	SelectVariety  *VarietySymbol `ko:"name=Select"`
	AugmentVariety *VarietySymbol `ko:"name=Augment"`
	InvokeVariety  *VarietySymbol `ko:"name=Invoke"`
	CombineVariety *VarietySymbol `ko:"name=Combine"`
}

func ExtractBooter(span *Span, a Symbol) (*Booter, error) {
	arg, ok := a.(*StructSymbol)
	if !ok {
		return nil, span.Errorf(nil, "booter must be a structure, got %v", a)
	}
	t := &Booter{Origin: span}
	if t.EnterVariety, ok = arg.Walk("Enter").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Enter must be a variety, got %v", arg.Walk("Enter"))
	}
	if t.LeaveVariety, ok = arg.Walk("Leave").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Leave must be a variety, got %v", arg.Walk("Leave"))
	}
	if t.LiteralVariety, ok = arg.Walk("Literal").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Literal must be a variety, got %v", arg.Walk("Literal"))
	}
	if t.LinkVariety, ok = arg.Walk("Link").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Link must be a variety, got %v", arg.Walk("Link"))
	}
	if t.SelectVariety, ok = arg.Walk("Select").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Select must be a variety, got %v", arg.Walk("Select"))
	}
	if t.AugmentVariety, ok = arg.Walk("Augment").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Augment must be a variety, got %v", arg.Walk("Augment"))
	}
	if t.InvokeVariety, ok = arg.Walk("Invoke").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Invoke must be a variety, got %v", arg.Walk("Invoke"))
	}
	if t.CombineVariety, ok = arg.Walk("Combine").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "booter Combine must be a variety, got %v", arg.Walk("Combine"))
	}
	return t, nil
}

func (b *Booter) delegateSpan(ctx *BootStepCtx, tag string) *Span {
	return RefineOutline(b.Origin, fmt.Sprintf("%s@%s", tag, ctx.Source))
}

func (b *Booter) combineSpan(summary *BootSummary, tag string) *Span {
	return RefineOutline(b.Origin, fmt.Sprintf("%s@%s", tag, summary.Source))
}

func (b *Booter) Enter(ctx *BootStepCtx, object Symbol) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "ENTER"),
		b.EnterVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Leave(ctx *BootStepCtx, object Symbol) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "LEAVE"),
		b.LeaveVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Link(ctx *BootStepCtx, object Symbol, name string, monadic bool) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "LINK"),
		b.LinkVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: BasicSymbol{name}},
			{Name: "monadic", Shape: BasicSymbol{monadic}},
		},
	)
}

func (b *Booter) Select(ctx *BootStepCtx, object Symbol, name string) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "SELECT"),
		b.SelectVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: BasicSymbol{name}},
		},
	)
}

func (b *Booter) Augment(ctx *BootStepCtx, object Symbol, fields BootFields) (*BootResidue, error) {
	delegatedSpan := b.delegateSpan(ctx, "AUGMENT")
	if deFields, err := fields.Deconstruct(b.Origin); err != nil {
		return nil, delegatedSpan.Errorf(err, "boot augmentation")
	} else {
		return b.delegate(
			delegatedSpan,
			b.AugmentVariety,
			Fields{
				{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
				{Name: "object", Shape: object},
				{Name: "fields", Shape: deFields},
			},
		)
	}
}

func (b *Booter) Invoke(ctx *BootStepCtx, object Symbol) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "INVOKE"),
		b.InvokeVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Literal(ctx *BootStepCtx, figure *BootFigure) (*BootResidue, error) {
	return b.delegate(
		b.delegateSpan(ctx, "LITERAL"),
		b.LiteralVariety,
		Fields{
			{Name: "ctx", Shape: ctx.Deconstruct(b.Origin)},
			{Name: "figure", Shape: figure.Deconstruct(b.Origin)},
		},
	)
}

func (b *Booter) Combine(summary *BootSummary, steps BootResidues) (*BootResidue, error) {
	combineSpan := b.combineSpan(summary, "COMBINE")
	if deSteps, err := steps.Deconstruct(combineSpan); err != nil {
		return nil, combineSpan.Errorf(err, "boot combining steps")
	} else {
		return b.delegate(
			combineSpan,
			b.CombineVariety,
			Fields{
				{Name: "summary", Shape: summary.Deconstruct(b.Origin)},
				{Name: "steps", Shape: deSteps},
			},
		)
	}
}
