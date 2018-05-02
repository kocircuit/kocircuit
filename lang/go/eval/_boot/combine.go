package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (b *BootController) Combine(
	bootSpan *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidues StepResidues,
) (Effect, error) {
	summary := b.BootSummary(returned)
	steps := b.BootResidues(stepResidues)
	if residue, err := b.Booter.Combine(summary, steps); err != nil {
		return nil, err
	} else {
		return b.WrapEffect(residue.Effect), nil
	}
}

func (b *BootController) BootResidues(stepResidues StepResidues) BootResidues {
	bootResidues := make(BootResidues, len(stepResidues))
	for i, stepResidue := range stepResidues {
		bootStep := NearestStep(stepResidue.Span)
		bootResidues[i] = &BootResidue{
			Step:     bootStep.Label,
			Logic:    bootStep.Logic.String(),
			Source:   bootStep.RegionString(),
			Returned: stepResidue.Shape.(BootObject).Object,
			Effect:   stepResidue.Effect.(BootEffect).Effect,
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
