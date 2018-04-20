package model

import (
	"bytes"
	"strings"
)

func (f *Func) DocShort() string {
	var w bytes.Buffer
	w.WriteString(f.Name)
	w.WriteString("(")
	args := f.Args()
	args2 := make([]string, len(args))
	for i, a := range args {
		if f.Monadic == a {
			args2[i] = a + "?"
		} else {
			args2[i] = a
		}
	}
	w.WriteString(strings.Join(args2, ", "))
	w.WriteString(")")
	return w.String()
}

func (f *Func) DocLong() string {
	var w bytes.Buffer
	w.WriteString(f.DocShort())
	w.WriteString("\n\n")
	w.WriteString(f.Doc)
	return w.String()
}
