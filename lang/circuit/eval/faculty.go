package eval

import (
	"bytes"
	"fmt"
	"path"
	"sort"

	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Ideals []Ideal

func (ss Ideals) Len() int { return len(ss) }

func (ss Ideals) Swap(i, j int) {
	ss[i], ss[j] = ss[j], ss[i]
}

func (ss Ideals) Less(i, j int) bool {
	if ss[i].Pkg < ss[j].Pkg {
		return true
	} else {
		return ss[i].Name < ss[j].Name
	}
}

func (ss Ideals) Sort() {
	sort.Sort(ss)
}

type Ideal struct {
	Pkg  string `ko:"name=pkg"`
	Name string `ko:"name=name"`
}

func (ideal Ideal) FullPath() string {
	return path.Join(ideal.Pkg, ideal.Name)
}

func (ideal Ideal) String() string {
	return fmt.Sprintf("%q.%s", ideal.Pkg, ideal.Name)
}

type Faculty map[Ideal]Macro

func (f Faculty) PkgNames() []string {
	n := map[string]bool{}
	for ideal := range f {
		n[ideal.Pkg] = true
	}
	r := make([]string, 0, len(n))
	for s := range n {
		r = append(r, s)
	}
	return r
}

func (f Faculty) SortedIdeals() (ideals Ideals) {
	for ideal := range f {
		ideals = append(ideals, ideal)
	}
	ideals.Sort()
	return
}

func (f Faculty) StringTable(header string) [][]string {
	ss := [][]string{}
	for ideal, macro := range f {
		ss = append(ss, []string{
			header,
			ideal.Pkg,
			ideal.Name,
			macro.Help(),
		})
	}
	return SortStringTable(ss)
}

func (f Faculty) ID() string {
	id := ""
	for k, v := range f {
		id = Mix(id, k.Pkg, k.Name, Mix(v.Help()))
	}
	return id
}

func (f Faculty) Add(key Ideal, value Macro) { f[key] = value }

func (f Faculty) AddExclusive(key Ideal, value Macro) {
	if _, ok := f[key]; ok {
		panic("duplicate faculty key")
	}
	f.Add(key, value)
}

func (f Faculty) String() string {
	line := []string{}
	for k, m := range f {
		line = append(line, fmt.Sprintf("%s: %v", k.String(), m.Help()))
	}
	sort.Strings(line)
	var w bytes.Buffer
	for _, l := range line {
		fmt.Fprintln(&w, l)
	}
	return w.String()
}

type MacroCases []Macro

func MergeFaculty(f ...Faculty) Faculty {
	r := Faculty{}
	for _, f := range f {
		for k, v := range f {
			r.Add(k, v)
		}
	}
	return r
}

func MergeFacultyExclusive(f ...Faculty) Faculty {
	r := Faculty{}
	for _, f := range f {
		for k, v := range f {
			r.AddExclusive(k, v)
		}
	}
	return r
}
