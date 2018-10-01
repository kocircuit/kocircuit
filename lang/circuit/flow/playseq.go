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

import "github.com/kocircuit/kocircuit/lang/circuit/model"

type StepPlayer interface {
	PlayStep(step *model.Step, gather []GatherEdge) (returns Edge, err error)
}

// Edge stands for data sent between the output of one step and the input of another during playing.
type Edge interface{}

type GatherEdge struct {
	Field string
	Edge  Edge
}

func playSeq(f *model.Func, StepPlayer StepPlayer) (map[*model.Step]Edge, error) {
	stepReturns := map[*model.Step]Edge{}
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
