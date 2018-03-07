package flow

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type StepPlayer interface {
	PlayStep(step *Step, gather []GatherEdge) (returns Edge, err error)
}

// Edge stands for data sent between the output of one step and the input of another during playing.
type Edge interface{}

type GatherEdge struct {
	Field string
	Edge  Edge
}

func playSeq(f *Func, StepPlayer StepPlayer) (map[*Step]Edge, error) {
	stepReturns := map[*Step]Edge{}
	for i := 0; i < len(f.Step); i++ {
		s := f.Step[len(f.Step)-1-i] // iterate steps in forward time order
		gather := make([]GatherEdge, len(s.Gather))
		for j, g := range s.Gather {
			gather[j] = GatherEdge{Field: g.Field, Edge: stepReturns[g.Step]}
		}
		sReturns, err := StepPlayer.PlayStep(s, gather)
		if err != nil {
			return nil, err
		}
		stepReturns[s] = sReturns
	}
	return stepReturns, nil
}
