package fast

import (
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func Insert(node *Node, value Symbol, less *VarietySymbol) *Node {
	if node == nil {
		return &Node{Value: value, Color: Red}
	} else {
		return insertInto(node, value, less)
	}
}

func insertInto(node *Node, value Symbol, less *VarietySymbol) *Node {
	if IsRedRed(node) {
		node = Flip(node)
	}
	node = placeValue(node, value, less)
	if IsBlackRed(node) {
		node = RotateLeft(node)
	}
	if IsLeftRedRed(node) {
		node = RotateRight(node)
	}
	return node
}

func IsRedRed(node *Node) bool {
	return node.Left != nil && node.Left.Color == Red &&
		node.Right != nil && node.Right.Color == Red
}

func IsBlackRed(node *Node) bool {
	return node.Left != nil && node.Left.Color == Black &&
		node.Right != nil && node.Right.Color == Red
}

func IsLeftRedRed(node *Node) bool {
	return node.Left != nil && node.Left.Left != nil &&
		node.Left.Color == Red && node.Left.Left.Color == Red
}

func placeValue(node *Node, value Symbol, less *VarietySymbol) *Node {
	smaller := evokeLess(value, node.Value, less)
	bigger := evokeLess(node.Value, value, less)
	switch {
	case smaller:
		return &Node{
			Value: node.Value,
			Color: node.Color,
			Left:  Insert(node.Left, value, less),
			Right: node.Right,
		}
	case bigger:
		return &Node{
			Value: node.Value,
			Color: node.Color,
			Left:  node.Left,
			Right: Insert(node.Right, value, less),
		}
	default:
		return &Node{
			Value: value,
			Color: node.Color,
			Left:  node.Left,
			Right: node.Right,
		}
	}
}

func evokeLess(u, v Symbol, less *VarietySymbol) bool {
	return AsBool(
		Evoke(
			less,
			EvokeArg{Name: "left", Value: u},
			EvokeArg{Name: "right", Value: v},
		),
	)
}
