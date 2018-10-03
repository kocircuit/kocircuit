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
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func IsEmptySymbol(sym Symbol) bool {
	_, isEmpty := sym.(EmptySymbol)
	return isEmpty
}

func IsEmptyType(t Type) bool {
	_, isEmpty := t.(EmptyType)
	return isEmpty
}

type EmptySymbol struct{}

func (empty EmptySymbol) Disassemble(span *model.Span) (*pb.Symbol, error) {
	return nil, nil
}

func (empty EmptySymbol) String() string {
	return tree.Sprint(empty)
}

func (empty EmptySymbol) Equal(span *model.Span, sym Symbol) bool {
	_, ok := sym.(EmptySymbol)
	return ok
}

func (empty EmptySymbol) Hash(span *model.Span) model.ID {
	return model.StringID("#empty")
}

func (empty EmptySymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return EmptySeries
}

func (empty EmptySymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to empty")
}

func (empty EmptySymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	return empty, nil, nil
}

func (empty EmptySymbol) Augment(span *model.Span, _ eval.Fields) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "empty value cannot be augmented")
}

func (empty EmptySymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "empty value cannot be invoked")
}

func (empty EmptySymbol) Type() Type {
	return EmptyType{}
}

func (empty EmptySymbol) Splay() tree.Tree {
	return tree.NoQuote{String_: "empty"}
}

type EmptyType struct{}

func (EmptyType) IsType() {}

func (EmptyType) String() string {
	return tree.Sprint(EmptyType{})
}

func (EmptyType) Splay() tree.Tree {
	return tree.NoQuote{String_: "Empty"}
}
