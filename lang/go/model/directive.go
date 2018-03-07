package model

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
)

type GoDirective struct {
	Label     string      `ko:"name=label"`
	GroupPath GoGroupPath `ko:"name=groupPath"`
	Inject    GoExpr      `ko:"name=expr"`
}

func (directive *GoDirective) DirectiveID() string {
	return Mix(directive.Label)
}

func (directive *GoDirective) RenderExpr(fileCtx GoFileContext) string {
	return directive.Inject.RenderExpr(fileCtx)
}

func DirectiveToExpr(directive ...*GoDirective) []GoExpr {
	r := make([]GoExpr, len(directive))
	for i, d := range directive {
		r[i] = d
	}
	return r
}

func FilterDirective(directive []*GoDirective, match GoGroupPath) []*GoDirective {
	r := []*GoDirective{}
	for _, d := range directive {
		if d.GroupPath == match {
			r = append(r, d)
		}
	}
	return r
}
