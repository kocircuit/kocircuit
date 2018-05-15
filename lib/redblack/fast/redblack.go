package fast

import (
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

// XXX: integrate/unify/deconstruct recognize symbol fields
type Node struct {
	Value Symbol `ko:"name=value"`
	Color Color  `ko:"name=color"`
	Left  *Node  `ko:"name=left"`
	Right *Node  `ko:"name=right"`
}

type Color int

const Red Color = 1
const Black Color = 2

func flipColor(c Color) Color {
	switch c {
	case Red:
		return Black
	case Black:
		return Red
	}
	panic("o")
}

func RotateLeft(node *Node) *Node {
	return &Node{
		Value: node.Right.Value,
		Color: node.Color,
		Left: &Node{
			Value: node.Value,
			Color: node.Right.Color,
			Left:  node.Left,
			Right: node.Right.Left,
		},
		Right: node.Right.Right,
	}
}

func RotateRight(node *Node) *Node {
	return &Node{
		Value: node.Left.Value,
		Color: node.Color,
		Left:  node.Left.Left,
		Right: &Node{
			Value: node.Value,
			Color: node.Left.Color,
			Left:  node.Left.Right,
			Right: node.Right,
		},
	}
}

func Flip(node *Node) *Node {
	return &Node{
		Value: node.Value,
		Color: flipColor(node.Color),
		Left:  flipNode(node.Left),
		Right: flipNode(node.Right),
	}
}

func flipNode(node *Node) *Node {
	return &Node{
		Value: node.Value,
		Color: flipColor(node.Color),
		Left:  node.Left,
		Right: node.Right,
	}
}
