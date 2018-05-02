package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) Combine(
	bootSpan *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidue StepResidues,
) (Effect, error) {
	summary := b.Controller.BootSummary(returned)
	steps := b.Controller.BootResidues(stepResidues)
	if residue, err := b.Booter.Combine(summary, steps); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.Wrap(residue.Effect), nil
	}
}

func (b *BootController) BootResidues(stepResidue StepResidues) BootResidues {
	bootResidues := make(BootResidues, len(stepResidues))
	for i := range stepResidues {
		bootResidues[i] = &BootResidue{
			Step:     XXX, //     string `ko:"name=step"`
			Logic:    XXX, //    string `ko:"name=logic"`
			Source:   XXX, //   string `ko:"name=source"`
			Returned: XXX, // Symbol `ko:"name=returned"`
			Effect:   XXX, //   Symbol `ko:"name=effect"`
		}
	}
	return bootResidues
}

func (b *BootController) Interpret(_ Evaluator, f *Func) Macro {
	return &BootFuncMacro{Func: f}
}

type BootFuncMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Func                 *Func `ko:"name=func"`
}
