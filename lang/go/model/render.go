package model

import (
	"bytes"
	"fmt"
)

func RenderGoPkgFile(pkgCtx GoPkgContext, optPkgName, fileName string, expr GoExpr) *SourceFile {
	fileCtx := pkgCtx.FileContext()
	body := expr.RenderExpr(fileCtx)
	imported := fileCtx.(*FileCtx).ImportedPkgs().Statement()
	var w bytes.Buffer
	if optPkgName == "" {
		optPkgName = pkgCtx.GroupPath().PkgName()
	}
	fmt.Fprintf(&w, "package %s\n\n", optPkgName)
	if len(imported) > 0 {
		fmt.Fprintf(&w, "import (\n")
		for _, clause := range imported {
			fmt.Fprintf(&w, "\t%s %q\n", clause.Alias, clause.PkgPath)
		}
		fmt.Fprintf(&w, ")\n\n")
	}
	fmt.Fprint(&w, body)
	return &SourceFile{
		Dir:  pkgCtx.UnifiedPath(pkgCtx.GroupPath()),
		Base: fileName,
		Body: w.String(),
	}
}

func renderCircuitFile(pkgCtx GoPkgContext, cir *GoCircuit) string {
	fileCtx := pkgCtx.FileContext()
	circuitImpl := cir.RenderImpl(fileCtx) // rendering must happen before computing imported pkgs
	imported := fileCtx.(*FileCtx).ImportedPkgs().Statement()
	return ApplyTmpl(
		circuitFileTmpl,
		M{
			"PkgComment":   fmt.Sprintf("Package %s implements circuit %s", pkgCtx.GroupPath().PkgName(), cir.Name()),
			"FuncBody":     cir.Origin.BodyString(),
			"PkgName":      pkgCtx.GroupPath().PkgName(),
			"ImportedPkgs": imported,
			"CircuitImpl":  circuitImpl,
		})
}

var circuitFileTmpl = ParseTmpl(`
{{- comment .PkgComment -}}
package {{.PkgName}}

import (
	__debug "runtime/debug"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
{{- if .ImportedPkgs}}
{{range .ImportedPkgs}}
	{{.Alias}} {{printf "%q" .PkgPath}}
{{- end}}
{{- end}}
)

/*
{{.FuncBody}}
*/

{{.CircuitImpl}}`)

func renderTypesFile(pkgCtx GoPkgContext, types []*GoAlias) string {
	fileCtx := pkgCtx.FileContext()
	m := []M{}
	for _, ductType := range types {
		if ductAlias := PeekAlias(ductType); ductAlias == nil { // peek through modifiers to alias type
			continue // skip types without aliases
		} else {
			m = append(m, M{"Comment": ductAlias.Doc(), "Def": ductAlias.RenderDef(fileCtx)})
		}
	}
	imported := fileCtx.(*FileCtx).ImportedPkgs().Statement()
	return ApplyTmpl(
		typesFileTmpl,
		M{
			"PkgComment":   fmt.Sprintf("Package %s defines types", pkgCtx.GroupPath().PkgName()),
			"PkgName":      pkgCtx.GroupPath().PkgName(),
			"DuctTypeDefs": m,
			"ImportedPkgs": imported,
		})
}

var typesFileTmpl = ParseTmpl(`
{{- comment .PkgComment -}}
package {{.PkgName}}

{{- if .ImportedPkgs}}
import (
{{range .ImportedPkgs}}
	{{.Alias}} {{printf "%q" .PkgPath}}
{{- end}}
)
{{- end}}
{{range .DuctTypeDefs}}
{{comment .Comment -}}
type {{.Def}}
{{end}}`)

func renderPlumbingFile(pkgName string) string {
	return ApplyTmpl(plumbingTmpl, M{"PkgName": pkgName})
}

var plumbingTmpl = ParseTmpl(`
package {{.PkgName}}

{{/* signals denote the sources of panics within a circuit */ -}}
type kill_signal struct{}
type upstream_signal struct{}
`)
