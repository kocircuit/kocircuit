package compile

import (
	"fmt"
	"path"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func MustCompileString(forPkg, fileName, fileText string) Repo {
	repo, err := CompileString(forPkg, fileName, fileText)
	if err != nil {
		panic("o")
	}
	return repo
}

func CompileString(forPkg, fileName, fileText string) (Repo, error) {
	fileSyntax, err := ParseFileString(fileName, fileText)
	if err != nil {
		return nil, fmt.Errorf("parsing (%v)", err)
	}
	pkg, err := GraftFile(forPkg, fileSyntax)
	if err != nil {
		return nil, fmt.Errorf("grafting (%v)", err)
	}
	if err = pkg.SweepSteps(
		func(step *Step) error {
			return resolvePkgScopeFuncRef(forPkg, pkg, step)
		},
	); err != nil {
		return nil, err
	}
	return Repo{forPkg: pkg}, nil
}

// Step logics include:
//	+ PkgFunc
//	+ Operator with other-than-two reference elements
func GraftFile(pkgPath string, file File) (pkg Package, err error) {
	pkg = Package{}
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
		func(step *Step) error {
			return rewritePkgAliasRef(file.Path, asPkg, step)
		},
	); err != nil {
		return nil, err
	}
	return pkg, nil
}

// rewritePkgAliasRef converts Operator into PkgFunc logics.
func rewritePkgAliasRef(filePath string, asPkg map[string]string, u *Step) error {
	switch ref := u.Logic.(type) {
	case Operator:
		if len(ref.Path) == 2 {
			pkg, ok := asPkg[ref.Path[0]]
			if !ok {
				return fmt.Errorf("%s not known at %s", ref.Path[0], u.RegionString())
			}
			u.Logic = PkgFunc{pkg, ref.Path[1]}
		}
	}
	return nil
}

func graftAsPkgMap(imp []Import) (pkg map[string]string, err error) {
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
