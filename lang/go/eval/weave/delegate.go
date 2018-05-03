package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

// delegate forwards invocation panics to its caller.
// errors in types are forwarded in the return value.
func (b *Weaver) delegate(delegated *Span, vty *VarietySymbol, fields Fields) (_ *WeaveStepResult, err error) {
	if vty == nil {
		return &WeaveStepResult{Returns: EmptySymbol{}, Effect: EmptySymbol{}}, nil
	}
	if stepResult, _, err := vty.Evoke(delegated, fields); err != nil {
		return nil, err
	} else {
		returns, _, err := stepResult.Select(delegated, Path{"returns"})
		if err != nil {
			return nil, delegated.Errorf(err, "expecting a structure with a returns field")
		}
		effect, _, err := stepResult.Select(delegated, Path{"effect"})
		if err != nil {
			return nil, delegated.Errorf(err, "expecting a structure with an effect field")
		}
		return &WeaveStepResult{Returns: returns.(Symbol), Effect: effect.(Symbol)}, nil
	}
}
