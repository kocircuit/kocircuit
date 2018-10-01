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
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/lex"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func (g *grafting) graftStep(label string, labelRegion lex.Region, asm syntax.Assembly) ([]*model.Step, error) {
	step, err := g.graftFieldLabelOrDefer("", labelRegion, asm.Sign)
	if err != nil {
		return nil, err
	}
	gather, err := g.graftGather(asm.Term)
	if err != nil {
		return nil, err
	}
	switch asm.Type {
	case "()":
		return g.expandDeferred(label, g.augment("", asm, step, gather)), nil
	case "[]":
		return g.augment(label, asm, step, gather), nil
	}
	return nil, fmt.Errorf("unknown bracket %s", asm.Type)
}

// graftFieldLabelOrDefer will return an empty step and nil error, if the reference does not resolve to a field or a label.
func (g *grafting) graftFieldLabelOrDefer(label string, labelRegion lex.Region, ref syntax.Ref) ([]*model.Step, error) {
	if len(ref.Path) == 0 {
		return []*model.Step{
			g.makeDeferRefStep(label, labelRegion, ref),
		}, nil // not an arg or label reference
	}
	if fieldStep, ok := g.arg[ref.Path[0]]; ok {
		return g.selectFrom(ref.Lex, label, []*model.Step{fieldStep}, ref.Path[1:]), nil
	}
	if _, ok := g.labelTerm[ref.Path[0]]; ok {
		labelStep, err := g.graftLabel(ref.Path[0])
		if err != nil {
			return nil, err
		}
		return g.selectFrom(ref.Lex, label, labelStep, ref.Path[1:]), nil
	}
	return []*model.Step{
		g.makeDeferRefStep(label, labelRegion, ref),
	}, nil // not an arg or label reference
}

func (g *grafting) graftGather(term []syntax.Term) (gather []*model.Gather, err error) {
	gather = []*model.Gather{}
	for _, t := range term {
		ss, err := g.graftTerm("", t.Label, t)
		if err != nil {
			return nil, err
		}
		for _, s := range ss {
			gather = append(gather, &model.Gather{Field: t.Label.Name(), Step: s})
		}
	}
	return
}

func (g *grafting) graftTerm(label string, labelRegion lex.Region, term syntax.Term) ([]*model.Step, error) {
	switch u := term.Hitch.(type) {
	case syntax.Literal:
		return []*model.Step{
			g.add(
				model.Step{
					Label:  g.takeLabel(label),
					Logic:  model.Number{Value: u.Value}, // LexString, LexInteger, LexFloat
					Syntax: labelRegion,
				},
			),
		}, nil
	case syntax.Ref:
		if len(u.Path) == 1 {
			if u.Path[0] == "true" {
				return []*model.Step{
					g.add(
						model.Step{
							Label:  g.takeLabel(label),
							Logic:  model.Number{Value: true},
							Syntax: labelRegion,
						},
					),
				}, nil
			}
			if u.Path[0] == "false" {
				return []*model.Step{
					g.add(
						model.Step{
							Label:  g.takeLabel(label),
							Logic:  model.Number{Value: false},
							Syntax: labelRegion,
						},
					),
				}, nil
			}
		}
		return g.graftFieldLabelOrDefer(label, labelRegion, u)
	case syntax.Assembly:
		return g.graftStep(label, labelRegion, u)
	}
	return nil, fmt.Errorf("unrecognized term at %v", term)
}
