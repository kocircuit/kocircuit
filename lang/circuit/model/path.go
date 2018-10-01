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
