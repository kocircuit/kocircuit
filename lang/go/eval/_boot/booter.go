package boot

import (
	"fmt"
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Booter struct {
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
	t := &Booter{}
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
	if deFields, err := fields.Deconstruct(b.Origin); err != nil {
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
