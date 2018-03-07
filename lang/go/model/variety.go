package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Varietal interface {
	// macro
	VarietyMacro() Macro
	VarietyExpr() GoExpr // corresponds to macro
	// augmentation
	VarietyBranch() []*GoBranch
	// carrier
	VarietyAddress() *GoAddress
	VarietyProjectionAddress() *GoAddress
}

type GoReal interface {
	Real() GoType
}

// GoVariety and (solved) GoUnknown are both GoVarietals.
type GoVarietal interface {
	GoType
	GoReal
	Varietal
}

// GoVariety is realized as GoNeverNilPtr{GoAlias{GoStruct}}
type GoVariety struct {
	Origin             *Span       `ko:"name=origin"`
	Macro              Macro       `ko:"name=macro"`  // either macro
	Branch             []*GoBranch `ko:"name=branch"` // or augmentation branch
	Alias_             *GoAlias    `ko:"name=alias"`
	AddressVariety_    *GoAddress  `ko:"name=addressVariety"`
	AddressProjection_ *GoAddress  `ko:"name=addressProjection"`
}

func NewGoVariety(span *Span, macro Macro, branch []*GoBranch) *GoVariety {
	real := varietyReal(branch)
	addressVariety := NewSpanAddress(span, "variety", real)
	addressProjection := NewSpanAddress(span, "projection", real)
	return &GoVariety{
		Origin:             span,
		Macro:              macro,
		Branch:             branch,
		Alias_:             NewGoAlias(addressVariety, real),
		AddressVariety_:    addressVariety,
		AddressProjection_: addressProjection,
	}
}

func (vty *GoVariety) TypeID() string { return vty.Real().TypeID() }

func (vty *GoVariety) VarietyAddress() *GoAddress { return vty.AddressVariety_ }

func (vty *GoVariety) VarietyProjectionAddress() *GoAddress { return vty.AddressProjection_ }

func (vty *GoVariety) Real() GoType {
	return NewGoNeverNilPtr(vty.Alias_)
}

func varietyReal(branch []*GoBranch) GoType {
	field := make([]*GoField, len(branch))
	for i, branch := range branch {
		field[i] = branch.BranchRealField()
	}
	return NewGoStruct(field...)
}

type variety struct {
	Origin *Span       `ko:"name=origin"`
	Real   string      `ko:"name=real"`   // real go type name
	Macro  Macro       `ko:"name=macro"`  // either macro
	Branch []*GoBranch `ko:"name=branch"` // or augmentation branch
}

func (vty variety) Relabel() Label {
	return Label{Name: "Variety"}
}

func (vty *GoVariety) Splay() Tree {
	v := variety{
		Origin: vty.Origin,
		Real:   vty.VarietyAddress().String(),
		Macro:  vty.Macro,
		Branch: vty.Branch,
	}
	return Splay(v)
}

func (vty *GoVariety) VarietyMacro() Macro { return vty.Macro }

func (vty *GoVariety) VarietyBranch() []*GoBranch { return vty.Branch }

func (vty *GoVariety) Doc() string { return Sprint(vty) }

func (vty *GoVariety) String() string { return Sprint(vty) }

func (vty *GoVariety) Sketch(ctx *GoSketchCtx) interface{} {
	return vty.Real().Sketch(ctx)
}

func (vty *GoVariety) Tag() []*GoTag { return nil }

func (vty *GoVariety) RenderDef(fileCtx GoFileContext) string {
	return vty.Real().RenderDef(fileCtx)
}

func (vty *GoVariety) RenderRef(fileCtx GoFileContext) string {
	return vty.Real().RenderRef(fileCtx)
}

func (vty *GoVariety) RenderZero(fileCtx GoFileContext) string {
	return fmt.Sprintf("&%s{}", vty.Alias_.RenderRef(fileCtx))
}

// RenderExpr implements GoExpr.
func (vty *GoVariety) VarietyExpr() GoExpr {
	return &GoZeroExpr{vty} // expression that renders as the zero value of vty
}

type WeaveCtx interface {
	Valve() *GoValve
}

// VarietalProject returns a shaper from variety to its flattened struct representation.
func VarietalProject(span *Span, varietal GoVarietal) (GoStructure, Shaper) {
	projected := VarietalProjectionReal(varietal)
	return projected, &ProjectShaper{
		Shaping: Shaping{Origin: span, From: varietal, To: projected},
	}
}

func VarietalProjectionReal(varietal GoVarietal) GoStructure {
	return NewGoNeverNilPtr(
		NewGoAlias(
			varietal.VarietyProjectionAddress(),
			NewGoStruct(VarietalProjectionLineFields(varietal)...),
		),
	)
}

func VarietalProjectionLine(varietal GoVarietal) []*GoProjectionLine {
	return NewRootBranch(varietal).ProjectionLine()
}

func VarietalProjectionLineFields(varietal GoVarietal) []*GoField {
	line := VarietalProjectionLine(varietal)
	field := make([]*GoField, len(line))
	for i, line := range line {
		field[i] = line.ProjectionRealField()
	}
	return field
}

func (vty *GoVariety) Augment(span *Span, with ...*GoField) (augmented *GoVariety) {
	return AugmentVarietalWithField(vty, span, with...)
}

func AugmentVarietalWithField(varietal GoVarietal, span *Span, with ...*GoField) (augmented *GoVariety) {
	return NewGoVariety(
		span,
		nil,
		[]*GoBranch{{
			Index:    0,
			Augments: varietal,
			With:     with,
		}},
	)
}

// verify no duplicate augmentations

func VerifyUniqueFieldAugmentations(span *Span, varietal GoVarietal) (FieldAugmentingGoVariety, error) {
	switch v := varietal.(type) {
	case Unknown:
		return nil, nil
	case *GoVariety:
		return v.VerifyUniqueFieldAugmentations(span)
	}
	panic("o")
}

func (vty *GoVariety) VerifyUniqueFieldAugmentations(span *Span) (FieldAugmentingGoVariety, error) {
	r := FieldAugmentingGoVariety{}
	for _, branch := range vty.Branch {
		if b, err := branch.VerifyUniqueFieldAugmentations(span, vty); err != nil {
			return nil, err
		} else {
			r = r.Merge(b)
		}
	}
	return r, nil
}

func (branch *GoBranch) VerifyUniqueFieldAugmentations(span *Span, parent *GoVariety) (FieldAugmentingGoVariety, error) {
	augments, err := VerifyUniqueFieldAugmentations(span, branch.Augments)
	if err != nil {
		return nil, err
	}
	for _, field := range branch.With {
		if len(augments[field.KoName()]) > 0 {
			return nil, span.Errorf(
				nil,
				"augmenting field %s at (%s), but field already augmented at %s",
				field.KoName(),
				parent.Origin.SourceLine(),
				augments[field.KoName()][0].Origin.SourceLine(),
			)
		} else {
			augments[field.KoName()] = []*GoVariety{parent}
		}
	}
	return augments, nil
}

type FieldAugmentingGoVariety map[string][]*GoVariety // field name -> descendant varieties that augment it

func (x FieldAugmentingGoVariety) Merge(y FieldAugmentingGoVariety) FieldAugmentingGoVariety {
	p := FieldAugmentingGoVariety{}
	for field, vties := range x {
		p[field] = vties
	}
	for field, vties := range y {
		p[field] = append(p[field], vties...)
	}
	return p
}
