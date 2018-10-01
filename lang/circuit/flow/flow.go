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

type Envelope interface {
	Enter(*model.Span) (Flow, error)                       // Enter returns the output of the Enter step of the enclosing function.
	Make(*model.Span, interface{}) (Flow, error)           // Make ...
	MakePkgFunc(*model.Span, string, string) (Flow, error) // MakePkgFunc ...
	MakeOp(*model.Span, []string) (Flow, error)            // MakeOp ...
}

type Flow interface {
	Link(*model.Span, string, bool) (Flow, error)    // Link ...
	Select(*model.Span, []string) (Flow, error)      // Select ...
	Augment(*model.Span, []GatherFlow) (Flow, error) // Augment ...
	Invoke(*model.Span) (Flow, error)                // Invoke ...
	Leave(*model.Span) (Flow, error)                 // Leave ...
}

type GatherFlow struct {
	Field string
	Flow  Flow
}

func PlaySeqFlow(span *model.Span, f *model.Func, env Envelope) (Flow, []Flow, error) {
	p := &flowPlayer{span: model.RefineFunc(span, f), fun: f, env: env}
	stepReturn, err := playSeq(f, p)
	if err != nil {
		return nil, nil, err
	}
	return stepReturnFlow(f, stepReturn)
}

func PlayParFlow(span *model.Span, f *model.Func, env Envelope) (Flow, []Flow, error) {
	p := &flowPlayer{span: model.RefineFunc(span, f), fun: f, env: env}
	stepReturn, err := playPar(f, p)
	if err != nil {
		return nil, nil, err
	}
	return stepReturnFlow(f, stepReturn)
}

func stepReturnFlow(f *model.Func, stepReturn map[*model.Step]Edge) (Flow, []Flow, error) {
	stepFlow := []Flow{}
	for _, step := range f.Step {
		stepFlow = append(stepFlow, stepReturn[step].(Flow))
	}
	return stepReturn[f.Leave].(Flow), stepFlow, nil
}
