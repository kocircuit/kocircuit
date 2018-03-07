package model

import (
	"bytes"
	"fmt"
)

type GoFile struct {
	Line []GoExpr `ko:"name=line"`
}

func (f *GoFile) RenderForPkg(pkgCtx GoPkgContext) string {
	fileCtx := pkgCtx.FileContext()
	block := &GoBlockExpr{Line: f.Line}
	body := block.RenderExpr(fileCtx)
	imported := fileCtx.(*FileCtx).ImportedPkgs()
	var w bytes.Buffer
	fmt.Fprintf(&w, "package %s\n\n", pkgCtx.GroupPath().PkgName())
	if len(imported.Import) > 0 {
		fmt.Fprintln(&w, imported.Render())
	}
	w.WriteString(body)
	return w.String()
}

type GoImportStatement struct {
	PkgPath string `ko:"name=pkgPath"`
	Alias   string `ko:"name=alias"`
}

func (expr *GoImportStatement) Render() string {
	return fmt.Sprintf("%s %q", expr.Alias, expr.PkgPath)
}

type GoImportClause struct {
	Import []*GoImportStatement `ko:"name=import"`
}

func (expr *GoImportClause) Statement() []*GoImportStatement {
	return expr.Import
}

func (expr *GoImportClause) Render() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "import (\n")
	for _, imp := range expr.Import {
		fmt.Fprintf(&w, "\t%s\n", imp.Render())
	}
	fmt.Fprintf(&w, ")\n")
	return w.String()
}

type SortGoImportStatement []*GoImportStatement

func (c SortGoImportStatement) Len() int { return len(c) }

func (c SortGoImportStatement) Less(i, j int) bool { return c[i].PkgPath < c[j].PkgPath }

func (c SortGoImportStatement) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
