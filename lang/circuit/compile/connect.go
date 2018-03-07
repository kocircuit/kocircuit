package compile

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func (g *grafting) graftStep(label string, labelRegion Region, asm Assembly) ([]*Step, error) {
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
func (g *grafting) graftFieldLabelOrDefer(label string, labelRegion Region, ref Ref) ([]*Step, error) {
	if len(ref.Path) == 0 {
		return []*Step{
			g.makeDeferRefStep(label, labelRegion, ref),
		}, nil // not an arg or label reference
	}
	if fieldStep, ok := g.arg[ref.Path[0]]; ok {
		return g.selectFrom(ref.Lex, label, []*Step{fieldStep}, ref.Path[1:]), nil
	}
	if _, ok := g.labelTerm[ref.Path[0]]; ok {
		labelStep, err := g.graftLabel(ref.Path[0])
		if err != nil {
			return nil, err
		}
		return g.selectFrom(ref.Lex, label, labelStep, ref.Path[1:]), nil
	}
	return []*Step{
		g.makeDeferRefStep(label, labelRegion, ref),
	}, nil // not an arg or label reference
}

func (g *grafting) graftGather(term []Term) (gather []*Gather, err error) {
	gather = []*Gather{}
	for _, t := range term {
		ss, err := g.graftTerm("", t.Label, t)
		if err != nil {
			return nil, err
		}
		for _, s := range ss {
			gather = append(gather, &Gather{Field: t.Label.Name(), Step: s})
		}
	}
	return
}

func (g *grafting) graftTerm(label string, labelRegion Region, term Term) ([]*Step, error) {
	switch u := term.Hitch.(type) {
	case Literal:
		return []*Step{
			g.add(
				Step{
					Label:  g.takeLabel(label),
					Logic:  Number{u.Value}, // LexString, LexInteger, LexFloat
					Syntax: labelRegion,
				},
			),
		}, nil
	case Ref:
		if len(u.Path) == 1 {
			if u.Path[0] == "true" {
				return []*Step{
					g.add(
						Step{
							Label:  g.takeLabel(label),
							Logic:  Number{true},
							Syntax: labelRegion,
						},
					),
				}, nil
			}
			if u.Path[0] == "false" {
				return []*Step{
					g.add(
						Step{
							Label:  g.takeLabel(label),
							Logic:  Number{false},
							Syntax: labelRegion,
						},
					),
				}, nil
			}
		}
		return g.graftFieldLabelOrDefer(label, labelRegion, u)
	case Assembly:
		return g.graftStep(label, labelRegion, u)
	}
	return nil, fmt.Errorf("unrecognized term at %v", term)
}
