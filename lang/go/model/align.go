package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func GoFieldsOnPath(field []*GoField, path []int) (pathField []*GoPathField) {
	pathField = make([]*GoPathField, len(field))
	for i, f := range field {
		pathField[i] = &GoPathField{Path: path, Field: f}
	}
	return
}

type FromToGoField struct {
	From *GoPathField `ko:"name=from"`
	To   *GoPathField `ko:"name=to"`
}

type OnlyToGoField struct {
	To *GoPathField `ko:"name=to"`
}

func AlignGoFields(span *Span, from, to []*GoPathField) (aligned []*FromToGoField) {
	aligned, from, to = alignMonadics(from, to)
	fromTo, toNotFrom := fieldDiff(from, to)
	for _, fromTo := range fromTo {
		aligned = append(aligned, fromTo)
	}
	for _, toNotFrom := range toNotFrom {
		aligned = append(aligned, &FromToGoField{
			From: &GoPathField{
				Path: nil,
				Field: &GoField{
					Name: fmt.Sprintf("Align_%s", toNotFrom.To.Field.Name),
					Type: NewGoEmpty(span),
					Tag:  toNotFrom.To.Field.Tag,
				},
			},
			To: toNotFrom.To,
		})
	}
	return
}

func alignMonadics(from, to []*GoPathField) (aligned []*FromToGoField, restFrom, restTo []*GoPathField) {
	fromMonadic, fromRest := findMonadic(from)
	toMonadic, toRest := findMonadic(to)
	if fromMonadic != nil && toMonadic != nil {
		return []*FromToGoField{{From: fromMonadic, To: toMonadic}}, fromRest, toRest
	} else {
		return nil, from, to
	}
}

func findMonadic(field []*GoPathField) (monadic *GoPathField, rest []*GoPathField) {
	for _, f := range field {
		if f.Field.IsMonadic() {
			if monadic != nil {
				panic("duplicate monadic")
			} else {
				monadic = f
			}
		} else {
			rest = append(rest, f)
		}
	}
	return
}

// fieldDiff returns the field set difference between from and to.
//	fromTo: field name -> fromField
//	toNotFrom: field name -> toField
func fieldDiff(from, to []*GoPathField) (
	fromTo map[string]*FromToGoField,
	toNotFrom map[string]*OnlyToGoField,
) {
	fromFields := map[string]*GoPathField{}
	for _, fromField := range from {
		fromFields[fromField.Field.KoName()] = fromField
	}
	fromTo, toNotFrom = map[string]*FromToGoField{}, map[string]*OnlyToGoField{}
	for _, toField := range to {
		if fromField, ok := fromFields[toField.Field.KoName()]; ok {
			fromTo[toField.Field.KoName()] = &FromToGoField{From: fromField, To: toField}
		} else {
			toNotFrom[toField.Field.KoName()] = &OnlyToGoField{To: toField}
		}
	}
	return
}
