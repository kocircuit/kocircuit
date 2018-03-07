package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (vty *VarietySymbol) Augment(span *Span, knot Knot) (Shape, Effect, error) {
	augmented, err := KnotToFieldSymbols(span, knot)
	if err != nil {
		return nil, nil, err
	}
	aggregate := append(append(FieldSymbols{}, vty.Arg...), augmented...)
	if err = VerifyNoDuplicateFieldSymbol(span, aggregate); err != nil {
		return nil, nil, span.Errorf(err, "augmenting %s", Sprint(vty))
	}
	return MakeVarietySymbol(vty.Macro, aggregate), nil, nil
}

func KnotToFieldSymbols(span *Span, knot Knot) (FieldSymbols, error) {
	ef := FieldSymbols{}
	for _, fieldGroup := range knot.FieldGroup() {
		fieldGroup = FilterEmptyFields(fieldGroup)
		switch len(fieldGroup) {
		case 0:
		case 1:
			y := fieldGroup[0].Shape.(Symbol)
			ef = append(ef,
				&FieldSymbol{
					Name:    fieldGroup[0].Name,
					Monadic: fieldGroup[0].Name == "",
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
				return nil, span.Errorf(err, "field group %s", fieldGroup[0].Name)
			}
			series := &SeriesSymbol{
				Type_: &SeriesType{Elem: unifiedElem},
				Elem:  fieldSymbols,
			}
			ef = append(ef,
				&FieldSymbol{
					Name:    fieldGroup[0].Name,
					Monadic: fieldGroup[0].Name == "",
					Value:   series,
				},
			)
		}
	}
	return ef, nil
}

func FilterEmptyFields(group []Field) (filtered []Field) {
	for _, field := range group {
		if !IsEmptySymbol(field.Shape.(Symbol)) {
			filtered = append(filtered, field)
		}
	}
	return
}

func VerifyNoDuplicateFieldSymbol(span *Span, fields FieldSymbols) error {
	seen := map[string]bool{}
	for _, f := range fields {
		if seen[f.Name] {
			return span.Errorf(nil, "augmenting duplicate field %s", f.Name)
		} else {
			seen[f.Name] = true
		}
	}
	return nil
}
