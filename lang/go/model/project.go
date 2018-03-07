package model

import (
	"fmt"
	"strconv"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// type Vty struct{
//	Branch0 *Vty0 // SEE: variety branch field
//	█
// }
//
// type VtyProj struct{
//	Projection_1_4_macro *struct{█} // SEE: line field
//	█
// }
//
// func ProjectVty(vty *Vty) *VtyProj {
// 	v0 := vty // #ProjectExprRenderingBegins
// 	switch { // #SwitchAcrossVarietyAncestors
// 	case v0.Branch0 != nil:
// 		v1 := v0.Branch0 // #ProjectExprRenderingBegins
// 		return &VtyProj{ // SEE: flatten variety to line
// 			Proj_1_4_macro: &struct{█}{
// 				A: v1.A // SEE: line branch argument name and assignment
// 				B: v2.B
// 				C: v2.C
// 				D: v3.D
// 				█
// 			}
// 		}
// 	case v0.Branch1 != nil:
// 		v1 := v0.Branch1 // (#ProjectExprRenderingBegins)
// 		█
// 	}
// }

type GoProjectionLine struct {
	Path      []int          `ko:"name=path"` // branch index path
	Root      *GoBranch      `ko:"name=root"`
	First     *GoBranch      `ko:"name=first"` // root child branch on the path to line
	Leaf      *GoBranch      `ko:"name=leaf"`  // leaf branch, pointing to macro variety
	Macro     Macro          `ko:"name=macro"`
	PathField []*GoPathField `ko:"name=pathField"` // arguments
}

type GoPathField struct {
	Path  []int    `ko:"name=path"`
	Field *GoField `ko:"name=field"`
}

func (gp *GoPathField) Depth() int { return len(gp.Path) }

func (line *GoProjectionLine) ProjectionRealField() *GoField { // SEE: line field
	return &GoField{
		Name: line.GoName(),
		Type: line.Real(),
		Tag:  KoTags(line.KoName(), false),
	}
}

func (line *GoProjectionLine) Real() GoType {
	return NewGoNeverNilPtr(NewGoStruct(FieldProject(line.PathField)...))
}

func ProjectField(path []int, field []*GoField) (projected []*GoPathField) {
	projected = make([]*GoPathField, len(field))
	for i, field := range field {
		projected[i] = &GoPathField{Path: path, Field: field}
	}
	return
}

func FieldProject(projected []*GoPathField) (field []*GoField) {
	field = make([]*GoField, len(projected))
	for i, projected := range projected {
		field[i] = projected.Field
	}
	return
}

func (line *GoProjectionLine) GoName() string {
	return fmt.Sprintf("Line_%s", line.Label())
}

func (line *GoProjectionLine) KoName() string {
	return fmt.Sprintf("line_%s", line.Label())
}

func (line *GoProjectionLine) Label() string {
	return fmt.Sprintf("%s_%s", line.PathString(), line.Macro.Label())
}

func (line *GoProjectionLine) PathString() string {
	s := make([]string, len(line.Path))
	for i, v := range line.Path {
		s[i] = strconv.Itoa(v)
	}
	return strings.Join(s, "_")
}

// ProjectShaper flattens a variety's deep structure into a flat struct.
type ProjectShaper struct {
	Shaping `ko:"name=shaping"` // From is GoVariety, To is GoNeverNil{GoStruct} (the line)
}

func (sh *ProjectShaper) ShaperID() string { return Mix("ProjectShaper", sh.ShapingID()) }

func (sh *ProjectShaper) String() string { return Sprint(sh) }

func (sh *ProjectShaper) vty() *GoVariety { return sh.Shaping.From.(*GoVariety) }

func (sh *ProjectShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *ProjectShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *ProjectShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("project_shaper_%s", sh.ShapingID()),
	}
}

func (sh *ProjectShaper) funcExpr() GoFuncExpr {
	argExpr := &GoVerbatimExpr{"vty"}
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    argExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line:       []GoExpr{sh.ProjectExpr(argExpr)},
	}
}

func (sh *ProjectShaper) CircuitEffect() *GoCircuitEffect {
	return &GoCircuitEffect{
		DuctFunc: []GoFuncExpr{sh},
	}
}

func (sh *ProjectShaper) ProjectExpr(shaperArg GoExpr) GoExpr {
	vty := sh.vty()
	root := NewRootBranch(vty)
	return root.ProjectExpr(
		&GoProjectCtx{
			Variety:       vty,
			Depth:         0,
			SwitchBranch:  ZeroExpr,
			SwitchVariety: shaperArg,
			Terminal:      BuildBranchProjectionIndex(root.ProjectionLine()),
		},
	)
}
