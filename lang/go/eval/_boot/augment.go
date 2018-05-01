package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b BootObject) Augment(bootSpan *Span, fields Fields) (Shape, Effect, error) {
	span := b.Controller.BootStepCtx(bootSpan)
	XXX
}

func GroupBootFields(span *Span, fields Fields) (BootFields, error) {
	XXX
	ef := FieldSymbols{}
	for _, fieldGroup := range fields.FieldGroup() {
		fieldGroupName := fieldGroup[0].Name
		fieldGroup = FilterEmptyFields(fieldGroup)
		switch len(fieldGroup) {
		case 0: // add a field with an empty symbol value (useful to All macro)
			ef = append(ef,
				&FieldSymbol{
					Name:    fieldGroupName,
					Monadic: fieldGroupName == "",
					Value:   EmptySymbol{},
				},
			)
		case 1:
			y := fieldGroup[0].Shape.(Symbol)
			ef = append(ef,
				&FieldSymbol{
					Name:    fieldGroupName,
					Monadic: fieldGroupName == "",
					Value:   y,
				},
			)
		default:
			fieldTypes, fieldSymbols := make([]Type, len(fieldGroup)), make(Symbols, len(fieldGroup))
			for i, f := range fieldGroup {
				y := f.Shape.(Symbol)
				fieldSymbols[i] = y
				fieldTypes[i] = y.Type()
			}
			unifiedElem, err := UnifyTypes(span, fieldTypes)
			if err != nil {
				return nil, span.Errorf(err, "field group %s", fieldGroupName)
			}
			series := &SeriesSymbol{
				Type_: &SeriesType{Elem: unifiedElem},
				Elem:  fieldSymbols,
			}
			ef = append(ef,
				&FieldSymbol{
					Name:    fieldGroupName,
					Monadic: fieldGroupName == "",
					Value:   series,
				},
			)
		}
	}
	return ef, nil
}
