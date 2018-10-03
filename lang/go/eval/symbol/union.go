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
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (ctx *typingCtx) UnifyStruct(x, y *StructType) (*StructType, error) {
	union := unionFields{}
	for _, xfield := range x.Field {
		union = union.AddXField(xfield)
	}
	for _, yfield := range y.Field {
		union = union.AddYField(yfield)
	}
	fields := []*FieldType{}
	for _, union := range union {
		if unified, err := ctx.Refine(union.Name).Unify(union.X, union.Y); err != nil {
			return nil, ctx.Errorf(nil, "cannot unify field %s types %s and %s",
				tree.Sprint(union.Name), tree.Sprint(union.X), tree.Sprint(union.Y))
		} else {
			fields = append(fields,
				&FieldType{
					Name:  union.Name,
					Type_: unified,
				})
		}
	}
	return &StructType{fields}, nil
}

type unionField struct {
	Name string `ko:"name=name"` // ko name
	X    Type   `ko:"name=x"`
	Y    Type   `ko:"name=y"`
}

type unionFields []*unionField

func (fields unionFields) Find(name string) (int, bool) {
	for i, field := range fields {
		if field.Name == name {
			return i, true
		}
	}
	return -1, false
}

func (fields unionFields) AddXField(xfield *FieldType) unionFields {
	if index, found := fields.Find(xfield.Name); found {
		fields[index].X = xfield.Type_
		return fields
	} else {
		return append(fields,
			&unionField{Name: xfield.Name, X: xfield.Type_, Y: EmptyType{}},
		)
	}
}

func (fields unionFields) AddYField(yfield *FieldType) unionFields {
	if index, found := fields.Find(yfield.Name); found {
		fields[index].Y = yfield.Type_
		return fields
	} else {
		return append(fields,
			&unionField{Name: yfield.Name, X: EmptyType{}, Y: yfield.Type_},
		)
	}
}
