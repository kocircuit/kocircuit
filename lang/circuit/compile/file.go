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
	"path"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

// MustCompileString compiles a file with given name & content in a package with given name.
// The function with panic is case of errors.
func MustCompileString(forPkg, fileName, fileText string) model.Repo {
	repo, err := CompileString(forPkg, fileName, fileText)
	if err != nil {
		panic(fmt.Sprintf("Compile failed: %s", err))
	}
	return repo
}

// CompileString compiles a file with given name & content in a package with given name.
func CompileString(forPkg, fileName, fileText string) (model.Repo, error) {
	fileSyntax, err := syntax.ParseFileString(fileName, fileText)
	if err != nil {
		return nil, fmt.Errorf("parsing (%v)", err)
	}
	pkg, err := GraftFile(forPkg, fileSyntax)
	if err != nil {
		return nil, fmt.Errorf("grafting (%v)", err)
	}
	if err = pkg.SweepSteps(
		func(step *model.Step) error {
			return resolvePkgScopeFuncRef(forPkg, pkg, step)
		},
	); err != nil {
		return nil, err
	}
	return model.Repo{forPkg: pkg}, nil
}

// Step logics include:
//	+ PkgFunc
//	+ Operator with other-than-two reference elements
func GraftFile(pkgPath string, file syntax.File) (pkg model.Package, err error) {
	pkg = model.Package{}
	for _, parsedFunc := range file.Design {
		if _, ok := pkg[parsedFunc.Name.Name()]; ok {
			return nil, fmt.Errorf(
				"function %s.%s (%s) already grafted in the same package",
				pkgPath, parsedFunc.Name.Name(),
				parsedFunc.RegionString(),
			)
		}
		f, err := graftFunc(pkgPath, parsedFunc)
		if err != nil {
			return nil, fmt.Errorf("grafting function at %v (%v)", parsedFunc.RegionString(), err)
		}
		pkg[parsedFunc.Name.Name()] = f
	}
	// expand pkgAlias.Func references to pkgPath.Func ones
	asPkg, err := graftAsPkgMap(file.Import)
	if err != nil {
		return nil, err
	}
	if err = pkg.SweepSteps(
		func(step *model.Step) error {
			return rewritePkgAliasRef(file.Path, asPkg, step)
		},
	); err != nil {
		return nil, err
	}
	return pkg, nil
}

// rewritePkgAliasRef converts Operator into PkgFunc logics.
func rewritePkgAliasRef(filePath string, asPkg map[string]string, u *model.Step) error {
	switch ref := u.Logic.(type) {
	case model.Operator:
		if len(ref.Path) == 2 {
			pkg, ok := asPkg[ref.Path[0]]
			if !ok {
				return fmt.Errorf("%s not known at %s", ref.Path[0], u.RegionString())
			}
			u.Logic = model.PkgFunc{Pkg: pkg, Func: ref.Path[1]}
		}
	}
	return nil
}

func graftAsPkgMap(imp []syntax.Import) (pkg map[string]string, err error) {
	pkg = make(map[string]string) // pkg alias -> package path
	for _, imp := range imp {
		var as string
		if imp.As.IsEmpty() {
			_, as = path.Split(imp.Path)
		} else {
			as = imp.As.Name()
		}
		if otherPkg, ok := pkg[as]; ok {
			return nil, fmt.Errorf("alias %s already used for package %s", as, otherPkg)
		}
		pkg[as] = imp.Path
	}
	return pkg, nil
}
