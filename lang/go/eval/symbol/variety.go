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
	"github.com/golang/protobuf/proto"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func MakeVarietySymbol(macro eval.Macro, arg FieldSymbols) *VarietySymbol {
	return &VarietySymbol{Macro: macro, Arg: arg}
}

func IsVarietySymbol(sym Symbol) bool {
	_, isVty := sym.Type().(VarietyType)
	return isVty
}

type VarietySymbol struct {
	Macro eval.Macro   `ko:"name=macro"`
	Arg   FieldSymbols `ko:"name=arg"`
}

var _ Symbol = &VarietySymbol{}

func (vty *VarietySymbol) Dismentle(span *model.Span) (pkgPath, funcName string, arg *StructSymbol, err error) {
	if interpretFu, ok := vty.Macro.(InterpretMacro); !ok { // if vty points to a circuit
		return "", "", nil, span.Errorf(nil, "variety is not underlied by a function")
	} else {
		pkgPath, funcName = interpretFu.InterpretFunc()
		arg = MakeStructSymbol(vty.Arg)
		return pkgPath, funcName, arg, nil
	}
}

type InterpretMacro interface {
	InterpretFunc() (pkgPath, funcName string)
}

func (vty *VarietySymbol) DisassembleToPB(span *model.Span) (*pb.Symbol, error) {
	if pkgPath, funcName, _, err := vty.Dismentle(span); err != nil {
		return nil, span.Errorf(err, "dismentling variety")
	} else {
		fields, err := disassembleFieldSymbolsToPB(span, vty.Arg)
		if err != nil {
			return nil, err
		}
		dis := &pb.SymbolVariety{
			PkgPath:  proto.String(pkgPath),
			FuncName: proto.String(funcName),
			Arg:      fields,
		}
		return &pb.Symbol{
			Symbol: &pb.Symbol_Variety{Variety: dis},
		}, nil
	}
}

func (vty *VarietySymbol) String() string {
	return tree.Sprint(vty)
}

func (vty *VarietySymbol) Equal(span *model.Span, sym Symbol) bool {
	if other, ok := sym.(*VarietySymbol); ok {
		return vty.Macro.MacroID() == other.Macro.MacroID() &&
			FieldSymbolsEqual(span, vty.Arg, other.Arg)
	} else {
		return false
	}
}

func (vty *VarietySymbol) Hash(span *model.Span) model.ID {
	return model.Blend(model.StringID(vty.Macro.MacroID()), FieldSymbolsHash(span, vty.Arg))
}

func (vty *VarietySymbol) LiftToSeries(span *model.Span) *SeriesSymbol {
	return singletonSeries(vty)
}

func (vty *VarietySymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to variety")
}

func (vty *VarietySymbol) Select(span *model.Span, path model.Path) (eval.Shape, eval.Effect, error) {
	if len(path) == 0 {
		return vty, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "variety %v cannot be selected into", vty)
	}
}

func (vty *VarietySymbol) Type() Type {
	return VarietyType{}
}

func (vty *VarietySymbol) Splay() tree.Tree {
	if len(vty.Arg) == 0 {
		return tree.NoQuote{String_: vty.Macro.Help()}
	} else {
		nameTrees := make([]tree.NameTree, len(vty.Arg))
		for i, field := range vty.Arg {
			nameTrees[i] = tree.NameTree{
				Name:    gate.KoGoName{Ko: field.Name},
				Monadic: field.Monadic,
				Tree:    field.Value.Splay(),
			}
		}
		return tree.Parallel{
			Label:   tree.Label{Name: vty.Macro.Help()},
			Bracket: "[]",
			Elem:    nameTrees,
		}
	}
}

type VarietyType struct{}

func (VarietyType) IsType() {}

func (VarietyType) String() string {
	return tree.Sprint(VarietyType{})
}

func (VarietyType) Splay() tree.Tree {
	return tree.NoQuote{String_: "Variety"}
}
