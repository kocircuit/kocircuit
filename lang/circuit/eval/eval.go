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

package eval

import (
	"github.com/kocircuit/kocircuit/lang/circuit/flow"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Evaluator interface {
	String() string
	EvalPar(*model.Span, *model.Func, Arg) (Return, Effect, error)
	EvalSeq(*model.Span, *model.Func, Arg) (Return, Effect, error)
	EvalIdiom() model.Repo
	EvalRepo() model.Repo
}

func SpanEvaluator(span *model.Span) Evaluator {
	return span.Hypervisor.(Evaluator)
}

// Program is an Evaluator.
type Program struct {
	Idiom  model.Repo `ko:"name=idiom"`
	Repo   model.Repo `ko:"name=repo"` // function repository from user sources
	System `ko:"name=system"`
}

type System struct {
	Faculty  Faculty  `ko:"name=faculty"`  // operator logic
	Boundary Boundary `ko:"name=boundary"` // interprets literals
	Combiner Combiner `ko:"name=combiner"` // generates function effect/description
}

func (prog Program) String() string { return tree.Sprint(prog) }

func (prog Program) EvalIdiom() model.Repo {
	return prog.Idiom
}

func (prog Program) EvalRepo() model.Repo {
	return prog.Repo
}

func (f evalFlow) StepResidue() *StepResidue {
	return &StepResidue{Span: f.Span, Shape: f.Shape, Effect: f.Effect}
}

func FlowResidues(stepFlow []flow.Flow) (stepResidue []*StepResidue) {
	stepResidue = make([]*StepResidue, len(stepFlow))
	for i, f := range stepFlow {
		stepResidue[i] = f.(evalFlow).StepResidue()
	}
	return
}

// func (prog Program) Eval(span *Span, f *Func, arg Arg) (Return, Effect, error) {
// 	return prog.EvalSeq(span, f, arg)
// }

func (prog Program) EvalSeq(span *model.Span, f *model.Func, arg Arg) (Return, Effect, error) {
	span = span.Attach(prog)
	envelope := evalEnvelope{Program: prog, Arg: arg, Span: span}
	if returnFlow, stepFlow, err := flow.PlaySeqFlow(span, f, envelope); err != nil {
		return nil, nil, err
	} else {
		returned := returnFlow.(evalFlow).Shape
		if eff, err := prog.Combiner.Combine(span, f, arg, returned, FlowResidues(stepFlow)); err != nil {
			return nil, nil, err
		} else {
			return returned, eff, nil
		}
	}
}

func (prog Program) EvalPar(span *model.Span, f *model.Func, arg Arg) (Return, Effect, error) {
	span = span.Attach(prog)
	envelope := evalEnvelope{Program: prog, Arg: arg, Span: span}
	if returnFlow, stepFlow, err := flow.PlayParFlow(span, f, envelope); err != nil {
		return nil, nil, err
	} else {
		returned := returnFlow.(evalFlow).Shape
		if eff, err := prog.Combiner.Combine(span, f, arg, returned, FlowResidues(stepFlow)); err != nil {
			return nil, nil, err
		} else {
			return returned, eff, nil
		}
	}
}
