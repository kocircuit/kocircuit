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

func (repo Repo) String() string {
	var w bytes.Buffer
	for _, pkg := range repo {
		fmt.Fprintln(&w, pkg.String())
	}
	return w.String()
}

func (repo Repo) BodyString() string {
	var w bytes.Buffer
	for _, pkg := range repo {
		fmt.Fprintln(&w, pkg.BodyString())
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

func (pkg Package) String() string {
	var w bytes.Buffer
	path := make([]string, 0, len(pkg)) // sort functions
	for p := range pkg {
		path = append(path, p)
	}
	sort.Strings(path)
	for _, p := range path {
		fmt.Fprintln(&w, pkg[p].String())
	}
	return w.String()
}

func (pkg Package) BodyString() string {
	var w bytes.Buffer
	path := make([]string, 0, len(pkg)) // sort functions
	for p := range pkg {
		path = append(path, p)
	}
	sort.Strings(path)
	for _, p := range path {
		fmt.Fprintln(&w, pkg[p].BodyString())
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
