package compile

import (
	"github.com/kocircuit/kocircuit/lang/circuit/kahnsort"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func sortStep(all []*Step) []*Step {
	start := []kahnsort.Node{}
	for _, s := range subtract(all, filterUsed(all)) {
		start = append(start, s)
	}
	sorted := kahnsort.Sort(start) // out and label steps are first
	r := make([]*Step, len(sorted))
	for i, _ := range sorted {
		r[i] = sorted[i].(*Step)
	}
	return r
}

func filterUsed(step []*Step) map[*Step]bool {
	u := map[*Step]bool{}
	for _, s := range step {
		for _, g := range s.Gather {
			u[g.Step] = true
		}
	}
	return u
}

func subtract(from []*Step, what map[*Step]bool) []*Step {
	var r []*Step
	for _, s := range from {
		if !what[s] {
			r = append(r, s)
		}
	}
	return r
}
