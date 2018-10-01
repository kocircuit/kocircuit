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

package model

import (
	"fmt"
	"sort"
	"strings"

	"github.com/kocircuit/kocircuit/lang/circuit/kahnsort"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
)

// Func describes a circuit.
type Func struct {
	Doc  string `ko:"name=doc"`
	ID   ID     `ko:"name=id"`
	Name string `ko:"name=name"`
	Pkg  string `ko:"name=pkg"`
	// Slice of all steps in complete topological order, aligning with the gather relationship.
	// Equivalently, this is a reverse-time ordering of the steps.
	Step []*Step `ko:"name=step"`
	// Spread maps a step to its downstream users.
	Spread map[*Step][]*Step `ko:"name=spread"`
	// Enter points to a step included in Step.
	Enter *Step `ko:"name=enter"`
	// Field points to steps included in Step.
	Field map[string]*Step `ko:"name=field"`
	// Name of monadic argument (if not empty)
	Monadic string `ko:"name=monadic"`
	// Leave points to a step included in Step.
	Leave         *Step `ko:"name=leave"`
	syntax.Syntax `ko:"name=syntax"`
}

func (f *Func) StepByLabel(label string) []*Step {
	found := []*Step{}
	for _, s := range f.Step {
		if s.Label == label {
			found = append(found, s)
		}
	}
	return found
}

// IsGlobal returns true if Name starts with a capital letter non-underscore character.
func (f *Func) IsGlobal() bool {
	return strings.ToLower(f.Name[:1]) != f.Name[:1]
}

func FuncID(pkgPath, name string) ID {
	return StringID(FuncFullPath(pkgPath, name))
}

func FuncFullPath(pkgPath, name string) string {
	return strings.Join([]string{pkgPath, name}, ".")
}

func (f *Func) FullPath() string {
	return FuncFullPath(f.Pkg, f.Name)
}

func (f *Func) FuncID() ID {
	return f.ID
}

func (f *Func) SheathID() *ID {
	return PtrID(f.FuncID())
}

func (f *Func) SheathLabel() *string {
	return util.PtrString(f.Name)
}

func (f *Func) SheathString() *string {
	return nil
}

func (f *Func) Args() []string {
	args := make([]string, 0, len(f.Field))
	for k := range f.Field {
		args = append(args, k)
	}
	sort.Strings(args)
	return args
}

// Step describes a step in a circuit.
type Step struct {
	ID            ID        `ko:"name=id"`
	Label         string    `ko:"name=label"`
	Gather        []*Gather `ko:"name=gather"`
	Logic         Logic     `ko:"name=logic"`
	syntax.Syntax `ko:"name=syntax"`
	Func          *Func `ko:"name=func"` // backlink to func
}

func (s *Step) StepID() ID {
	return s.ID
}

func (s *Step) SheathID() *ID {
	return PtrID(s.StepID())
}

func (s *Step) SheathLabel() *string {
	return util.PtrString(s.Label)
}

func (s *Step) SheathString() *string {
	return util.PtrString(fmt.Sprintf("%s:%s", s.Func.Name, s.Label))
}

// Gather describes a link between two steps.
type Gather struct {
	Field string `ko:"name=field"`
	// output of steps feed into this parameter
	Step *Step `ko:"name=step"`
}

func (s *Step) Down() []kahnsort.Node {
	var r []kahnsort.Node
	for _, g := range s.Gather {
		r = append(r, g.Step)
	}
	return r
}

// Logics are designators of transformation logic. Logics are the types:
// 	Enter{}, Leave{}, Number{}, Select{}, Augment{}, Invoke{}
// Unresolved logics:
// 	Operator{}, PkgFunc{}
type Logic interface {
	String() string
}

// Operator is a logic described by a syntactic reference.
type Operator struct {
	Path []string `ko:"name=path"`
}

func (x Operator) String() string {
	return strings.Join(x.Path, ".")
}

// PkgFunc is a placeholder logic describing a function, to be deferred, by package and name.
// Upon resolution, it is replaced by a Defer{*Func}.
type PkgFunc struct {
	Pkg  string `ko:"name=pkg"`
	Func string `ko:"name=func"`
}

func (x PkgFunc) String() string {
	return fmt.Sprintf("%q.%s", x.Pkg, x.Func)
}

type Enter struct{}

func (x Enter) String() string {
	return "ENTER"
}

type Leave struct{}

func (x Leave) String() string {
	return "LEAVE"
}

type Number struct {
	// Value is one of: bool, LexString, LexInteger or LexFloat.
	Value interface{} `ko:"name=value"`
}

func (x Number) String() string {
	switch t := x.Value.(type) {
	case string:
		return fmt.Sprintf("NUMBER(%T:%q)", t, t)
	}
	return fmt.Sprintf("NUMBER(%T:%v)", x.Value, x.Value)
}

// MainFlowLabel is the label of the function field, carrying the “main” input flow.
// It is a symbol that cannot come from the syntactic path.
const MainFlowLabel = "█"

// Field <MainFlowLabel> carries a value to select from.
type Select struct {
	Path []string `ko:"name=path"`
}

func (x Select) String() string {
	return fmt.Sprintf("SELECT(%s)", strings.Join(x.Path, "."))
}

// Field <MainFlowLabel> carries the function argument structure.
type Link struct {
	Name    string `ko:"name=name"`
	Monadic bool   `ko:"monadic=true"`
}

func (x Link) String() string {
	monadic := ""
	if x.Monadic {
		monadic = "?"
	}
	return fmt.Sprintf("ARG(%s%s)", x.Name, monadic)
}

// Augment attaches to a lambda, in field <MainFlowLabel>, all other fields.
type Augment struct{}

func (x Augment) String() string {
	return "AUGMENT"
}

// Invoke invokes the lambda passed as field <MainFlowLabel>.
type Invoke struct{}

func (x Invoke) String() string {
	return "INVOKE"
}
