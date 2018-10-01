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
	"github.com/kocircuit/kocircuit/lang/circuit/kahnsort"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

func sortStep(all []*model.Step) []*model.Step {
	start := []kahnsort.Node{}
	for _, s := range subtract(all, filterUsed(all)) {
		start = append(start, s)
	}
	sorted := kahnsort.Sort(start) // out and label steps are first
	r := make([]*model.Step, len(sorted))
	for i := range sorted {
		r[i] = sorted[i].(*model.Step)
	}
	return r
}

func filterUsed(step []*model.Step) map[*model.Step]bool {
	u := map[*model.Step]bool{}
	for _, s := range step {
		for _, g := range s.Gather {
			u[g.Step] = true
		}
	}
	return u
}

func subtract(from []*model.Step, what map[*model.Step]bool) []*model.Step {
	var r []*model.Step
	for _, s := range from {
		if !what[s] {
			r = append(r, s)
		}
	}
	return r
}
