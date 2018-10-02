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
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

// CompileRepo compiles all files in the given packages in a repository
// contained in given directories into a Repo.
func CompileRepo(repoDirs []string, pkgPaths []string) (repo model.Repo, err error) {
	local := syntax.NewLocalRepository(repoDirs)
	parsedPkgFiles, err := syntax.ParseRepo(local, pkgPaths)
	if err != nil {
		return nil, err
	}
	return GraftRepo(parsedPkgFiles)
}

// Step logics include Operator logics with 0, 1, 3 or more elements, as well as PkgFunc logics.
func GraftRepo(pkgFiles map[string][]syntax.File) (repo model.Repo, err error) {
	repo = model.Repo{}
	for pkgPath, file := range pkgFiles { // for each package
		repo[pkgPath] = model.Package{}
		for _, file := range file { // and each file in the package directory
			fileFunc, err := GraftFile(pkgPath, file)
			if err != nil {
				return nil, err
			}
			if repo[pkgPath], err = unionPkg(repo[pkgPath], fileFunc); err != nil {
				return nil, fmt.Errorf("duplicate function %v in package %s", err, pkgPath)
			}
		}
		if err = repo[pkgPath].SweepSteps(
			func(step *model.Step) error {
				return resolvePkgScopeFuncRef(pkgPath, repo[pkgPath], step)
			},
		); err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func unionPkg(x, y model.Package) (model.Package, error) {
	z := model.Package{}
	for k, v := range x {
		z[k] = v
	}
	for k, v := range y {
		if _, ok := z[k]; ok {
			return nil, fmt.Errorf("%s", k)
		}
		z[k] = v
	}
	return z, nil
}

func resolvePkgScopeFuncRef(pkg string, pkgFunc model.Package, u *model.Step) error {
	switch ref := u.Logic.(type) {
	case model.Operator:
		if len(ref.Path) == 1 {
			if f := pkgFunc[ref.Path[0]]; f != nil {
				u.Logic = model.PkgFunc{Pkg: pkg, Func: f.Name}
			}
		}
	}
	return nil
}
