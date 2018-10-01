//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package flow

import (
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

type flowPlayer struct {
	env  Envelope
	fun  *model.Func
	span *model.Span
}

func (p *flowPlayer) PlayStep(step *model.Step, gather []GatherEdge) (returns Edge, err error) {
	// Put the step we are in on the stack. Remove it at the end of this function.
	span := model.RefineStep(p.span, step)
	// Convert edges into flows.
	mainFlow, augmentedFlow, err := splitEdgeFlow(gather)
	if err != nil {
		return nil, p.span.Errorf(err, "splitting flow")
	}
	// Perform logic-dependent flow manipulation
	switch u := step.Logic.(type) {
	case model.Enter:
		return p.env.Enter(span)
	case model.Leave:
		return mainFlow.Leave(span)
	case model.Number:
		return p.env.Make(span, u.Value)
	case model.Operator:
		return p.env.MakeOp(span, u.Path)
	case model.PkgFunc:
		return p.env.MakePkgFunc(span, u.Pkg, u.Func)
	case model.Link:
		return mainFlow.Link(span, u.Name, u.Monadic)
	case model.Select:
		return mainFlow.Select(span, u.Path)
	case model.Augment:
		return mainFlow.Augment(span, augmentedFlow)
	case model.Invoke:
		return mainFlow.Invoke(span)
	}
	panic("unknown logic")
}

func splitEdgeFlow(gatherEdge []GatherEdge) (main Flow, remainder []GatherFlow, err error) {
	for _, ge := range gatherEdge {
		if ge.Field == model.MainFlowLabel {
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
