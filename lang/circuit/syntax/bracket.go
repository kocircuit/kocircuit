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

// Bracket is a Lex.
type Bracket struct {
	Left  lex.Token `ko:"name=left"`
	Right lex.Token `ko:"name=right"`
	Body  []lex.Lex `ko:"name=body"` // []Token
}

func (bra Bracket) FilePath() string {
	return bra.Left.FilePath()
}

func (bra Bracket) StartPosition() lex.Position {
	return bra.Left.StartPosition()
}

func (bra Bracket) EndPosition() lex.Position {
	return bra.Right.EndPosition()
}

func (bra Bracket) RegionString() string {
	return bra.Left.RegionString()
}

func (bra Bracket) Type() string { // "{}" or "[]" or "()"
	return bra.Left.Char.(lex.Punc).String + bra.Right.Char.(lex.Punc).String
}

// An alternative parallel global algorithm for bracket folding:
//	(a) filter non-brackets
//	(b) greedily collapse matching brackets, resulting in bracket structures, recurse

func FoldBracket(suffix []lex.Lex) (r []lex.Lex, err error) { // [](Token or Bracket)
	for len(suffix) > 0 {
		var block []lex.Lex
		if block, suffix = sliceBlock(suffix); len(block) > 0 {
			r = append(r, block...)
			continue
		}
		var bracket []lex.Lex
		if bracket, suffix, err = sliceBracket(suffix); err != nil {
			return nil, err
		} else if len(bracket) > 0 {
			if fold, err := FoldBracket(bracket[1 : len(bracket)-1]); err != nil {
				return nil, err
			} else {
				r = append(r,
					Bracket{
						Left:  bracket[0].(lex.Token),
						Right: bracket[len(bracket)-1].(lex.Token),
						Body:  fold,
					},
				)
				continue
			}
		}
	}
	return r, nil
}

// sliceBlock splits suffix into two pieces: block and tail.
// The block captures all tokens before the first bracket token,
// The tail captures the rest.
func sliceBlock(suffix []lex.Lex) (block, tail []lex.Lex) {
	for i, z := range suffix {
		if len(isBracketChar(z.(lex.Token).Char)) > 0 {
			return suffix[:i], suffix[i:]
		}
	}
	return suffix, nil
}

// isBracketChar returns a non-empty bracket string if z is Punc{bracket}.
func isBracketChar(z lex.Char) string {
	p, ok := z.(lex.Punc)
	if !ok {
		return ""
	}
	switch p.String {
	case "{", "}", "[", "]", "(", ")":
		return p.String
	}
	return ""
}

func sliceBracket(suffix []lex.Lex) (bracket, tail []lex.Lex, err error) {
	if len(suffix) == 0 {
		return nil, nil, nil
	}
	var stack []string // stack of brackets
	for i, z := range suffix {
		tok := z.(lex.Token)
		if bra := isBracketChar(tok.Char); len(bra) > 0 {
			if len(stack) == 0 || !cancelBracket(stack[len(stack)-1], bra) {
				stack = append(stack, bra)
			} else {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 {
					return suffix[:i+1], suffix[i+1:], nil
				}
			}
		}
	}
	if len(stack) == 0 {
		panic("unexpected")
	}
	return nil, nil, fmt.Errorf("bracket imbalance within %s", lex.LexUnion(suffix...).RegionString())
}

func cancelBracket(left, right string) bool {
	switch {
	case left == "{" && right == "}":
		return true
	case left == "[" && right == "]":
		return true
	case left == "(" && right == ")":
		return true
	}
	return false
}
