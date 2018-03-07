package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (ctx *AssignCtx) assignVarietyVariety(from *GoVariety, to *GoVariety) (shaper Shaper, err error) {
	fromTree, toTree := VarietalVariationTree(from), VarietalVariationTree(to)
	fromVtns := fromTree.Range()
	vtnShapers := make([]*VariationShaper, len(fromVtns))
	for i, fromVtn := range fromVtns { // for each from-variation
		if toVtnWalk := toTree.FindWalk(fromVtn.Macro); len(toVtnWalk) == 0 { // find corresponding to-variation
			return nil, ctx.Errorf(
				ctx.Errorf(nil, "no to-variation for %s", Sprint(fromVtn.Macro)),
				"assigning variety %s to %s", Sprint(from), Sprint(to),
			)
		} else {
			if fieldShapers, err := ctx.assignInsert(fromVtn.PathField, toVtnWalk.Leaf().Variation.PathField); err != nil {
				return nil, ctx.Errorf(err, "assigning variety %s to %s", Sprint(from), Sprint(to))
			} else {
				vtnShapers[i] = &VariationShaper{
					Macro:  fromVtn.Macro,
					WalkTo: toVtnWalk, // tree path to the to-variant
					Field:  fieldShapers.VariationFieldShapers(),
				}
			}
		}
	}
	return &VarietyShaper{
		Shaping:   Shaping{Origin: ctx.Origin, From: from, To: to},
		FromTree:  fromTree,
		Variation: vtnShapers, // XXX: vtnShapers contains a nil shaper
	}, nil
}

type VarietyShaper struct {
	Shaping   `ko:"name=shaping"`
	FromTree  *GoVariationTree   `ko:"name=from"`      // (rendering) traverse this to switch on non-nil from-variation (tree leaf)
	Variation []*VariationShaper `ko:"name=variation"` // (rendering) describes to-assignments, for a non-nil from-variation
}

func (sh *VarietyShaper) FindVariationShaper(vtn *GoVariation) *VariationShaper {
	for _, vs := range sh.Variation {
		if vs.Macro.MacroID() == vtn.Macro.MacroID() {
			return vs
		}
	}
	panic("o")
}

type VariationShaper struct {
	Macro  Macro                   `ko:"name=macro"`
	WalkTo GoVariationWalk         `ko:"name=pathTo"` // tree nodes on way to to-variation
	Field  []*VariationFieldShaper `ko:"name=field"`
}

type VariationFieldShaper struct {
	From   *GoPathField `ko:"name=from"`   // location of from-field
	To     *GoPathField `ko:"name=to"`     // location of to-field
	Shaper Shaper       `ko:"name=shaper"` // field shaper
}

func (sh *VarietyShaper) ShaperID() string {
	return Mix("VarietyShaper", sh.ShapingID())
}

func (sh *VarietyShaper) String() string {
	return Sprint(sh)
}

func (sh *VarietyShaper) Reverse() Shaper {
	panic("irreversible")
}

func (sh *VarietyShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *VarietyShaper) CircuitEffect() *GoCircuitEffect {
	effect := &GoCircuitEffect{DuctFunc: []GoFuncExpr{sh}} // VarietyShaper is a duct function
	for _, vtn := range sh.Variation {
		if vtn == nil {
			println("VTNNIL")
		}
		for _, field := range vtn.Field {
			effect = effect.Aggregate(field.Shaper.CircuitEffect())
		}
	}
	return effect
}

func (sh *VarietyShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("variety_shaper_%s", sh.ShapingID()),
	}
}

func (sh *VarietyShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *VarietyShaper) funcExpr() GoFuncExpr {
	fromExpr := &GoVerbatimExpr{"from"}
	return &GoShaperFuncExpr{
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line: []GoExpr{
			(&varietySwitchCtx{Shaper: sh}).TreeExpr(sh.FromTree, nil, fromExpr),
			&GoPanicExpr{}, // panic("o")
		},
	}
}

type varietySwitchCtx struct {
	Shaper *VarietyShaper `ko:"name=shaper"`
}

func (ctx *varietySwitchCtx) TreeExpr(from *GoVariationTree, branchExpr, vtyExpr GoExpr) GoExpr {
	if from.Leaf != nil {
		return ctx.FillExpr(from, branchExpr, vtyExpr)
	} else {
		return ctx.SwitchExpr(from, branchExpr, vtyExpr)
	}
}

