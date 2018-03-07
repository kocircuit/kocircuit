package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (local *GoLocal) Augment(span *Span, knot Knot) (Shape, Effect, error) {
	vty, ok := local.Image().(GoVarietal) // GoVariety or Unknown
	if !ok {
		return nil, nil, span.Errorf(nil, "augment expects a variety, got %s", Sprint(local.Image()))
	}
	span = AugmentSpanCache(span, local.Cached) // merge span and upstream caches
	knotSlots, err := KnotToSlots(span, knot)
	if err != nil {
		return nil, nil, err
	}
	var aug *Augmentation
	if aug, err = GoAugment(span, vty, knotSlots.With); err != nil {
		return nil, nil, err
	}
	step := &GoStep{
		Span:  span,
		Label: NearestStep(span).Label,
		Arrival: append(
			[]*GoArrival{
				ArrivalFromLocal(span, local, RootSlot{}),
			},
			knotSlots.Arrival...,
		),
		Returns: aug.Varietal,
		Logic:   aug.Logic,
		Cached:  knotSlots.Cached,
	}
	return local.Inherit(span, step, aug.Varietal),
		&GoStepEffect{
			Step:          step,
			CircuitEffect: aug.CircuitEffect().AggregateDuctType(knotSlots.Duct...),
		}, nil
}

type Augmentation struct {
	Origin   *Span           `ko:"name=origin"`
	Varietal GoVarietal      `ko:"name=varietal"`
	Logic    *GoAugmentLogic `ko:"name=logic"`
}

func (aug *Augmentation) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return aug.Logic.FormExpr(arg...)
}

func (aug *Augmentation) CircuitEffect() *GoCircuitEffect {
	if aug == nil {
		return nil
	}
	return &GoCircuitEffect{
		DuctType: []GoType{aug.Varietal},
	}
}

func (aug *Augmentation) ProgramEffect() *GoProgramEffect {
	return nil
}

func WithGoField(field ...*GoField) []*GoAugmentField {
	augField := make([]*GoAugmentField, len(field))
	for i, f := range field {
		augField[i] = &GoAugmentField{
			Field: f,
			Slot: []Slot{
				NameSlot{f.KoName()},
			},
		}
	}
	return augField
}

func GoAugment(span *Span, vty GoVarietal, with []*GoAugmentField) (aug *Augmentation, err error) {
	aug = &Augmentation{Origin: span}
	aug.Varietal = AugmentVarietalWithField(vty, span, StripAugmentField(with)...)
	aug.Logic = &GoAugmentLogic{From: vty, To: aug.Varietal, With: with}
	if _, err = VerifyUniqueFieldAugmentations(span, aug.Varietal); err != nil {
		return nil, err
	}
	return aug, nil
}

type KnotSlots struct {
	Span    *Span             `ko:"name=span"`
	Arrival []*GoArrival      `ko:"name=arrival"`
	With    []*GoAugmentField `ko:"name=with"`
	Duct    []GoType          `ko:"name=duct"`
	Cached  *AssignCache      `ko:"name=cached"`
}

// Knot (a: X, a: Y, b: Z) becomes
//	&branch{█}{
//		a: []GCT(X, Y){
//			shape_x(slot_field_a0),	// X->slot_field_a0
//			shape_y(slot_field_a1),	// Y->slot_field_a1
//		},
//		b: shape_z(slot_field_b0),		// Z->slot_field_b0
//	}
// Knot (X, Y, Z) becomes
//	&branch{█}{
//		Etc: []GCT(X, Y, Z){
//			shape_x(slot_etc_0),	// X->slot_etc_0
//			shape_y(slot_etc_1),	// Y->slot_etc_1
//			shape_z(slot_etc_2),	// Z->slot_etc_2
//		},
//	}
func KnotToSlots(span *Span, knot Knot) (knotSlots *KnotSlots, err error) {
	knotSlots = &KnotSlots{Span: span}
	assn := NewAssignCtx(span)
	for _, fieldGroup := range knot.FieldGroup() { // fieldGroup []Field
		fieldGroup = RemoveFieldGroupGoEmpties(span, fieldGroup)
		var groupType GoType
		var groupArrival []*GoArrival
		switch len(fieldGroup) {
		case 0:
			// leaves groupType == nil and groupArrival == nil
		case 1:
			local := fieldGroup[0].Shape.(*GoLocal)
			groupType = local.Image()
			groupArrival = []*GoArrival{
				ArrivalFromLocal(span, local, KnotFieldSlot{fieldGroup[0], 0}),
			}
		default:
			var elem GoType
			if elem, err = GeneralizeKnotFieldTypes(span, fieldGroup); err != nil {
				return nil, err
			} else {
				elem = LiftNumber(elem)
			}
			knotSlots.Duct = append(knotSlots.Duct, elem)
			groupType = NewGoSlice(elem) // arrays are not appropriate here, as they dont generalize with singleton types
			for i, field := range fieldGroup {
				local := field.Shape.(*GoLocal)
				bridge, err := assn.Assign(local.Image(), elem)
				if err != nil {
					panic("o") // GeneralizeKnotFieldTypes above guarantees success here
				}
				groupArrival = append(
					groupArrival,
					ArrivalFromLocal(span, local.Extend(span, bridge), KnotFieldSlot{field, i}),
				)
			}
		}
		if groupType != nil {
			knotSlots.Arrival = append(knotSlots.Arrival, groupArrival...)
			knotSlots.With = append(
				knotSlots.With,
				&GoAugmentField{
					Field: BuildGoField(span, fieldGroup[0].Name, groupType, fieldGroup[0].Name == NoLabel),
					Slot:  ArrivalSlots(groupArrival),
				},
			)
		}
	}
	knotSlots.Cached = assn.Flush()
	return
}

func RemoveFieldGroupGoEmpties(span *Span, fg []Field) (dense []Field) {
	for _, field := range fg {
		simple, _ := Simplify(span, field.Shape.(*GoLocal).Image())
		if _, isEmpty := simple.(*GoEmpty); !isEmpty {
			dense = append(dense, field)
		}
	}
	return
}

func LiftNumber(t GoType) GoType {
	switch n := t.(type) {
	case GoNumber:
		return n.Builtin()
	default:
		return n
	}
}

func ArrivalSlots(arrival []*GoArrival) []Slot {
	s := make([]Slot, len(arrival))
	for i, arrival := range arrival {
		s[i] = arrival.Slot
	}
	return s
}

func GeneralizeKnotFieldTypes(span *Span, field []Field) (generalized GoType, err error) {
	generalized = field[0].Shape.(*GoLocal).Image()
	for _, field := range field {
		if generalized, err = Generalize(span, generalized, field.Shape.(*GoLocal).Image()); err != nil {
			return nil, err
		}
	}
	return generalized, nil
}

func FieldGroupFields(group []Field) (field []*GoField) {
	field = make([]*GoField, len(group))
	for i, f := range group {
		local := f.Shape.(*GoLocal)
		field[i] = BuildGoField(nil, f.Name, local.Image(), false)
	}
	return
}
