package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type Translator struct {
	Enter   *VarietySymbol `ko:"name=Enter"`
	Leave   *VarietySymbol `ko:"name=Leave"`
	Literal *VarietySymbol `ko:"name=Literal"`
	Link    *VarietySymbol `ko:"name=Link"`
	Select  *VarietySymbol `ko:"name=Select"`
	Augment *VarietySymbol `ko:"name=Augment"`
	Invoke  *VarietySymbol `ko:"name=Invoke"`
	Combine *VarietySymbol `ko:"name=Combine"`
}

func ExtractTranslator(span *Span, a Symbol) (*Translator, error) {
	arg, ok := a.(*StructSymbol)
	if !ok {
		return nil, span.Errorf(nil, "translator must be a structure, got %v", a)
	}
	t := &Translator{}
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
