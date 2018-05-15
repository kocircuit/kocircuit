package fast

import (
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func Find(node *Node, value Symbol, less *VarietySymbol) Symbol {
	if node == nil {
		return EmptySymbol{}
	} else {
		smaller := evokeLess(value, node.Value, less)
		bigger := evokeLess(node.Value, value, less)
		switch {
		case smaller:
			return Find(node.Left, value, less)
		case bigger:
			return Find(node.Right, value, less)
		default:
			return node.Value
		}
	}
}
