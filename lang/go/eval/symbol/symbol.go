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

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Symbol implementations:
//	BasicSymbol, EmptySymbol,
// *SeriesSymbol, *StructSymbol
// *NamedSymbol, *OpaqueSymbol, *MapSymbol
// *VarietySymbol,
// *BlobSymbol
type Symbol interface {
	eval.Shape   // String, Select, Augment, Invoke, Link
	tree.Splayer // Splay
	Type() Type
	Hash(*model.Span) model.ID
	Equal(*model.Span, Symbol) bool
	LiftToSeries(*model.Span) *SeriesSymbol
	Disassemble(*model.Span) (*pb.Symbol, error)
}

var (
	symbolPtr    *Symbol
	typeOfSymbol = reflect.TypeOf((*Symbol)(nil)).Elem()
)

type Symbols []Symbol

func (syms Symbols) Types() Types {
	types := make(Types, len(syms))
	for i, sym := range syms {
		types[i] = sym.Type()
	}
	return types
}
