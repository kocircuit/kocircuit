package model

import (
	"bytes"
	"fmt"
	"sort"

	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func CombineRepo(first, second Repo) Repo {
	r := Repo{}
	for k, v := range first {
		r[k] = v
	}
	for k, v := range second {
		r[k] = v
	}
	return r
}

// Repo captures the chain transform: package path -> function name -> function.
type Repo map[string]Package

func (repo Repo) PkgNames() []string {
	r := make([]string, 0, len(repo))
	for s := range repo {
		r = append(r, s)
	}
	return r
}

func (repo Repo) Has(pkg string) bool {
	_, has := repo[pkg]
	return has
}

func (repo Repo) StringTable(header string) [][]string {
	ss := [][]string{}
	for p, pkg := range repo {
		for f := range pkg {
			ss = append(ss, []string{header, p, f, "circuit"})
		}
	}
	return SortStringTable(ss)
}

func (repo Repo) DocPackage(pkgPath string) (string, bool) {
	if pkg, ok := repo[pkgPath]; ok {
		return pkg.DocPackage(), true
	} else {
		return "", false
	}
}

func (repo Repo) DocFunc(pkgPath, funcName string) (string, bool) {
	if pkg, ok := repo[pkgPath]; ok {
		return pkg.DocFunc(funcName)
	} else {
		return "", false
	}
}

func (repo Repo) SortedPackagePaths() []string {
	pkgPath := make([]string, 0, len(repo))
	for p := range repo {
		pkgPath = append(pkgPath, p)
	}
	sort.Strings(pkgPath)
	return pkgPath
}

func (repo Repo) String() string {
	var w bytes.Buffer
	for _, pkgPath := range repo.SortedPackagePaths() {
		fmt.Fprintln(&w, repo[pkgPath].String())
	}
	return w.String()
}

func (repo Repo) BodyString() string {
	var w bytes.Buffer
	for _, pkgPath := range repo.SortedPackagePaths() {
		fmt.Fprintln(&w, repo[pkgPath].BodyString())
	}
	return w.String()
}

func (repo Repo) Lookup(pkg string, fu string) *Func {
	if repo == nil {
		return nil
	}
	if p := repo[pkg]; p == nil {
		return nil
	} else {
		return p[fu]
	}
}

type Package map[string]*Func

func (pkg Package) DocPackage() string {
	var w bytes.Buffer
	for _, name := range pkg.SortedFuncNames() {
		w.WriteString(pkg[name].DocShort())
		w.WriteString("\n")
	}
	return w.String()
}

func (pkg Package) DocFunc(name string) (string, bool) {
	if fu := pkg[name]; fu == nil {
		return "", false
	} else {
		return fu.DocLong(), true
	}
}

func (pkg Package) SortedFuncNames() []string {
	names := make([]string, 0, len(pkg))
	for n := range pkg {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func (pkg Package) String() string {
	var w bytes.Buffer
	for _, fuName := range pkg.SortedFuncNames() {
		fmt.Fprintln(&w, pkg[fuName].String())
	}
	return w.String()
}

func (pkg Package) BodyString() string {
	var w bytes.Buffer
	for _, fuName := range pkg.SortedFuncNames() {
		fmt.Fprintln(&w, pkg[fuName].BodyString())
	}
	return w.String()
}

func (pkg Package) SweepSteps(sweep func(*Step) error) error {
	for _, f := range pkg {
		for _, s := range f.Step {
			if err := sweep(s); err != nil {
				return err
			}
		}
	}
	return nil
}
