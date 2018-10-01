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
	"strconv"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

type grafting struct {
	in            *model.Step
	arg           map[string]*model.Step // steps returning the function argument fields
	labelTerm     map[string][]syntax.Term
	pendingLabel  []string                 // labels whose grafting hasn't started yet
	graftingLabel map[string]bool          // labels being grafted (on the stack)
	labelStep     map[string][]*model.Step // grafted label -> step
	all           []*model.Step            // all steps
	auto          int
	monadic       string
}

func newGrafting() *grafting {
	return &grafting{
		labelTerm:     map[string][]syntax.Term{},
		graftingLabel: map[string]bool{},
		labelStep:     map[string][]*model.Step{},
	}
}

func (g *grafting) add(step model.Step) *model.Step {
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
