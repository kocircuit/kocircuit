package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func GoVarietyUnion(span *Span, x, y *GoVariety) (_ *GoVariety, err error) {
	// build variations for each leaf in x and y, and
	xVtns := VarietalVariationTree(x).Range()
	yVtns := VarietalVariationTree(y).Range()
	// unify variations for duplicate macros
	vMat := (GoVariationMatrix{}).Merge(xVtns).Merge(yVtns)
	vtns := make(GoVariations, len(vMat))
	for i, corrVtns := range vMat {
		switch len(corrVtns) {
		case 1:
			vtns[i] = corrVtns[0]
		case 2:
			if unifiedArg, err := Generalize(span, corrVtns[0].Arg(), corrVtns[1].Arg()); err != nil {
				return nil, span.Errorf(err, "merging augmentations")
			} else {
				vtns[i] = &GoVariation{
					Macro:     vtns[0].Macro,
					PathField: GoFieldsOnPath(unifiedArg.(*GoStruct).Field, nil),
				}
			}
		default:
			panic("o")
		}
	}
	return vtns.GoVariety(span), nil
}

type GoVariationMatrix []GoVariations

func (in GoVariationMatrix) Copy() (copied GoVariationMatrix) {
	copied = make(GoVariationMatrix, len(in))
	for i := range in {
		copied[i] = in[i].Copy()
	}
	return copied
}

func (in GoVariationMatrix) Merge(vtns GoVariations) (merged GoVariationMatrix) {
	merged = in.Copy()
	for _, vtn := range vtns {
		if vtn == nil {
			panic("XXX") // not met
		}
		if index, found := merged.Find(vtn.Macro); found {
			merged[index] = append(merged[index], vtn)
		} else {
			merged = append(merged, GoVariations{vtn})
		}
	}
	return
}

func (vmat GoVariationMatrix) Find(macro Macro) (int, bool) {
	id := macro.MacroID()
	for i, vtns := range vmat {
		if vtns[0].Macro.MacroID() == id {
			return i, true
		}
	}
	return -1, false
}

func VarietalVariationTree(v GoVarietal) (tree *GoVariationTree) {
	ctx := &variationTreeCtx{}
	tree = ctx.variationTree(v)
	tree.verifyNoDuplicateMacros()
	return tree
}

type variationTreeCtx struct {
	PathFromRoot []int          `ko:"name=pathFromRoot"`
	PathField    []*GoPathField `ko:"name=pathField"`
}

func (ctx *variationTreeCtx) Depth() int {
	return len(ctx.PathFromRoot)
}

func (ctx *variationTreeCtx) PathToBranch(index int) []int {
	return append(append([]int{}, ctx.PathFromRoot...), index)
}

func (ctx *variationTreeCtx) Push(index int, pathField []*GoPathField) *variationTreeCtx {
	return &variationTreeCtx{
		PathFromRoot: append(append([]int{}, ctx.PathFromRoot...), index),
		PathField:    append(append([]*GoPathField{}, ctx.PathField...), pathField...),
	}
}

func (ctx *variationTreeCtx) variationTree(v GoVarietal) (tree *GoVariationTree) {
	if v.VarietyMacro() != nil {
		return &GoVariationTree{
			Varietal: v, Depth: ctx.Depth(),
			Leaf: &GoVariation{PathField: ctx.PathField, Macro: v.VarietyMacro()},
		}
	} else {
		variationBranches := make([]*GoVariationTreeBranch, len(v.VarietyBranch()))
		for i, branch := range v.VarietyBranch() {
			pathToBranch := ctx.PathToBranch(branch.Index)
			branchFields := make([]*GoPathField, len(branch.With))
			for j, withField := range branch.With {
				branchFields[j] = &GoPathField{Path: pathToBranch, Field: withField}
			}
			variationBranches[i] = &GoVariationTreeBranch{
				Tree:   ctx.Push(branch.Index, branchFields).variationTree(branch.Augments),
				Branch: branch,
			}
		}
		return &GoVariationTree{
			Varietal: v, Depth: ctx.Depth(),
			Branch: variationBranches,
		}
	}
}

