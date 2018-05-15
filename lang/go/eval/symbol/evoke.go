package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type EvokeArg struct {
	Name  string `ko:"name=name"`
	Value Symbol `ko:"name=value"`
}

func Evoke(vty *VarietySymbol, args ...EvokeArg) Symbol {
	ff := make(Fields, len(args))
	for i, arg := range args {
		ff[i] = Field{Name: arg.Name, Shape: arg.Value}
	}
	if returns, _, err := vty.Evoke(NewSpan(), ff); err != nil {
		panic(err)
	} else {
		return returns.(Symbol)
	}
}

func (vty *VarietySymbol) Evoke(span *Span, fields Fields) (Symbol, Effect, error) {
	if augmented, _, err := vty.Augment(span, fields); err != nil {
		return nil, nil, err
	} else if returns, _, err := augmented.Invoke(span); err != nil {
		return nil, nil, err
	} else {
		return returns.(Symbol), nil, nil
	}
}
