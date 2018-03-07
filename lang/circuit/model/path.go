package model

import (
	"path"
	"strings"
)

type Path []string

func (p Path) String() string { return strings.Join(p, ".") }

func (p Path) Slash() string { return path.Join(p...) }

func (p Path) Reverse() Path {
	q := make(Path, len(p))
	for i := range p {
		q[i] = p[len(p)-1-i]
	}
	return q
}

func (p Path) Extend(q ...string) Path {
	x := make(Path, len(p)+len(q))
	copy(x[:len(p)], p)
	copy(x[len(p):], q)
	return x
}

func EqualPath(p, q Path) bool {
	if len(p) != len(q) {
		return false
	}
	for i := range p {
		if p[i] != q[i] {
			return false
		}
	}
	return true
}

func JoinPath(path ...Path) Path {
	r := Path{}
	for _, p := range path {
		r = append(r, p...)
	}
	return r
}
