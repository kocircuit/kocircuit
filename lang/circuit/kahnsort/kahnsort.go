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

// Package kahnsort provides topological sorting using Kahn's algorithm.
package kahnsort

import (
	"log"
)

// Node represents a node in a directed acyclic graph.
// Node implementations must be comparable using go equality.
type Node interface {
	// Down returns a list of downstream nodes adjacent to this node.
	// The list of adjacent nodes cannot have duplicates.
	Down() []Node
}

// Sort computes a topological sorting of a downstream-directed acyclic graph.
// The set of starting nodes must be maximal in topological order.
func Sort(start []Node) []Node {
	// Initialize algorithm.
	s := &sorting{
		sorted:  make(map[Node]bool),
		staged:  make(map[Node]bool),
		reached: make(map[Node]*uplinks),
	}
	for _, n := range start {
		s.Reach(n)
		delete(s.reached, n)
	}
	for _, n := range start {
		s.Stage(n)
	}
	// Greedily move staged vertices to sorted order.
	for len(s.staged) > 0 {
		u := s.TakeStaged() // Staged nodes have no uplinks to sorted nodes or among themselves.
		for _, d := range u.Down() {
			if s.RemoveEdge(u, d) {
				s.Stage(d)
			}
		}
		s.sorted[u], s.result = true, append(s.result, u)
	}
	if len(s.reached) != 0 {
		log.Panic("not all reached nodes sorted")
	}
	return s.result
}

type sorting struct {
	sorted  map[Node]bool
	staged  map[Node]bool
	queue   []Node            // queue of staged nodes
	reached map[Node]*uplinks // known multi-uplinks for unstaged/unsorted nodes
	result  []Node
}

type uplinks struct {
	Node []Node
}

func (u *uplinks) Empty() bool {
	return len(u.Node) == 0
}

func (u *uplinks) Remove(m Node) bool {
	for i, n := range u.Node {
		if n == m {
			u.Node[i] = u.Node[len(u.Node)-1]
			u.Node = u.Node[:len(u.Node)-1]
			return true
		}
	}
	return false
}

func (s *sorting) Stage(u Node) {
	if s.staged[u] || s.sorted[u] {
		log.Panic("node already staged")
	}
	s.staged[u] = true
	s.queue = append(s.queue, u)
}

// Reach creates uplinks for all nodes transitively reachable from u.
// Reach returns the slice of uplinks for r.
func (s *sorting) Reach(r Node) *uplinks {
	if u, ok := s.reached[r]; ok {
		return u
	}
	for _, d := range r.Down() {
		du := s.Reach(d)
		du.Node = append(du.Node, r)
	}
	s.reached[r] = &uplinks{}
	return s.reached[r]
}

func (s *sorting) TakeStaged() Node {
	taken := s.queue[0]
	s.queue = s.queue[1:]
	delete(s.staged, taken)
	return taken
}

func (s *sorting) RemoveEdge(up, down Node) (stageable bool) {
	du := s.reached[down]
	if !du.Remove(up) {
		log.Panic("upstream node is not within downstream node's known uplinks")
	}
	if du.Empty() {
		delete(s.reached, down)
		return true
	}
	return false
}
