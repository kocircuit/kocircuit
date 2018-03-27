package symbol

import (
	"github.com/golang/protobuf/proto"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func MakeVarietySymbol(macro Macro, arg FieldSymbols) *VarietySymbol {
	return &VarietySymbol{Macro: macro, Arg: arg}
}

func IsVarietySymbol(sym Symbol) bool {
	_, isVty := sym.Type().(VarietyType)
	return isVty
}

type VarietySymbol struct {
	Macro Macro        `ko:"name=macro"`
	Arg   FieldSymbols `ko:"name=arg"`
}

type InterpretMacro interface {
	InterpretFunc() (pkgPath, funcName string)
}

func (vty *VarietySymbol) Disassemble(span *Span) *pb.Symbol {
	if interpretFu, ok := vty.Macro.(InterpretMacro); ok { // if vty points to a circuit
		pkgPath, funcName := interpretFu.InterpretFunc()
		dis := &pb.SymbolVariety{
			PkgPath:  proto.String(pkgPath),
			FuncName: proto.String(funcName),
			Arg:      DisassembleFieldSymbols(span, vty.Arg),
		}
		return &pb.Symbol{
			Symbol: &pb.Symbol_Variety{Variety: dis},
		}
	} else {
		return nil // non-function macros are not disassembled
	}
}

func (vty *VarietySymbol) String() string {
	return Sprint(vty)
}

func (vty *VarietySymbol) Equal(sym Symbol) bool {
	if other, ok := sym.(*VarietySymbol); ok {
		return vty.Macro.MacroID() == other.Macro.MacroID() && FieldSymbolsEqual(vty.Arg, other.Arg)
	} else {
		return false
	}
}

func (vty *VarietySymbol) Hash() string {
	return Mix(vty.Macro.MacroID(), FieldSymbolsHash(vty.Arg))
}

func (vty *VarietySymbol) LiftToSeries(span *Span) *SeriesSymbol {
	return singletonSeries(vty)
}

func (vty *VarietySymbol) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return vty, nil, nil
	} else {
		return nil, nil, span.Errorf(nil, "variety %v cannot be selected into", vty)
	}
}

func (vty *VarietySymbol) Type() Type {
	return VarietyType{}
}

func (vty *VarietySymbol) Splay() Tree {
	nameTrees := make([]NameTree, len(vty.Arg))
	for i, field := range vty.Arg {
		nameTrees[i] = NameTree{
			Name:    gate.KoGoName{Ko: field.Name},
			Monadic: field.Monadic,
			Tree:    field.Value.Splay(),
		}
	}
	return Parallel{
		Label:   Label{Name: vty.Macro.Help()},
		Bracket: "[]",
		Elem:    nameTrees,
	}
}

type VarietyType struct{}

func (VarietyType) IsType() {}

func (VarietyType) String() string {
	return Sprint(VarietyType{})
}

func (VarietyType) Splay() Tree {
	return NoQuote{"Variety"}
}
