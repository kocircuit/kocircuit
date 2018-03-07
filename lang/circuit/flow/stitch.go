package flow

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type flowPlayer struct {
	env   Envelope
	fun   *Func
	frame *Span
}

func (p *flowPlayer) PlayStep(step *Step, gather []GatherEdge) (returns Edge, err error) {
	// Put the step we are in on the stack. Remove it at the end of this function.
	frame := RefineStep(p.frame, step)
	// Convert edges into flows.
	mainFlow, augmentedFlow, err := splitEdgeFlow(gather)
	if err != nil {
		return nil, p.frame.Errorf(err, "splitting flow")
	}
	// Perform logic-dependent flow manipulation
	switch u := step.Logic.(type) {
	case Enter:
		return p.env.Enter(frame)
	case Leave:
		return mainFlow.Leave(frame)
	case Number:
		return p.env.Make(frame, u.Value)
	case Operator:
		return p.env.MakeOp(frame, u.Path)
	case PkgFunc:
		return p.env.MakePkgFunc(frame, u.Pkg, u.Func)
	case Select:
		return mainFlow.Select(frame, u.Path)
	case Augment:
		return mainFlow.Augment(frame, augmentedFlow)
	case Invoke:
		return mainFlow.Invoke(frame)
	}
	panic("unknown logic")
}

func splitEdgeFlow(gatherEdge []GatherEdge) (main Flow, remainder []GatherFlow, err error) {
	for _, ge := range gatherEdge {
		if ge.Field == MainFlowLabel {
			if main != nil {
				return nil, nil, fmt.Errorf("multiple main flows")
			}
			main = ge.Edge.(Flow)
		} else {
			remainder = append(remainder, GatherFlow{Field: ge.Field, Flow: ge.Edge.(Flow)})
		}
	}
	return
}
