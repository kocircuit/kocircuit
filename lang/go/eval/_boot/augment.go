package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b BootObject) Augment(bootSpan *Span, fields Fields) (Shape, Effect, error) {
	ctx := b.Controller.BootStepCtx(bootSpan)
	delegatedSpan := ctx.DelegateSpan()
	bootFields, err := b.Controller.GroupBootFields(delegatedSpan, fields)
	if err != nil {
		return nil, nil, err
	}
	if residue, err := b.Booter.Augment(ctx, b.Object, bootFields); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.WrapEffect(residue.Effect), nil
	}
}

func (b *BootController) GroupBootFields(span *Span, fields Fields) (BootFields, error) {
	bootFields := BootFields{}
	for _, fieldGroup := range fields.FieldGroup() {
		groupName := fieldGroup[0].Name
		fieldGroup = FilterEmptyBootFields(fieldGroup)
		switch len(fieldGroup) {
		case 0: // add a field with an empty symbol value (useful to All macro)
			bootFields = append(bootFields,
				&BootField{Name: groupName, Monadic: groupName == "", Objects: EmptySymbol{}},
			)
		case 1:
			y := fieldGroup[0].Shape.(BootSymbol)
			bootFields = append(bootFields,
				&BootField{Name: groupName, Monadic: groupName == "", Objects: y.Object},
			)
		default:
			repeatedSymbols := make(Symbols, len(fieldGroup))
			for i, f := range fieldGroup {
				repeatedSymbols[i] = f.Shape.(BootSymbol).Object
			}
			if repeated, err := MakeSeriesSymbol(span, repeatedSymbols); err != nil {
				return nil, span.Errorf(err, "boot repeated field group %s", groupName)
			} else {
				bootFields = append(bootFields,
					&FieldSymbol{Name: groupName, Monadic: groupName == "", Objects: repeated},
				)
			}
		}
	}
	return bootFields, nil
}

func FilterEmptyBootFields(group []Field) (filtered []Field) {
	for _, field := range group {
		if !IsEmptySymbol(field.Shape.(BootSymbol).Object) {
			filtered = append(filtered, field)
		}
	}
	return
}