type GoVariationTree struct {
	Varietal GoVarietal               `ko:"name=varietal"`
	Depth    int                      `ko:"name=depth"`  // depth from root
	Branch   []*GoVariationTreeBranch `ko:"name=branch"` // either branch
	Leaf     *GoVariation             `ko:"name=leaf"`   // or leaf
}

type GoVariationTreeBranch struct {
	Tree   *GoVariationTree `ko:"name=tree"`
	Branch *GoBranch        `ko:"name=branch"`
}

type GoVariationWalk []*GoVariationStep

func (vw GoVariationWalk) Leaf() *GoVariationStep {
	return vw[len(vw)-1]
}

type GoVariationStep struct {
	Varietal  GoVarietal   `ko:"name=varietal"`
	Branch    *GoBranch    `ko:"name=branch"`    // either branch
	Variation *GoVariation `ko:"name=variation"` // or variation
}

// FindWalk returns the set of nodes on the path to the variant leaf
func (tree *GoVariationTree) FindWalk(macro Macro) GoVariationWalk {
	if tree.Leaf != nil {
		if tree.Leaf.Macro.MacroID() == macro.MacroID() {
			return GoVariationWalk{{Varietal: tree.Varietal, Variation: tree.Leaf}}
		}
	} else {
		for _, branch := range tree.Branch {
			if p := branch.Tree.FindWalk(macro); p != nil {
				return append(GoVariationWalk{{Varietal: tree.Varietal, Branch: branch.Branch}}, p...)
			}
		}
	}
	return nil
}

func (tree *GoVariationTree) verifyNoDuplicateMacros() {
	dup := map[string]bool{}
	for _, m := range tree.rangeMacros() {
		if dup[m.MacroID()] {
			panic("o")
		} else {
			dup[m.MacroID()] = true
		}
	}
}

func (tree *GoVariationTree) rangeMacros() (macros []Macro) {
	if tree.Leaf != nil {
		return []Macro{tree.Leaf.Macro}
	} else {
		for _, branch := range tree.Branch {
			macros = append(macros, branch.Tree.rangeMacros()...)
		}
		return
	}
}

func (tree *GoVariationTree) Range() (leaves GoVariations) {
	if tree.Leaf != nil {
		return GoVariations{tree.Leaf}
	} else {
		for _, branch := range tree.Branch {
			leaves = append(leaves, branch.Tree.Range()...)
		}
		return
	}
}

type GoVariations []*GoVariation

func (vtns GoVariations) Copy() (copied GoVariations) {
	copied = make(GoVariations, len(vtns))
	for i := range vtns {
		copied[i] = vtns[i]
	}
	return copied
}

func (vtns GoVariations) Seen(macro Macro) bool {
	for _, vtn := range vtns {
		if vtn.Macro.MacroID() == macro.MacroID() {
			return true
		}
	}
	return false
}

func (vtns GoVariations) GoVariety(span *Span) *GoVariety {
	branches := make([]*GoBranch, len(vtns))
	for i, vtn := range vtns {
		branches[i] = &GoBranch{
			Index:    i,
			Augments: NewGoVariety(span, vtn.Macro, nil),
			With:     vtn.Field(),
		}
	}
	return NewGoVariety(span, nil, branches)
}

type GoVariation struct {
	PathField []*GoPathField `ko:"name=field"` // field (with corresponding path) list on path to leaf variety
	Macro     Macro          `ko:"name=macro"` // terminating macro
}

// Field returns augmented fields along the path to the leaf variety.
func (vtn *GoVariation) Field() []*GoField {
	field := make([]*GoField, len(vtn.PathField))
	for i, pathField := range vtn.PathField {
		field[i] = pathField.Field
	}
	return field
}

func (vtn *GoVariation) Arg() *GoStruct {
	return NewGoStruct(vtn.Field()...)
}
