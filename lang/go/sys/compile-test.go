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

package sys

import (
	"fmt"
	"strings"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// CompileTest is a service that compiles a repository and runs all tests
// in it afterward.
type CompileTest struct {
	Repo    []string     `ko:"name=repo"`
	Pkg     string       `ko:"name=pkg"`
	Func    string       `ko:"name=func"`
	Faculty eval.Faculty `ko:"name=faculty"`
	Show    bool         `ko:"name=show"`
}

// Play runs the service.
func (arg *CompileTest) Play(ctx *runtime.Context) *PlayResult {
	c := &Compile{
		RepoDirs: arg.Repo,
		PkgPath:  arg.Pkg,
		Show:     arg.Show,
	}
	compiled := c.Play(ctx)
	if compiled.Error != nil {
		return &PlayResult{Error: compiled.Error}
	}
	// Build play service
	ev := go_eval.NewEvaluator(arg.Faculty, compiled.Repo)
	w := &PlayFuncEval{
		Func: ev.Program.Idiom.Lookup("idiom", "RunTests"),
		Eval:ev,
	}
	if (w.Func == nil ) {
		return &PlayResult{Error: fmt.Errorf("cannot find idiom.RunTests")}
	}
	// Find test functions
	var testsValue symbol.Symbols
	for _, pkgName := range compiled.Repo.SortedPackagePaths() {
		pkg := compiled.Repo[pkgName]
		for _, fName := range pkg.SortedFuncNames() {
			if strings.HasPrefix(fName, "Test") {
				f := pkg[fName] // of type model.Func
				vSym := symbol.MakeVarietySymbol(&evalTestFuncMacro{
					Func:   f,
					Parent: w.Eval.Program,
				}, nil)
				entry := symbol.MakeStructSymbol(symbol.FieldSymbols{
					&symbol.FieldSymbol{Name: "name", Value: symbol.MakeBasicSymbol(model.NewSpan(), fName)},
					&symbol.FieldSymbol{Name: "func", Value: vSym},
				})
				testsValue = append(testsValue, entry)
			}
		}
	}
	testsValueSym, err := symbol.MakeSeriesSymbol(model.NewSpan(), testsValue)
	if err != nil {
		return &PlayResult{Error: err}
	}
	tests := &symbol.FieldSymbol{
		Name:    "tests",
		Monadic: true,
		Value:   testsValueSym,
	}
	// Run the tests
	w.Arg = symbol.MakeStructSymbol(symbol.FieldSymbols{tests})
	return w.Play(ctx)
}

// evalTestFuncMacro is a macro that plays an underlying test function with the parent evaluator.
type evalTestFuncMacro struct {
	Func   *model.Func
	Parent eval.Evaluator
}

func (m *evalTestFuncMacro) Splay() tree.Tree {
	return tree.Quote{String_: m.Help()}
}

func (m *evalTestFuncMacro) MacroID() string { return m.Help() }

func (m *evalTestFuncMacro) MacroSheathString() *string { return nil }

func (m *evalTestFuncMacro) Label() string { return "evaltest" }

func (m *evalTestFuncMacro) Help() string {
	return fmt.Sprintf("Test(%s)", m.Func.FullPath())
}

func (m *evalTestFuncMacro) Doc() string {
	return m.Func.DocLong()
}

func (m *evalTestFuncMacro) Invoke(span *model.Span, arg eval.Arg) (eval.Return, eval.Effect, error) {
	return m.Parent.EvalSeq(span, m.Func, arg)
}
