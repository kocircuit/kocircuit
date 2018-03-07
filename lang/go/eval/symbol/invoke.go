package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (vty *VarietySymbol) Invoke(span *Span) (Shape, Effect, error) {
	if r, eff, err := vty.Macro.Invoke(
		RefineMacro(span, vty.Macro),
		MakeStructSymbol(vty.Arg),
	); err != nil {
		return nil, nil, err
	} else {
		return r.(Symbol), eff, nil
	}
}
