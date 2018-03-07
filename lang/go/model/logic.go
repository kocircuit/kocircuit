// Package model provides a language model for Go implementations of Ko programs.
package model

// GoStepLogic is responsible for rendering the synchronous, core of the step.
type GoStepLogic interface {
	// Arguments of step logic are the received input values after shaping.
	// Arguments are accessible via GoArrival.SlotExpr.
	GoSlotForm
}

// GoEnterLogic returns the argument structure of the enclosing circuit function.
type GoEnterLogic struct{}

func (enter *GoEnterLogic) FormExpr(...*GoSlotExpr) GoExpr {
	return &GoVerbatimExpr{`arg`}
}

// GoLeaveLogic returns its default slot argument.
type GoLeaveLogic struct{}

func (leave *GoLeaveLogic) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return FindSlotExpr(arg, RootSlot{})
}

// GoAugmentLogic extends a variety structure.
type GoAugmentLogic struct {
	From GoVarietal        `ko:"name=from"`
	To   GoVarietal        `ko:"name=to"`
	With []*GoAugmentField `ko:"name=with"`
}

type GoAugmentField struct {
	Field *GoField `ko:"name=field"`
	Slot  []Slot   `ko:"name=slot"`
}

func StripAugmentField(augField []*GoAugmentField) []*GoField {
	f := make([]*GoField, len(augField))
	for i, augField := range augField {
		f[i] = augField.Field
	}
	return f
}

// The empty-string field (in receive) carries the base variety.
// The remaining (named) fields carry the augmentations.
func (augment *GoAugmentLogic) FormExpr(arg ...*GoSlotExpr) GoExpr {
	// toVty has a singleton branch
	toVty := augment.To.(*GoVariety)
	branchFieldExpr := []*GoFieldExpr{{
		Field: toVty.Branch[0].AugmentsField(),
		Expr:  FindSlotExpr(arg, RootSlot{}),
	}}
	// for _, field := range toVty.Branch[0].With {
	for _, with := range augment.With {
		branchFieldExpr = append(branchFieldExpr, with.Expr(arg))
	}
	expr := &GoMakeStructExpr{ // make variety
		For: augment.To.Real(), // variety
		Field: []*GoFieldExpr{{ // variety branch assignment expression
			Field: toVty.Branch[0].BranchRealField(), // variety's branch field
			Expr: &GoMakeStructExpr{ // make branch
				For:   toVty.Branch[0].Real(), // branch
				Field: branchFieldExpr,        // augment + with field assignment expressions
			}},
		},
	}
	return expr
}

func (augField *GoAugmentField) Expr(arg []*GoSlotExpr) *GoFieldExpr {
	switch len(augField.Slot) {
	case 0:
		panic("o")
	case 1:
		return &GoFieldExpr{
			Field: augField.Field,
			Expr:  FindSlotExpr(arg, augField.Slot[0]),
		}
	default:
		elemExpr := make([]GoExpr, len(augField.Slot))
		for i, augSlot := range augField.Slot {
			elemExpr[i] = FindSlotExpr(arg, augSlot)
		}
		return &GoFieldExpr{
			Field: augField.Field,
			Expr: &GoMakeSequenceExpr{
				Type: augField.Field.Type,
				Elem: elemExpr,
			},
		}
	}
}

type GoInvokeLogic struct {
	DuctFunc  *GoDuctFuncExpr `ko:"name=ductFunc"`
	Projector Shaper          `ko:"name=projector"` // projects variety to a variety projection
}

func (invoke *GoInvokeLogic) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallDuctFuncExpr{
		Func: invoke.DuctFunc,
		Arg: &GoShapeExpr{
			Shaper: invoke.Projector,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		},
	}
}
