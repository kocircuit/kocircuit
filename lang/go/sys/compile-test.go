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
	"strings"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
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
	// Find test functions
	tests := &symbol.FieldSymbol{
		Name:    "tests",
		Monadic: true,
		Value:   nil,
	}
	for _, pkgName := range compiled.Repo.SortedPackagePaths() {
		pkg := compiled.Repo[pkgName]
		for _, fName := range pkg.SortedFuncNames() {
			if strings.HasPrefix(fName, "Test") {
				//f := pkg[fName] // of type model.Func
				// How to make a Variety Symbol from `f` ?
				//vSym := symbol.MakeVarietySymbol(f, nil)
				// TODO add vSym as (name: fName, func: vSym) pair to tests.Value.
			}
		}
	}
	// Run the tests
	w := &Play{
		Pkg:     "idiom",
		Func:    "RunTests",
		Repo:    compiled.Repo,
		Faculty: arg.Faculty,
		Arg:     symbol.MakeStructSymbol(symbol.FieldSymbols{tests}),
	}
	return w.Play(ctx)
}
