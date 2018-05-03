package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (b *WeaveController) Combine(
	weaveSpan *Span,
	_ *Func,
	_ Arg, // BoolObject
	returns Return, // BoolObject
	stepResidues StepResidues,
) (Effect, error) {
	summary := b.WeaveSummary(b.Unwrap(returns))
	steps := b.WeaveResidues(stepResidues)
	if residue, err := b.Weaver.Combine(summary, steps); err != nil {
		return nil, err
	} else {
		return b.WrapEffect(residue.Effect), nil
	}
}

func (b *WeaveController) WeaveResidues(stepResidues StepResidues) WeaveResidues {
	weaveResidues := make(WeaveResidues, len(stepResidues))
	for i, stepResidue := range stepResidues {
		weaveStep := NearestStep(stepResidue.Span)
		weaveResidues[i] = &WeaveResidue{
			Step:    weaveStep.Label,
			Logic:   weaveStep.Logic.String(),
			Source:  weaveStep.RegionString(),
			Returns: stepResidue.Shape.(WeaveObject).Object,
			Effect:  stepResidue.Effect.(WeaveEffect).Effect,
		}
	}
	return weaveResidues
}

func (b *WeaveController) Interpret(_ Evaluator, f *Func) Macro {
	return &WeaveFuncMacro{Func: f}
}

type WeaveFuncMacro struct {
	WeavePlaceholderMacro `ko:"name=placeholder"`
	Func                 *Func `ko:"name=func"`
}
