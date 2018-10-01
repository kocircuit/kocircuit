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

package compile

import (
	"github.com/kocircuit/kocircuit/lang/circuit/lex"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func (g *grafting) selectFrom(region lex.Region, label string, step []*model.Step, path []string) []*model.Step {
	if len(path) == 0 {
		return step
	}
	r := make([]*model.Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			model.Step{
				Label:  g.takeLabel(label),
				Gather: []*model.Gather{{Field: model.MainFlowLabel, Step: s}},
				Logic:  model.Select{Path: path},
				Syntax: region,
			},
		)
	}
	return r
}

// Given a slice of input steps, expandDeferred returns a slice of corresponding steps,
// invoking the step's return value.
func (g *grafting) expandDeferred(label string, step []*model.Step) []*model.Step {
	r := make([]*model.Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			model.Step{
				Label:  g.takeLabel(label),
				Gather: []*model.Gather{{Field: model.MainFlowLabel, Step: s}},
				Logic:  model.Invoke{},
				Syntax: s.Syntax,
			},
		)
	}
	return r
}

// Given a slice of input steps, augment returns a slice of corresponding steps,
// augmented with the gather. If gather is empty, augment is identity on the steps.
func (g *grafting) augment(label string, region lex.Region, step []*model.Step, gather []*model.Gather) []*model.Step {
	if len(gather) == 0 {
		return step
	}
	r := make([]*model.Step, len(step))
	for i, s := range step {
		r[i] = g.add(
			model.Step{
				Label:  g.takeLabel(label),
				Gather: append([]*model.Gather{{Field: model.MainFlowLabel, Step: s}}, gather...),
				Logic:  model.Augment{},
				Syntax: region,
			},
		)
	}
	return r
}

func (g *grafting) makeDeferRefStep(label string, labelRegion lex.Region, ref syntax.Ref) *model.Step {
	return g.add(
		model.Step{
			Label:  g.takeLabel(label), // use label if non-empty, otherwise generate a unique label
			Logic:  model.Operator{Path: ref.Path},
			Syntax: ref, // associate a step with its label syntax
		},
	)
}
