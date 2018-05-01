package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) Combine(
	span *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidue StepResidues,
) (Effect, error) {
	XXX //XXX
	return nil, nil
}

func (b *BootController) Interpret(_ Evaluator, f *Func) Macro {
	return &BootFuncMacro{Func: f}
}

type BootFuncMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Func                 *Func `ko:"name=func"`
}