// 	v0 := █
// 	switch {
// 	case v0.Branch1 != nil:
// 		b1 := v0.Branch1
// 		v1 := b1.Augments
// 		...
// 	case v0.Branch2 != nil:
// 		b1 := v0.Branch2
// 		v1 := b1.Augments
// 		...
// 	default:
// 		panic("o")
// 	}
func (ctx *varietySwitchCtx) SwitchExpr(from *GoVariationTree, branchExpr, vtyExpr GoExpr) GoExpr {
	vtyName, initExpr := ctx.InitExpr(from, branchExpr, vtyExpr)
	cases := make([]*GoSwitchCaseExpr, len(from.Branch)) // switch cases
	for i, branch := range from.Branch {
		nextBranchExpr := &GoSelectExpr{Into: vtyName, Field: branch.Branch.GoName()} // v0.Branch1
		cases[i] = &GoSwitchCaseExpr{
			Predicate: &GoInequalityExpr{nextBranchExpr, NilExpr}, // v0.Branch1 != nil
			Expr:      ctx.TreeExpr(branch.Tree, nextBranchExpr, nil),
		}
	}
	return GoBlock(
		initExpr,
		&GoSwitchExpr{
			Over:    EmptyExpr,
			Case:    cases,
			Default: &GoPanicExpr{},
		},
	)
}

func (ctx *varietySwitchCtx) InitExpr(from *GoVariationTree, branchExpr, vtyExpr GoExpr) (vtyName GoExpr, initExpr GoExpr) {
	branchName := &GoVerbatimExpr{fmt.Sprintf("b%d", from.Depth)} // b0
	vtyName = &GoVerbatimExpr{fmt.Sprintf("v%d", from.Depth)}     // v0
	var renameBranchExpr, renameVtyExpr GoExpr
	switch {
	case branchExpr != nil && vtyExpr == nil:
		renameBranchExpr = &GoColonAssignExpr{Left: branchName, Right: branchExpr} // b0 := █
		renameVtyExpr = &GoColonAssignExpr{                                        // v0 := b0.Augments
			Left:  vtyName,
			Right: &GoSelectExpr{Into: branchName, Field: AugmentsFieldName},
		}
	case branchExpr == nil && vtyExpr != nil:
		if from.Depth != 0 { // only at depth 0, branch is not set
			panic("o")
		}
		renameBranchExpr = nil
		renameVtyExpr = &GoColonAssignExpr{Left: vtyName, Right: vtyExpr} // v0 := █
	default:
		panic("o")
	}
	return vtyName, GoBlock(renameBranchExpr, renameVtyExpr)
}

func (ctx *varietySwitchCtx) FillExpr(from *GoVariationTree, branchExpr, vtyExpr GoExpr) GoExpr {
	_, initExpr := ctx.InitExpr(from, branchExpr, vtyExpr)
	return GoBlock(initExpr, ctx.insertExpr(from))
}

func (ctx *varietySwitchCtx) insertExpr(fromLeaf *GoVariationTree) GoExpr {
	vs := ctx.Shaper.FindVariationShaper(fromLeaf.Leaf)
	// 	w0 := &W0{}
	// 	c1 := &C1{}
	// 	...
	// 	ck := &Ck{}
	// 	wk := &Wk{}
	line := []GoExpr{}
	walkTo := vs.WalkTo
	for i, step := range walkTo {
		line = append(line,
			&GoColonAssignExpr{ // w0 := &W0{}
				Left:  GoVerbatimf("w%d", i),             // w0
				Right: &GoZeroExpr{step.Varietal.Real()}, // &W0{}
			},
		)
		switch {
		case step.Branch != nil:
			line = append(line,
				&GoColonAssignExpr{ // c1 := &C1{}
					Left:  GoVerbatimf("c%d", i+1),         // c1
					Right: &GoZeroExpr{step.Branch.Real()}, // &C1{}
				},
			)
		case step.Variation != nil: // nop
		default:
			panic("o")
		}
	}
	// w0.BranchK0 = c1
	// c1.Augments = w1
	// w1.BranchK1 = c2
	// ...
	// cN.Augments = wN
	for i, step := range walkTo {
		switch {
		case step.Branch != nil:
			line = append(line,
				&GoAssignExpr{ // w0.BranchK = c1
					Left:  &GoSelectExpr{Into: GoVerbatimf("w%d", i), Field: step.Branch.GoName()}, // w0.BranchK
					Right: GoVerbatimf("c%d", i+1),                                                 // c1
				},
			)
		case step.Variation != nil: // nop
		default:
			panic("o")
		}
	}
	// c1.FieldN = Shaper_c1FieldN(bK.FieldM)
	// ...
	for _, vfs := range vs.Field {
		line = append(line,
			&GoAssignExpr{
				Left: &GoSelectExpr{
					Into:  GoVerbatimf("c%d", vfs.To.Depth()+1),
					Field: vfs.To.Field.Name,
				}, // c1.FieldN
				Right: &GoShapeExpr{ // Shaper_c1FieldN(bK.FieldM)
					Shaper: vfs.Shaper,
					Expr: &GoSelectExpr{
						Into:  GoVerbatimf("b%d", vfs.From.Depth()+1),
						Field: vfs.From.Field.Name,
					},
				},
			},
		)
	}
	return GoBlock(line...)
}
