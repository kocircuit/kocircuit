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

package syntax

import (
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/lex"
)

var ErrUnexpectedEnd = SyntaxError{Msg: "unexpected end of file"}

type SyntaxError struct {
	Remainder []lex.Lex     `ko:"name=remainder"` // remaining tokens at error
	Msg       string        `ko:"name=msg"`
	Cause     []SyntaxError `ko:"name=cause"`
}

func (e SyntaxError) Depth() int {
	var d int
	for _, f := range e.Cause {
		d = max(d, 1+f.Depth())
	}
	return d
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (e SyntaxError) BottomCause() SyntaxError {
	if len(e.Cause) == 0 {
		return e
	}
	return e.DeepestCause().BottomCause()
}

func (e SyntaxError) DeepestCause() SyntaxError {
	var d SyntaxError
	for _, f := range e.Cause {
		if f.Depth() > d.Depth() {
			d = f
		}
	}
	return d
}

func (e SyntaxError) ShortestRemainder() SyntaxError {
	h := e
	for _, f := range e.Cause {
		g := f.ShortestRemainder()
		if len(g.Remainder) < len(h.Remainder) {
			h = g
		}
	}
	return h
}

func (e SyntaxError) Error() string {
	farthest := e.ShortestRemainder()
	if len(farthest.Remainder) == 0 {
		return fmt.Sprintf("(end of file) %s", e.Msg)
	}
	return fmt.Sprintf("%v: %s", lex.LexUnion(farthest.Remainder...).RegionString(), e.Msg)
}
