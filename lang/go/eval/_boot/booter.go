package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Booter struct {
	Enter   *VarietySymbol `ko:"name=Enter"`
	Leave   *VarietySymbol `ko:"name=Leave"`
	Literal *VarietySymbol `ko:"name=Literal"`
	Link    *VarietySymbol `ko:"name=Link"`
	Select  *VarietySymbol `ko:"name=Select"`
	Augment *VarietySymbol `ko:"name=Augment"`
	Invoke  *VarietySymbol `ko:"name=Invoke"`
	Combine *VarietySymbol `ko:"name=Combine"`
}

func ExtractBooter(span *Span, a Symbol) (*Booter, error) {
	arg, ok := a.(*StructSymbol)
	if !ok {
		return nil, span.Errorf(nil, "translator must be a structure, got %v", a)
	}
	t := &Booter{}
	if t.Enter, ok = arg.Walk("Enter").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Enter must be a variety, got %v", arg.Walk("Enter"))
	}
	if t.Leave, ok = arg.Walk("Leave").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Leave must be a variety, got %v", arg.Walk("Leave"))
	}
	if t.Literal, ok = arg.Walk("Literal").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Literal must be a variety, got %v", arg.Walk("Literal"))
	}
	if t.Link, ok = arg.Walk("Link").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Link must be a variety, got %v", arg.Walk("Link"))
	}
	if t.Select, ok = arg.Walk("Select").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Select must be a variety, got %v", arg.Walk("Select"))
	}
	if t.Augment, ok = arg.Walk("Augment").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Augment must be a variety, got %v", arg.Walk("Augment"))
	}
	if t.Invoke, ok = arg.Walk("Invoke").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Invoke must be a variety, got %v", arg.Walk("Invoke"))
	}
	if t.Combine, ok = arg.Walk("Combine").(*VarietySymbol); !ok {
		return nil, span.Errorf(nil, "translator.Combine must be a variety, got %v", arg.Walk("Combine"))
	}
	return t, nil
}

func (b *Booter) Enter(ctx *BootStepCtx, object Symbol) *BootResidue {
	XXX
}

func (b *Booter) Leave(ctx *BootStepCtx, object Symbol) *BootResidue {
	XXX
}

func (b *Booter) Link(ctx *BootStepCtx, object Symbol, name string, monadic bool) *BootResidue {
	XXX
}

func (b *Booter) Select(ctx *BootStepCtx, object Symbol, name string) *BootResidue {
	XXX
}

func (b *Booter) Augment(ctx *BootStepCtx, object Symbol, fields []*BootField) *BootResidue {
	XXX
}

func (b *Booter) Invoke(ctx *BootStepCtx, object Symbol) *BootResidue {
	XXX
}

func (b *Booter) Literal(ctx *BootStepCtx, figure *BootFigure) *BootResidue {
	XXX
}

func (b *Booter) Combine(summary *BootSummary, stepResidues []*BootResidue) *BootResidue {
	XXX
}
