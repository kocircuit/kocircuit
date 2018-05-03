package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func (b WeaveObject) Augment(weaveSpan *Span, fields Fields) (Shape, Effect, error) {
	ctx := b.Controller.WeaveStepCtx(weaveSpan)
	delegatedSpan := ctx.DelegateSpan()
	weaveFields, err := b.Controller.GroupWeaveFields(delegatedSpan, fields)
	if err != nil {
		return nil, nil, err
	}
	if residue, err := b.Controller.Weaver.Augment(ctx, b.Object, weaveFields); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returns), b.Controller.WrapEffect(residue.Effect), nil
	}
}

func (b *WeaveController) GroupWeaveFields(span *Span, fields Fields) (WeaveFields, error) {
	weaveFields := WeaveFields{}
	for _, fieldGroup := range fields.FieldGroup() {
		groupName := fieldGroup[0].Name
		fieldGroup = FilterEmptyWeaveFields(fieldGroup)
		switch len(fieldGroup) {
		case 0: // add a field with an empty symbol value (useful to All macro)
			weaveFields = append(weaveFields,
				&WeaveField{Name: groupName, Monadic: groupName == "", Objects: EmptySymbol{}},
			)
		case 1:
			y := fieldGroup[0].Shape.(WeaveObject)
			weaveFields = append(weaveFields,
				&WeaveField{Name: groupName, Monadic: groupName == "", Objects: y.Object},
			)
		default:
			repeatedSymbols := make(Symbols, len(fieldGroup))
			for i, f := range fieldGroup {
				repeatedSymbols[i] = f.Shape.(WeaveObject).Object
			}
			if repeated, err := MakeSeriesSymbol(span, repeatedSymbols); err != nil {
				return nil, span.Errorf(err, "weave repeated field group %s", groupName)
			} else {
				weaveFields = append(weaveFields,
					&WeaveField{Name: groupName, Monadic: groupName == "", Objects: repeated},
				)
			}
		}
	}
	return weaveFields, nil
}

func FilterEmptyWeaveFields(group []Field) (filtered []Field) {
	for _, field := range group {
		if !IsEmptySymbol(field.Shape.(WeaveObject).Object) {
			filtered = append(filtered, field)
		}
	}
	return
}
