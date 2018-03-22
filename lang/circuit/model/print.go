package model

import (
	"bytes"
	"fmt"
	"sort"
)

const KoIndent = "\t" // used to be 3 spaces

func (f *Func) String() string {
	return fmt.Sprintf("%s.%s", f.Pkg, f.Name)
}

func (f *Func) SortedFieldNames() []string {
	fieldNames := make([]string, 0, len(f.Field))
	for p := range f.Field {
		fieldNames = append(fieldNames, p)
	}
	sort.Strings(fieldNames)
	return fieldNames
}

func (f *Func) BodyString() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%s(", fmt.Sprintf("%s.%s", f.Pkg, f.Name))
	var numField int
	sortedFieldNames := f.SortedFieldNames()
	for _, fieldName := range sortedFieldNames {
		fmt.Fprint(&w, f.Field[fieldName])
		numField++
		if numField < len(sortedFieldNames) {
			fmt.Fprint(&w, ", ")
		}
	}
	fmt.Fprintf(&w, ") {\n")
	for _, s := range f.Step {
		fmt.Fprintf(&w, "%s // %s\n", s.String(), s.Syntax.RegionString())
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

func (s *Step) String() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%s%s: %s(", KoIndent, s.Label, s.Logic.String())
	var n int
	for _, g := range s.Gather {
		fmt.Fprintf(&w, "%s:%s, ", g.Field, g.Step.Label)
		n++
	}
	if n > 0 {
		w.Truncate(w.Len() - 2)
	}
	fmt.Fprint(&w, ")")
	return w.String()
}
