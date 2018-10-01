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
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Combiner interface {
	Interpret(Evaluator, *model.Func) Macro
	Combine(*model.Span, *model.Func, Arg, Return, StepResidues) (Effect, error)
}

type StepResidue struct {
	Span   *model.Span `ko:"name=span"`
	Shape  Shape       `ko:"name=shape"`
	Effect Effect      `ko:"name=effect"`
}

type StepResidues []*StepResidue

func (sr StepResidues) String() string {
	return tree.Sprint(sr)
}

type IdentityCombiner struct{}

func (IdentityCombiner) Combine(_ *model.Span, _ *model.Func, _ Arg, _ Return, effect StepResidues) (Effect, error) {
	return effect, nil
}

func (IdentityCombiner) Interpret(eval Evaluator, f *model.Func) Macro {
	return &evalFixedFuncMacro{Func: f, Parent: eval}
}

// evalFixedFuncMacro is a macro that plays an underlying circuit function with the parent evaluator.
type evalFixedFuncMacro struct {
	Func   *model.Func
	Parent Evaluator
}

func (m *evalFixedFuncMacro) Splay() tree.Tree {
	return tree.Quote{String_: m.Help()}
}

func (m *evalFixedFuncMacro) MacroID() string { return m.Help() }

func (m *evalFixedFuncMacro) MacroSheathString() *string { return nil }

func (m *evalFixedFuncMacro) Label() string { return "evalfixed" }

func (m *evalFixedFuncMacro) Help() string {
	return fmt.Sprintf("Eval(%s)", m.Func.FullPath())
}

func (m *evalFixedFuncMacro) Doc() string {
	return m.Func.DocLong()
}

func (m *evalFixedFuncMacro) Invoke(span *model.Span, arg Arg) (Return, Effect, error) {
	if arg == nil {
		return m.Parent.EvalSeq(span, m.Func, nil)
	} else {
		return m.Parent.EvalSeq(span, m.Func, arg)
	}
}
