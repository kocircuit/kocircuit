package compile

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func (g *grafting) selectFrom(region Region, label string, step []*Step, path []string) []*Step {
	if len(path) == 0 {
		return step
	}
	r := make([]*Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			Step{
				Label:  g.takeLabel(label),
				Gather: []*Gather{{Field: MainFlowLabel, Step: s}},
				Logic:  Select{Path: path},
				Syntax: region,
			},
		)
	}
	return r
}

// Given a slice of input steps, expandDeferred returns a slice of corresponding steps,
// invoking the step's return value.
func (g *grafting) expandDeferred(label string, step []*Step) []*Step {
	r := make([]*Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			Step{
				Label:  g.takeLabel(label),
				Gather: []*Gather{{Field: MainFlowLabel, Step: s}},
				Logic:  Invoke{},
				Syntax: s.Syntax,
			},
		)
	}
	return r
}

// Given a slice of input steps, augment returns a slice of corresponding steps,
// augmented with the gather. If gather is empty, augment is identity on the steps.
func (g *grafting) augment(label string, region Region, step []*Step, gather []*Gather) []*Step {
	if len(gather) == 0 {
		return step
	}
	r := make([]*Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			Step{
				Label:  g.takeLabel(label),
				Gather: append([]*Gather{{Field: MainFlowLabel, Step: s}}, gather...),
				Logic:  Augment{},
				Syntax: region,
			},
		)
	}
	return r
}

func (g *grafting) makeDeferRefStep(label string, labelRegion Region, ref Ref) *Step {
	return g.add(
		Step{
			Label:  g.takeLabel(label), // use label if non-empty, otherwise generate a unique label
			Logic:  Operator{Path: ref.Path},
			Syntax: ref, // associate a step with its label syntax
		},
	)
}
