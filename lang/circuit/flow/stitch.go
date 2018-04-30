package flow

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type flowPlayer struct {
	env  Envelope
	fun  *Func
	span *Span
}

func (p *flowPlayer) PlayStep(step *Step, gather []GatherEdge) (returns Edge, err error) {
	// Put the step we are in on the stack. Remove it at the end of this function.
	span := RefineStep(p.span, step)
	// Convert edges into flows.
	mainFlow, augmentedFlow, err := splitEdgeFlow(gather)
	if err != nil {
		return nil, p.span.Errorf(err, "splitting flow")
	}
	// Perform logic-dependent flow manipulation
	switch u := step.Logic.(type) {
	case Enter:
		return p.env.Enter(span)
	case Leave:
		return mainFlow.Leave(span)
	case Number:
		return p.env.Make(span, u.Value)
	case Operator:
		return p.env.MakeOp(span, u.Path)
	case PkgFunc:
		return p.env.MakePkgFunc(span, u.Pkg, u.Func)
	case Link:
		return mainFlow.Link(span, u.Name, u.Monadic)
	case Select:
		return mainFlow.Select(span, u.Path)
	case Augment:
		return mainFlow.Augment(span, augmentedFlow)
	case Invoke:
		return mainFlow.Invoke(span)
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
