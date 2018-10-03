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
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type OpaqueSymbol struct {
	Value reflect.Value `ko:"name=value"`
}

var _ Symbol = &OpaqueSymbol{}

func (opaque *OpaqueSymbol) Interface() interface{} {
	return opaque.Value.Interface()
}

// DisassembleToGo converts a Ko value into a Go value
func (opaque *OpaqueSymbol) DisassembleToGo(span *model.Span) (reflect.Value, error) {
	return reflect.Value{}, span.Errorf(nil, "cannot disassemble opaque symbol %v", opaque)
}

// DisassembleToPB converts a Ko value into a protobuf
func (opaque *OpaqueSymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	return nil, span.Errorf(nil, "cannot disassemble opaque symbol %v", opaque)
}

func (opaque *OpaqueSymbol) String() string {
	return tree.Sprint(opaque)
}

func (opaque *OpaqueSymbol) Equal(span *model.Span, sym Symbol) bool {
	if other, ok := sym.(*OpaqueSymbol); ok {
		return opaque.Value.Interface() == other.Value.Interface()
	} else {
		return false
	}
}

func (opaque *OpaqueSymbol) Hash(span *model.Span) model.ID {
	return model.StringID("#█")
}

func (opaque *OpaqueSymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return singletonSeries(opaque)
}

func (opaque *OpaqueSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to opaque")
}

func (opaque *OpaqueSymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return opaque, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "opaque value %v cannot be selected into", opaque)
	}
}

func (opaque *OpaqueSymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "opaque value %v cannot be augmented", opaque)
}

func (opaque *OpaqueSymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "opaque value %v cannot be invoked", opaque)
}

func (opaque *OpaqueSymbol) GoType() reflect.Type {
	return opaque.Value.Type()
}

func (opaque *OpaqueSymbol) Type() Type {
	return &OpaqueType{Type: opaque.Value.Type()}
}

func (opaque *OpaqueSymbol) Splay() tree.Tree {
	return opaque.Type().Splay()
}

type OpaqueType struct {
	Type reflect.Type `ko:"name=type"`
}

func (opaque *OpaqueType) IsType() {}

func (opaque *OpaqueType) String() string {
	return tree.Sprint(opaque)
}

func (opaque *OpaqueType) Splay() tree.Tree {
	return tree.NoQuote{String_: fmt.Sprintf("Opaque<%v>", opaque.Type)}
}
