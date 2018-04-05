package compile

import (
	"fmt"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Step logics include Operator.
func graftFunc(pkg string, parsed Design) (f *Func, err error) {
	g := newGrafting()
	if err = g.graftArgs(parsed); err != nil {
		return nil, fmt.Errorf("grafting fields of %s.%v (%v)", pkg, parsed.Name, err)
	}
	if err = g.graftBody(parsed.Returns); err != nil {
		return nil, fmt.Errorf("grafting body of %s.%v (%v)", pkg, parsed.Name.Name(), err)
	}
	returned, ok := g.labelStep["return"]
	if !ok {
		return nil, fmt.Errorf(
			"missing return in %s of %s",
			Sprint(pkg),
			Sprint(parsed.Name),
		)
	}
	var gatherReturned []*Gather
	for _, s := range returned {
		gatherReturned = append(gatherReturned, &Gather{Field: MainFlowLabel, Step: s})
	}
	f = &Func{
		Doc:     strings.Trim(parsed.Comment, " \t\n\r"),
		ID:      FuncID(pkg, parsed.Name.Name()),
		Name:    parsed.Name.Name(),
		Pkg:     pkg,
		Enter:   g.in,
		Field:   g.arg,
		Monadic: g.monadic,
		Leave:   g.add(Step{Label: "0_leave", Gather: gatherReturned, Logic: Leave{}, Syntax: parsed}),
		Step:    sortStep(g.all),
		Spread:  nil, // filled later
		Syntax:  parsed,
	}
	if err = computeStepIDs(f); err != nil {
		return nil, err
	}
	BacklinkFunc(f)
	return f, nil
}

func BacklinkFunc(fu *Func) {
	fu.Spread = map[*Step][]*Step{}
	for _, s := range fu.Step {
		s.Func = fu // add backlink to the func
		for _, g := range s.Gather {
			fu.Spread[g.Step] = append(fu.Spread[g.Step], s)
		}
	}
}

func computeStepIDs(f *Func) error {
	label := map[string]bool{}
	for _, s := range f.Step {
		if label[s.Label] {
			return fmt.Errorf("duplicate step label %s in %s", s.Label, f.FullPath())
		}
		label[s.Label] = true
		s.ID = StepID(s.Label)
	}
	return nil
}

func StepID(label string) ID {
	return StringID(label)
}

func (g *grafting) graftArgs(parsed Design) error {
	g.in = g.add(Step{Label: "0_enter", Logic: Enter{}, Syntax: parsed})
	g.arg = map[string]*Step{}
	for _, factor := range parsed.Factor {
		if factor.Monadic {
			g.monadic = factor.Name.Name()
		}
		if g.arg[factor.Name.Name()] != nil {
			return fmt.Errorf("duplicate field name %s", factor.Name.Name())
		}
		g.arg[factor.Name.Name()] = g.add(
			Step{
				Label:  fmt.Sprintf("0_enter_%s", factor.Name.Name()),
				Gather: []*Gather{{Field: MainFlowLabel, Step: g.in}},
				Logic:  Select{Path: []string{factor.Name.Name()}},
				Syntax: parsed,
			},
		)
	}
	return nil
}

func (g *grafting) graftBody(asm Assembly) error {
	if len(asm.Sign.Path) != 0 || asm.Type != "{}" {
		return fmt.Errorf("function body syntax near %v", asm)
	}
	for _, t := range asm.Term {
		if t.Label.Name() == "" {
			return fmt.Errorf("step label cannot be empty (near %s)", asm.RegionString())
		}
		g.labelTerm[t.Label.Name()] = append(g.labelTerm[t.Label.Name()], t)
		g.pendingLabel = append(g.pendingLabel, t.Label.Name())
	}
	for len(g.pendingLabel) > 0 {
		label := g.pendingLabel[0]
		g.pendingLabel = g.pendingLabel[1:]
		if _, err := g.graftLabel(label); err != nil {
			return fmt.Errorf("grafting label %s at %v (%v)", label, asm, err)
		}
	}
	return nil
}

func (g *grafting) graftLabel(l string) (step []*Step, err error) {
	// caching and cyclical references
	if step, ok := g.labelStep[l]; ok {
		return step, nil
	}
	if g.graftingLabel[l] {
		return nil, fmt.Errorf("label %s involved in cyclical reference", l)
	}
	g.graftingLabel[l] = true
	defer delete(g.graftingLabel, l)
	// grafting
	for _, t := range g.labelTerm[l] {
		s, err := g.graftTerm(t.Label.Name(), t.Label, t)
		if err != nil {
			return nil, err
		}
		g.labelStep[l] = append(g.labelStep[l], s...)
	}
	return g.labelStep[l], nil
}
