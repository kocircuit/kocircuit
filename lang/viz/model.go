// Package viz provides a protocol for defining circuit graphs for visualisation.
package viz

type Graph struct {
	Node []*Node `ko:"name=node"`
}

type Node struct {
	Type  string   `ko:"name=type"`
	Name  string   `ko:"name=name"`
	Input []*Input `ko:"name=input"`
}

type Input struct {
	Label string `ko:"name=label"`
	Node  *Node  `ko:"name=node"`
}
