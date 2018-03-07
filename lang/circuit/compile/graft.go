package compile

import (
	"strconv"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

type grafting struct {
	in            *Step
	arg           map[string]*Step // steps returning the function argument fields
	labelTerm     map[string][]Term
	pendingLabel  []string           // labels whose grafting hasn't started yet
	graftingLabel map[string]bool    // labels being grafted (on the stack)
	labelStep     map[string][]*Step // grafted label -> step
	all           []*Step            // all steps
	auto          int
	monadic       string
}

func newGrafting() *grafting {
	return &grafting{
		labelTerm:     map[string][]Term{},
		graftingLabel: map[string]bool{},
		labelStep:     map[string][]*Step{},
	}
}

func (g *grafting) add(step Step) *Step {
	g.all = append(g.all, &step)
	return &step
}

func (g *grafting) takeLabel(d string) string {
	if d != "" {
		return d
	}
	g.auto++
	return strconv.Itoa(g.auto)
}
