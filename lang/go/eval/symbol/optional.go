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
	"reflect"

	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type OptionalType struct {
	Elem Type `ko:"name=elem"`
}

var _ Type = &OptionalType{}

func (*OptionalType) IsType() {}

func (ot *OptionalType) String() string {
	return tree.Sprint(ot)
}

func (ot *OptionalType) Splay() tree.Tree {
	return tree.Sometimes{Elem: ot.Elem.Splay()}
}

// GoType returns the Go equivalent of the type.
func (ot *OptionalType) GoType() reflect.Type {
	return reflect.PtrTo(ot.Elem.GoType())
}

// Optionally makes a type optional, unless it is already optional or series.
func Optionally(t Type) Type {
	switch t.(type) {
	case EmptyType:
		return t
	case *OptionalType:
		return t
	case *SeriesType:
		return t
	case BasicType, *StructType, VarietyType, NamedType, *OpaqueType:
		return &OptionalType{Elem: t}
	}
	panic("o")
}
