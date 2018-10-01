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
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/ir"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/eval/macros"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// CompilePlay is a service that compiles a repository and plays it afterward.
type CompilePlay struct {
	Repo    []string             `ko:"name=repo"`
	Pkg     string               `ko:"name=pkg"`
	Func    string               `ko:"name=func"`
	Faculty eval.Faculty         `ko:"name=faculty"`
	Arg     *symbol.StructSymbol `ko:"name=arg"` // arg can be nil
	Show    bool                 `ko:"name=show"`
}

// Play runs the service.
func (arg *CompilePlay) Play(ctx *runtime.Context) *PlayResult {
	c := &Compile{
		RepoDirs: arg.Repo,
		PkgPaths:  []string{arg.Pkg},
		Show:     arg.Show,
	}
	compiled := c.Play(ctx)
	if compiled.Error != nil {
		return &PlayResult{Error: compiled.Error}
	}
	w := &Play{
		Pkg:     arg.Pkg,
		Func:    arg.Func,
		Repo:    compiled.Repo,
		Faculty: PostCompileFaculty(arg.Faculty, arg.Repo, compiled.Repo),
		Arg:     arg.Arg,
	}
	return w.Play(ctx)
}

func PostCompileFaculty(baseFaculty eval.Faculty, repoPaths []string, repo model.Repo) eval.Faculty {
	repoProto, repoProtoBytes, err := ir.SerializeEncodeRepo(repo)
	if err != nil {
		panic(err)
	}
	return eval.MergeFaculty(
		eval.Faculty{
			eval.Ideal{Pkg: "repo", Name: "Path"}:       &macros.EvalGoValueMacro{Value: repoPaths[0]},
			eval.Ideal{Pkg: "repo", Name: "Roots"}:      &macros.EvalGoValueMacro{Value: repoPaths},
			eval.Ideal{Pkg: "repo", Name: "Proto"}:      &macros.EvalGoValueMacro{Value: repoProto},
			eval.Ideal{Pkg: "repo", Name: "ProtoBytes"}: &macros.EvalGoValueMacro{Value: repoProtoBytes},
		},
		baseFaculty,
	)
}
