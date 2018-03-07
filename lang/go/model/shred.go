package model

import (
	"fmt"
)

func (ctx *RenderCtx) Shred(circuit []*GoCircuit, directive []*GoDirective, exclude ...GoGroupPath) (files []*SourceFile) {
	circuitFiles := ctx.ShredCircuits(circuit)
	files = append(files, circuitFiles...)
	directiveFiles := ctx.ShredDirectives(directive, exclude...)
	files = append(files, directiveFiles...)
	duct := ExtractDuctTypes(circuit)
	typeFiles := ctx.ShredTypes(duct)
	files = append(files, typeFiles...)
	plumbingFiles := ctx.ShredPlumbing(ExtractGroupPaths(duct))
	files = append(files, plumbingFiles...)
	return
}

func (ctx *RenderCtx) ShredPlumbing(groups []GoGroupPath) (files []*SourceFile) {
	for _, groupPath := range groups {
		pkgCtx := ctx.PkgContext(groupPath)
		files = append(files,
			&SourceFile{
				Dir:  pkgCtx.UnifiedPath(pkgCtx.GroupPath()),
				Base: "plumbing.go",
				Body: renderPlumbingFile(pkgCtx.GroupPath().PkgName()),
			},
		)
	}
	return
}

func (ctx *RenderCtx) ShredCircuits(circuit []*GoCircuit) (files []*SourceFile) {
	for _, circuitPkg := range IndexCircuits(circuit) { // group circuits by package
		pkgCtx := ctx.PkgContext(circuitPkg.GroupPath)
		for _, circuit := range circuitPkg.Circuit { // for each circuit in a package
			files = append(files,
				&SourceFile{
					Dir:  pkgCtx.UnifiedPath(pkgCtx.GroupPath()),
					Base: fmt.Sprintf("%s_circuit.go", SpanFilebase(circuit.Valve.Origin)),
					Body: renderCircuitFile(pkgCtx, circuit),
				},
			)
		}
	}
	return
}

func (ctx *RenderCtx) ShredDirectives(directive []*GoDirective, exclude ...GoGroupPath) (files []*SourceFile) {
	for _, directivePkg := range IndexDirectives(directive) {
		if GroupPathExcluded(directivePkg.GroupPath, exclude...) {
			continue
		}
		pkgCtx := ctx.PkgContext(directivePkg.GroupPath)
		line := []GoExpr{}
		for _, directive := range directivePkg.Directive {
			line = append(line, DirectiveToExpr(directive)...)
		}
		expr := &GoFile{Line: line}
		files = append(files,
			&SourceFile{
				Dir:  pkgCtx.UnifiedPath(pkgCtx.GroupPath()),
				Base: "directives.go",
				Body: expr.RenderForPkg(pkgCtx),
			},
		)
	}
	return
}

func GroupPathExcluded(gp GoGroupPath, exclude ...GoGroupPath) bool {
	for _, ex := range exclude {
		if ex == gp {
			return true
		}
	}
	return false
}

// types must not contain duplicates.
func (ctx *RenderCtx) ShredTypes(types []*GoAlias) (files []*SourceFile) {
	for _, typePkg := range IndexTypes(types) {
		pkgCtx := ctx.PkgContext(typePkg.GroupPath)
		for _, typeFile := range typePkg.File {
			files = append(files,
				&SourceFile{
					Dir:  pkgCtx.UnifiedPath(pkgCtx.GroupPath()),
					Base: fmt.Sprintf("%s_types.go", typeFile.Filebase),
					Body: renderTypesFile(pkgCtx, typeFile.Type),
				},
			)
		}
	}
	return
}
