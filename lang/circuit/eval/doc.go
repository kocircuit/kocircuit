package eval

import (
	"bytes"
)

func (f Faculty) DocPackage(pkgPath string) (string, bool) {
	var w bytes.Buffer
	found := false
	for _, ideal := range f.SortedIdeals() {
		if ideal.Pkg == pkgPath {
			w.WriteString(ideal.Name)
			w.WriteString("\n")
			found = true
		}
	}
	if found {
		return w.String(), true
	} else {
		return "", false
	}
}

func (f Faculty) DocFunc(pkgPath, funcName string) (string, bool) {
	if macro, ok := f[Ideal{Pkg: pkgPath, Name: funcName}]; ok {
		return macro.Doc(), true
	} else {
		return "", false
	}
}
