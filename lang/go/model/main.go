package model

import (
	"bytes"
	"fmt"
)

type GoMainExpr struct {
	Valve     *GoValve       `ko:"name=valve"`
	Directive []*GoDirective `ko:"name=directive"`
}

//		// package main
//		// import "github.com/kocircuit/kocircuit/lang/ko/cmd"
//		func main() {
//			cmd.Execute()
//		}
func (expr *GoMainExpr) RenderExpr(fileCtx GoFileContext) string {
	mainAddr := &GoAddress{
		GroupPath: GoGroupPath{
			Group: GoHereditaryPkgGroup,
			Path:  "github.com/kocircuit/kocircuit/lang/ko/cmd",
		},
		Name: "Execute",
	}
	body := &GoBlockExpr{
		Line: append(
			[]GoExpr{
				&GoMainFuncExpr{
					Line: []GoExpr{
						&GoCallExpr{Func: mainAddr, Arg: []GoExpr{}},
					},
				},
			},
			DirectiveToExpr(expr.Directive...)...,
		),
	}
	return body.RenderExpr(fileCtx)
}

type GoMainFuncExpr struct {
	Line []GoExpr `ko:"name=line"`
}

func (expr *GoMainFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "func main() {")
	bodyExpr := &GoBlockExpr{Line: expr.Line}
	if body := bodyExpr.RenderExpr(fileCtx); body != "" {
		fmt.Fprint(&w, "\n\t")
		fmt.Fprintln(&w, Indent(body))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}
