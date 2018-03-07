package model

import (
	"fmt"
	"path"
	"sort"
	"unicode"
)

type GoRenderContext interface {
	// Unify returns the unified path for the package group, in this rendering context.
	UnifiedPath(GoGroupPath) (unified string)
	PkgContext(pkgGroup GoGroupPath) GoPkgContext
}

type GoPkgContext interface {
	GoRenderContext
	// PkgGroup returns the non-unified Go package group and path of the context package.
	GroupPath() GoGroupPath
	FileContext() GoFileContext
	Walk(string) GoPkgContext
}

type GoFileContext interface {
	GoPkgContext
	// Import returns the import alias for the requested Go package group/path.
	Import(GoGroupPath) (alias string)
}

type RenderCtx struct {
	Root string `ko:"name=root"` // pkg root relative to GOROOT
}

// UnifiedPath returns the GOROOT-relative path of the pkgPlace.
func (ctx *RenderCtx) UnifiedPath(pkgPlace GoGroupPath) string {
	switch pkgPlace.Group {
	case KoPkgGroup:
		return path.Join(ctx.Root, pkgPlace.Path) // go generated from ko
	case GoHereditaryPkgGroup:
		return pkgPlace.Path
	}
	panic("o")
}

func (ctx *RenderCtx) PkgContext(pkgGroup GoGroupPath) GoPkgContext {
	return &PkgCtx{RenderCtx: ctx, pkgGroup: pkgGroup}
}

type PkgCtx struct {
	*RenderCtx
	pkgGroup GoGroupPath // package group of the package of this file context
}

func (pkgCtx *PkgCtx) GroupPath() GoGroupPath {
	return pkgCtx.pkgGroup
}

func (pkgCtx *PkgCtx) FileContext() GoFileContext {
	return &FileCtx{
		PkgCtx:   pkgCtx,
		aliasFor: map[GoGroupPath]string{},
	}
}

func (pkgCtx *PkgCtx) Walk(subPkg string) GoPkgContext {
	return &PkgCtx{
		RenderCtx: pkgCtx.RenderCtx,
		pkgGroup: GoGroupPath{
			Group: pkgCtx.pkgGroup.Group,
			Path:  path.Join(pkgCtx.pkgGroup.Path, subPkg),
		},
	}
}

type FileCtx struct {
	*PkgCtx
	aliasFor map[GoGroupPath]string // package group to import alias
}

func (fileCtx *FileCtx) Import(pkgGroup GoGroupPath) string {
	alias, ok := fileCtx.aliasFor[pkgGroup]
	if !ok {
		alias = fmt.Sprintf(
			"pkg_%s_%d",
			SanitizeForPkgName(pkgGroup.PkgName()),
			len(fileCtx.aliasFor)+1,
		)
		fileCtx.aliasFor[pkgGroup] = alias
	}
	return alias
}

func SanitizeForPkgName(s string) string {
	buf := []byte(s)
	for i, b := range buf {
		if unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b)) || b == '_' {
			continue
		} else {
			buf[i] = '_'
		}
	}
	return string(buf)
}

func (fileCtx *FileCtx) ImportedPkgs() *GoImportClause {
	if len(fileCtx.aliasFor) == 0 {
		return &GoImportClause{}
	}
	stmt := make([]*GoImportStatement, 0, len(fileCtx.aliasFor))
	for pkgGroup, alias := range fileCtx.aliasFor {
		stmt = append(stmt,
			&GoImportStatement{PkgPath: fileCtx.UnifiedPath(pkgGroup), Alias: alias},
		)
	}
	sort.Sort(SortGoImportStatement(stmt))
	return &GoImportClause{Import: stmt}
}
