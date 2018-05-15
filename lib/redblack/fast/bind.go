package fast

import (
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoInsert))
	RegisterEvalGate(new(GoFind))
}

type GoInsert struct {
	Node  *Node  `ko:"name=node"`
	Value Symbol `ko:"name=value"`
	Less  Symbol `ko:"name=Less"`
}

func (g *GoInsert) Play(ctx *runtime.Context) *Node {
	return Insert(g.Node, g.Value, g.Less.(*VarietySymbol))
}

type GoFind struct {
	Node  *Node  `ko:"name=node"`
	Value Symbol `ko:"name=value"`
	Less  Symbol `ko:"name=Less"`
}

func (g *GoFind) Play(ctx *runtime.Context) Symbol {
	return Find(g.Node, g.Value, g.Less.(*VarietySymbol))
}
