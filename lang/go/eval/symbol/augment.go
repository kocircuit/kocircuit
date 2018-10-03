//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package symbol

import (
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (vty *VarietySymbol) Augment(span *model.Span, fields eval.Fields) (eval.Shape, eval.Effect, error) {
	augmented, err := GroupFieldsToSymbols(span, fields)
	if err != nil {
		return nil, nil, err
	}
	aggregate := append(append(FieldSymbols{}, vty.Arg...), augmented...)
	if err = VerifyNoDuplicateFieldSymbol(span, aggregate); err != nil {
		return nil, nil, span.Errorf(err, "augmenting %s", tree.Sprint(vty))
	}
	return MakeVarietySymbol(vty.Macro, aggregate), nil, nil
}

func GroupFieldsToSymbols(span *model.Span, fields eval.Fields) (FieldSymbols, error) {
	evalFields := FieldSymbols{}
	for _, fieldGroup := range fields.FieldGroup() {
		groupName := fieldGroup[0].Name
		fieldGroup = FilterEmptyEvalFields(fieldGroup)
		switch len(fieldGroup) {
		case 0: // add a field with an empty symbol value (useful to All macro)
			evalFields = append(evalFields,
				&FieldSymbol{Name: groupName, Monadic: groupName == "", Value: EmptySymbol{}},
			)
		case 1:
			y := fieldGroup[0].Shape.(Symbol)
			evalFields = append(evalFields,
				&FieldSymbol{Name: groupName, Monadic: groupName == "", Value: y},
			)
		default:
			repeatedSymbols := make(Symbols, len(fieldGroup))
			for i, f := range fieldGroup {
				repeatedSymbols[i] = f.Shape.(Symbol)
			}
			if repeated, err := MakeSeriesSymbol(span, repeatedSymbols); err != nil {
				return nil, span.Errorf(err, "repeated field group %s", groupName)
			} else {
				evalFields = append(evalFields,
					&FieldSymbol{Name: groupName, Monadic: groupName == "", Value: repeated},
				)
			}
		}
	}
	return evalFields, nil
}

func FilterEmptyEvalFields(group []eval.Field) (filtered []eval.Field) {
	for _, field := range group {
		if !IsEmptySymbol(field.Shape.(Symbol)) {
			filtered = append(filtered, field)
		}
	}
	return
}

func VerifyNoDuplicateFieldSymbol(span *model.Span, fieldSymbols FieldSymbols) error {
	seen := map[string]bool{}
	for _, f := range fieldSymbols {
		if seen[f.Name] {
			return span.Errorf(nil, "augmenting duplicate field %s", f.Name)
		} else {
			seen[f.Name] = true
		}
	}
	return nil
}
