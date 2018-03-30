package compile

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func CompileRepo(repoDir, pkgPath string) (repo Repo, err error) {
	local := NewLocalRepository(repoDir)
	parsedPkgFiles, err := ParseRepo(local, pkgPath)
	if err != nil {
		return nil, err
	}
	return GraftRepo(parsedPkgFiles)
}

// Step logics include Operator logics with 0, 1, 3 or more elements, as well as PkgFunc logics.
func GraftRepo(pkgFiles map[string][]File) (repo Repo, err error) {
	repo = Repo{}
	for pkgPath, file := range pkgFiles { // for each package
		repo[pkgPath] = Package{}
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
			func(step *Step) error {
				return resolvePkgScopeFuncRef(pkgPath, repo[pkgPath], step)
			},
		); err != nil {
			return nil, err
		}
	}
	return repo, nil
}

func unionPkg(x, y Package) (Package, error) {
	z := Package{}
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

func resolvePkgScopeFuncRef(pkg string, pkgFunc Package, u *Step) error {
	switch ref := u.Logic.(type) {
	case Operator:
		if len(ref.Path) == 1 {
			if f := pkgFunc[ref.Path[0]]; f != nil {
				u.Logic = PkgFunc{Pkg: pkg, Func: f.Name}
			}
		}
	}
	return nil
}
