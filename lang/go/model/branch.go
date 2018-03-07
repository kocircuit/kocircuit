package model

import (
	"fmt"
)

type GoBranch struct {
	Index    int        `ko:"name=index"`
	Augments GoVarietal `ko:"name=augments"`
	With     []*GoField `ko:"name=with"`
}

func NewRootBranch(vty GoVarietal) *GoBranch {
	return &GoBranch{Augments: vty}
}

func (branch *GoBranch) GoName() string {
	return fmt.Sprintf("Branch%d", branch.Index+1)
}

func (branch *GoBranch) KoName() string {
	return fmt.Sprintf("branch%d", branch.Index+1)
}

type GoBranchProjectionLine map[*GoBranch]*GoProjectionLine

const AugmentsFieldName = "Augments" // knot-orignating fields' names are Field_*

func (branch *GoBranch) AugmentsField() *GoField {
	return &GoField{
		Name: AugmentsFieldName,
		Type: branch.Augments,
		Tag:  KoTags("!augments", false), // user cannot access "!"
	}
}

func (branch *GoBranch) Real() GoType {
	return NewGoNeverNilPtr(
		NewGoStruct(
			append(
				[]*GoField{branch.AugmentsField()},
				branch.With...,
			)...,
		),
	)
}

func (branch *GoBranch) BranchRealField() *GoField { // field within real variety
	return &GoField{
		Name: branch.GoName(),
		Type: branch.Real(),
		Tag:  KoTags(branch.KoName(), false),
	}
}

func (branch *GoBranch) ProjectionLine() (line []*GoProjectionLine) {
	if _, ok := branch.Augments.(Unknown); ok {
		return []*GoProjectionLine{{
			Macro:     GoUnknownMacro{},
			Root:      branch,
			Leaf:      branch,
			PathField: ProjectField([]int{}, branch.With),
		}}
	}
	if macro := branch.Augments.VarietyMacro(); macro != nil {
		return []*GoProjectionLine{{
			Macro:     macro,
			Root:      branch,
			Leaf:      branch,
			PathField: ProjectField([]int{}, branch.With),
		}}
	} else {
		for i, child := range branch.Augments.VarietyBranch() {
			for _, partial := range child.ProjectionLine() {
				path := append(append([]int{}, partial.Path...), i)
				// path := append([]int{i}, partial.Path...)
				line = append(
					line,
					&GoProjectionLine{
						Path:  path,
						Macro: partial.Macro,
						Root:  branch,
						First: child, // root child branch on the path to line
						Leaf:  partial.Leaf,
						PathField: append(
							append([]*GoPathField{}, partial.PathField...),
							ProjectField(path, branch.With)...,
						),
					},
				)
			}
		}
		return line
	}
}

func BuildBranchProjectionIndex(projection []*GoProjectionLine) (index GoBranchProjectionLine) {
	index = GoBranchProjectionLine{}
	for _, p := range projection {
		index[p.Leaf] = p
	}
	return index
}

type GoProjectCtx struct {
	Variety       *GoVariety             `ko:"name=variety"`      // root
	Depth         int                    `ko:"name=depth"`        // depth from root
	SwitchBranch  GoExpr                 `ko:"name=switchBranch"` // expr to switch over
	SwitchVariety GoExpr                 `ko:"name=switchVariety"`
	Terminal      GoBranchProjectionLine `ko:"name=terminal"` // projection structures at leaves
}

type BranchExprRenderer interface {
	RenderBranchExpr(depth int) GoExpr
}

func (GoBranch) RenderBranchExpr(depth int) GoExpr {
	return &GoVerbatimExpr{fmt.Sprintf("b%d", depth)}
}

// ProjectExpr returns the go expression projecting a real (go) variety type to its real (go) projection type.
// Projection is determined by a root branch (captured within ctx) and the current branch (captured by branch).
func (branch *GoBranch) ProjectExpr(ctx *GoProjectCtx) GoExpr { // #ProjectExprRenderingBegins
	var case_ []*GoSwitchCaseExpr
	var default_ GoExpr = &GoPanicExpr{}
	overVariety := &GoVerbatimExpr{fmt.Sprintf("v%d", ctx.Depth)}
	overBranch := branch.RenderBranchExpr(ctx.Depth)
	if macro := branch.Augments.VarietyMacro(); macro != nil {
		// SEE: flatten variety to projection
		projection := ctx.Terminal[branch]
		default_ = &GoReturnExpr{
			&GoMakeStructExpr{ // make projection structure
				For: VarietalProjectionReal(ctx.Variety),
				Field: []*GoFieldExpr{{ // SEE: projection branch argument name and assignment
					Field: projection.ProjectionRealField(),
					Expr: &GoMakeStructExpr{
						For:   projection.Real(),
						Field: projection.ProjectionFieldExpr(ctx.Depth, branch),
					},
				}},
			},
		}
	} else { // #SwitchAcrossVarietyAncestors
		for _, subBranch := range branch.Augments.VarietyBranch() {
			branchSelection := &GoSelectExpr{Into: overVariety, Field: subBranch.BranchRealField().Name}
			varietySelection := &GoSelectExpr{Into: branchSelection, Field: AugmentsFieldName}
			case_ = append(
				case_,
				&GoSwitchCaseExpr{
					Predicate: &GoInequalityExpr{
						Left:  branchSelection,
						Right: NilExpr,
					},
					Expr: subBranch.ProjectExpr( // recurse
						&GoProjectCtx{
							Variety:       ctx.Variety,
							Depth:         ctx.Depth + 1,
							SwitchBranch:  branchSelection,
							SwitchVariety: varietySelection,
							Terminal:      ctx.Terminal,
						},
					),
				},
			)
		}
	}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoColonPairAssignExpr{
				Left:  [2]GoExpr{overBranch, overVariety},
				Right: [2]GoExpr{ctx.SwitchBranch, ctx.SwitchVariety},
			}, // #ProjectExprRenderingBegins
			&GoPairAssignExpr{
				Left:  [2]GoExpr{UnderlineExpr, UnderlineExpr},
				Right: [2]GoExpr{overBranch, overVariety},
			},
			&GoSwitchExpr{
				Over:    EmptyExpr,
				Case:    case_,
				Default: default_,
			}, // #SwitchAcrossVarietyAncestors
		},
	}
}

func (projection *GoProjectionLine) ProjectionFieldExpr(depth int, over BranchExprRenderer) (fieldExpr []*GoFieldExpr) {
	fieldExpr = make([]*GoFieldExpr, len(projection.PathField))
	for i, pathField := range projection.PathField {
		fieldExpr[i] = &GoFieldExpr{
			Field: pathField.Field,
			Expr: &GoSelectExpr{
				Into:  over.RenderBranchExpr(depth - pathField.Depth()),
				Field: pathField.Field.Name,
			},
		}
	}
	return
}
