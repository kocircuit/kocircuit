package boot

import (
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
		return nil, span.Errorf(nil, "translator must be a structure, got %v", a)
	}
	t := &Booter{Origin: span}
	if t.EnterVariety, ok = arg.Walk("Enter").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Enter must be a variety, got %v", arg.Walk("Enter"))
	}
	if t.LeaveVariety, ok = arg.Walk("Leave").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Leave must be a variety, got %v", arg.Walk("Leave"))
	}
	if t.LiteralVariety, ok = arg.Walk("Literal").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Literal must be a variety, got %v", arg.Walk("Literal"))
	}
	if t.LinkVariety, ok = arg.Walk("Link").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Link must be a variety, got %v", arg.Walk("Link"))
	}
	if t.SelectVariety, ok = arg.Walk("Select").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Select must be a variety, got %v", arg.Walk("Select"))
	}
	if t.AugmentVariety, ok = arg.Walk("Augment").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Augment must be a variety, got %v", arg.Walk("Augment"))
	}
	if t.InvokeVariety, ok = arg.Walk("Invoke").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Invoke must be a variety, got %v", arg.Walk("Invoke"))
	}
	if t.CombineVariety, ok = arg.Walk("Combine").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Combine must be a variety, got %v", arg.Walk("Combine"))
	}
	return t, nil
}

func (b *Booter) Enter(ctx *BootStepCtx, object Symbol) *BootResidue {
	return b.delegate(
		b.EnterVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Leave(ctx *BootStepCtx, object Symbol) *BootResidue {
	return b.delegate(
		b.LeaveVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Link(ctx *BootStepCtx, object Symbol, name string, monadic bool) *BootResidue {
	return b.delegate(
		b.LinkVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: DeconstructInterface(b.Origin, name)},
			{Name: "monadic", Shape: DeconstructInterface(b.Origin, monadic)},
		},
	)
}

func (b *Booter) Select(ctx *BootStepCtx, object Symbol, name string) *BootResidue {
	return b.delegate(
		b.SelectVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
			{Name: "name", Shape: DeconstructInterface(b.Origin, name)},
		},
	)
}

func (b *Booter) Augment(ctx *BootStepCtx, object Symbol, fields []*BootField) *BootResidue {
	return b.delegate(
		b.AugmentVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
			{Name: "fields", Shape: DeconstructInterface(b.Origin, fields)},
		},
	)
}

func (b *Booter) Invoke(ctx *BootStepCtx, object Symbol) *BootResidue {
	return b.delegate(
		b.InvokeVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "object", Shape: object},
		},
	)
}

func (b *Booter) Literal(ctx *BootStepCtx, figure *BootFigure) *BootResidue {
	return b.delegate(
		b.LiteralVariety,
		Fields{
			{Name: "ctx", Shape: DeconstructInterface(b.Origin, ctx)},
			{Name: "figure", Shape: DeconstructInterface(b.Origin, figure)},
		},
	)
}

func (b *Booter) Combine(summary *BootSummary, steps []*BootResidue) *BootResidue {
	return b.delegate(
		b.CombineVariety,
		Fields{
			{Name: "summary", Shape: DeconstructInterface(b.Origin, summary)},
			{Name: "steps", Shape: DeconstructInterface(b.Origin, steps)},
		},
	)
}
