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
	"fmt"
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type NamedSymbol struct {
	Value reflect.Value `ko:"name=value"`
}

var _ Symbol = &NamedSymbol{}

// DisassembleToGo converts a Ko value into a Go value
func (named *NamedSymbol) DisassembleToGo(span *model.Span) (reflect.Value, error) {
	return named.Value, nil
}

// DisassembleToPB converts a Ko value into a protobuf
func (named *NamedSymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	return DeconstructKind(span, named.Value).DisassembleToPB(span)
}

func (named *NamedSymbol) String() string {
	return tree.Sprint(named.Value.Interface())
}

func (named *NamedSymbol) Equal(span *model.Span, sym Symbol) bool {
	if other, ok := sym.(*NamedSymbol); ok {
		return reflect.DeepEqual(named.Value.Interface(), other.Value.Interface())
	} else {
		return false
	}
}

func (named *NamedSymbol) Splay() tree.Tree {
	return tree.Splay(named.Value.Interface())
}

func (named *NamedSymbol) Hash(span *model.Span) model.ID {
	return DeconstructKind(span, named.Value).Hash(span)
}

func (named *NamedSymbol) GoType() reflect.Type {
	return named.Value.Type()
}

func (named *NamedSymbol) Type() Type {
	return NamedType{
		Type: named.GoType(),
	}
}

func (named *NamedSymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return singletonSeries(named)
}

func (named *NamedSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to named")
}

func (named *NamedSymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return named, nil, nil
	} else {
		if step, err := named.Walk(span, path[0]); err != nil {
			return nil, nil, err
		} else {
			return step.Select(span, path[1:])
		}
	}
}

func (named *NamedSymbol) Walk(span *model.Span, field string) (Symbol, error) {
	v := named.Value
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return EmptySymbol{}, nil
	}
	if v.Kind() != reflect.Struct {
		return nil, span.Errorf(nil, "cannot select %s into %v", field, named)
	}
	if fieldIndex, ok := gate.StripFields(v.Type()).FieldByKoName(field); ok {
		return Deconstruct(span, v.Field(fieldIndex)), nil
	} else {
		return EmptySymbol{}, nil
	}
}

func (named *NamedSymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "named value %v cannot be augmented", named)
}

func (named *NamedSymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "named value %v cannot be invoked", named)
}

type NamedType struct {
	Type reflect.Type `ko:"name=type"`
}

func (named NamedType) IsType() {}

func (named NamedType) String() string {
	return tree.Sprint(named)
}

func (named NamedType) Splay() tree.Tree {
	return tree.NoQuote{String_: fmt.Sprintf("Named<%s.%s>", named.Type.PkgPath(), named.Type.Name())}
}
