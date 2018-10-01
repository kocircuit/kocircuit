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
	"strings"

	"github.com/kocircuit/kocircuit/lang/circuit/flow"
	"github.com/kocircuit/kocircuit/lang/circuit/lex"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

type evalEnvelope struct {
	Program Program
	Arg     Arg
	Span    *model.Span
}

func (env evalEnvelope) boundary() Boundary { return env.Program.System.Boundary }

func (env evalEnvelope) newFlow(span *model.Span, r Shape, eff Effect) evalFlow {
	return evalFlow{env: env, Span: span, Shape: r, Effect: eff}
}

func (env evalEnvelope) Enter(span *model.Span) (flow.Flow, error) {
	arg, effect, err := env.boundary().Enter(span, env.Arg)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

func (env evalEnvelope) Make(span *model.Span, v interface{}) (flow.Flow, error) {
	var w Figure
	switch u := v.(type) {
	case lex.LexString:
		w = String{u.String}
	case bool: // comes from graftTerm
		w = Bool{u}
	case lex.LexInteger:
		w = Integer{u.Int64}
	case lex.LexFloat:
		w = Float{u.Float64}
	default:
		panic("unrecognized figure")
	}
	arg, effect, err := env.boundary().Figure(span, w)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

func (env evalEnvelope) MakePkgFunc(span *model.Span, pkg string, name string) (flow.Flow, error) {
	if fn := env.Program.Repo.Lookup(pkg, name); fn != nil { // user can overwrite idiomatic functions
		return env.makeMacroFigure(
			span,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if fn := env.Program.Idiom.Lookup(pkg, name); fn != nil { // otherwise hereditary idiom used
		return env.makeMacroFigure(
			span,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if macro := env.Program.Faculty[Ideal{Pkg: pkg, Name: name}]; macro != nil {
		return env.makeMacroFigure(span, macro)
	}
	return nil, span.Errorf(nil, "no function or macro %q.%s", pkg, name)
}

const IdiomRootPkg = "idiom"

func (env evalEnvelope) MakeOp(span *model.Span, ref []string) (flow.Flow, error) {
	// first check if idiomatic repo implements the operator
	// thus, idiomatic circuits overwrite hard-coded macros
	if len(ref) > 0 {
		idiomFu := ref[len(ref)-1]
		idiomPkg := strings.Join(
			append([]string{IdiomRootPkg}, ref[:len(ref)-1]...),
			".",
		)
		if idiomFn := env.Program.Idiom.Lookup(idiomPkg, idiomFu); idiomFn != nil {
			return env.makeMacroFigure(
				span,
				env.Program.System.Combiner.Interpret(env.Program, idiomFn),
			)
		}
	}
	ideal := Ideal{Name: strings.Join(ref, ".")}
	if macro := env.Program.Faculty[ideal]; macro != nil {
		return env.makeMacroFigure(span, macro)
	}
	return nil, span.Errorf(nil, "macro %v not known", ideal)
}

func (env evalEnvelope) makeMacroFigure(span *model.Span, macro Macro) (flow.Flow, error) {
	arg, effect, err := env.boundary().Figure(span, macro)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

type evalFlow struct {
	env    evalEnvelope
	Shape  Shape
	Effect Effect
	Span   *model.Span // span during when this flow was created
}

func (f evalFlow) Link(span *model.Span, name string, monadic bool) (flow.Flow, error) {
	returns, effect, err := f.Shape.Link(span, name, monadic)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Select(span *model.Span, path []string) (flow.Flow, error) {
	if len(path) == 0 {
		return f, nil
	}
	returns, effect, err := f.Shape.Select(span, model.Path(path))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Augment(span *model.Span, gather []flow.GatherFlow) (flow.Flow, error) {
	returns, effect, err := f.Shape.Augment(span, gatherFlowArg(gather))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func gatherFlowArg(gather []flow.GatherFlow) Fields {
	var s Fields
	for _, g := range gather {
		fieldFlow := g.Flow.(evalFlow)
		s = append(s, Field{
			Name:  g.Field,
			Shape: fieldFlow.Shape,
		})
	}
	return s
}

func (f evalFlow) Invoke(span *model.Span) (flow.Flow, error) {
	returns, effect, err := f.Shape.Invoke(span)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Leave(span *model.Span) (flow.Flow, error) {
	r, effect, err := f.env.boundary().Leave(span, f.Shape)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, r, effect), nil
}
