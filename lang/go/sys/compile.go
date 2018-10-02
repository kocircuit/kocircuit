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
	"github.com/kocircuit/kocircuit/lang/circuit/compile"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// Compile is a service to compile a Ko repository.
type Compile struct {
	RepoDirs []string `ko:"name=repoDir"`
	PkgPaths []string `ko:"name=pkgPaths"`
	Show     bool     `ko:"name=show"`
}

// CompileResult holds the result of a compile service.
type CompileResult struct {
	Compile *Compile         `ko:"name=compile"`
	Repo    model.Repo       `ko:"name=repo"`
	Stats   *model.RepoStats `ko:"name=stats"`
	Error   error            `ko:"name=error"`
}

// Play runs the compilation service
func (c *Compile) Play(ctx *runtime.Context) *CompileResult {
	r := &CompileResult{Compile: c}
	if r.Repo, r.Error = compile.CompileRepo(
		c.RepoDirs,
		c.PkgPaths,
	); r.Error != nil {
		return r
	}
	r.Stats = r.Repo.Stats()
	ctx.Printf(
		"compiled functions=%d steps=%d steps-per-function=%0.2f",
		r.Stats.TotalFunc, r.Stats.TotalStep, r.Stats.StepPerFunc,
	)
	if c.Show {
		ctx.Printf("%s\n", r.Repo.BodyString())
	}
	return r
}
