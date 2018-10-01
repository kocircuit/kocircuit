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
	"bytes"
	"fmt"
	"sort"

	"github.com/kocircuit/kocircuit/lang/go/kit/util"
)

// CombineRepo creates a new Repo containing all packages
// from the given Repo's combined.
// If a package is contained in both Repo's, the package from
// the second repo wins.
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

// PkgNames returns all package names in the given Repo.
func (repo Repo) PkgNames() []string {
	r := make([]string, 0, len(repo))
	for s := range repo {
		r = append(r, s)
	}
	return r
}

// Has returns true if the given package is contained in the given Repo
// or false otherwise.
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
	return util.SortStringTable(ss)
}

// DocPackage returns the documentation for the package with given name
// in the given repo.
// The returned boolean is true when the package is found or false otherwise.
func (repo Repo) DocPackage(pkgPath string) (string, bool) {
	if pkg, ok := repo[pkgPath]; ok {
		return pkg.DocPackage(), true
	} else {
		return "", false
	}
}

// DocFunc returns the documentation for the function with given name in the package with given name
// in the given repo.
// The returned boolean is true when the function is found or false otherwise.
func (repo Repo) DocFunc(pkgPath, funcName string) (string, bool) {
	if pkg, ok := repo[pkgPath]; ok {
		return pkg.DocFunc(funcName)
	} else {
		return "", false
	}
}

// SortedPackagePaths returns a sorted list of the names of all package
// contained in the given Repo.
func (repo Repo) SortedPackagePaths() []string {
	list := repo.PkgNames()
	sort.Strings(list)
	return list
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

// Lookup a function with given package name and given function name.
// Returns nil if not found.
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

// Package is a collection of functions identified by name.
type Package map[string]*Func

// DocPackage returns the documentation of the given package.
func (pkg Package) DocPackage() string {
	var w bytes.Buffer
	for _, name := range pkg.SortedFuncNames() {
		w.WriteString(pkg[name].DocShort())
		w.WriteString("\n")
	}
	return w.String()
}

// DocFunc returns the documentation for the function with given name in the given package.
// The returned boolean is true when the function is found or false otherwise.
func (pkg Package) DocFunc(name string) (string, bool) {
	if fu := pkg[name]; fu != nil {
		return fu.DocLong(), true
	}
	return "", false
}

// SortedFuncNames returns a sort list of function names in the given package.
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
