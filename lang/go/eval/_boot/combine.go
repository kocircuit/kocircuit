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
	if residue, err := b.Booter.Combine(b.Controller.BootStepCtx(bootSpan), b.Object, path); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.Wrap(residue.Effect), nil
	}
}

func (b *BootController) Interpret(_ Evaluator, f *Func) Macro {
	return &BootFuncMacro{Func: f}
}

type BootFuncMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Func                 *Func `ko:"name=func"`
}
