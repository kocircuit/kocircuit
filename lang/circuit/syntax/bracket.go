package syntax

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
)

// Bracket is a Lex.
type Bracket struct {
	Left  Token `ko:"name=left"`
	Right Token `ko:"name=right"`
	Body  []Lex `ko:"name=body"` // []Token
}

func (bra Bracket) FilePath() string {
	return bra.Left.FilePath()
}

func (bra Bracket) StartPosition() Position {
	return bra.Left.StartPosition()
}

func (bra Bracket) EndPosition() Position {
	return bra.Right.EndPosition()
}

func (bra Bracket) RegionString() string {
	return bra.Left.RegionString()
}

func (bra Bracket) Type() string { // "{}" or "[]" or "()"
	return bra.Left.Char.(Punc).String + bra.Right.Char.(Punc).String
}

// An alternative parallel global algorithm for bracket folding:
//	(a) filter non-brackets
//	(b) greedily collapse matching brackets, resulting in bracket structures, recurse

func FoldBracket(suffix []Lex) (r []Lex, err error) { // [](Token or Bracket)
	for len(suffix) > 0 {
		var block []Lex
		if block, suffix = sliceBlock(suffix); len(block) > 0 {
			r = append(r, block...)
			continue
		}
		var bracket []Lex
		if bracket, suffix, err = sliceBracket(suffix); err != nil {
			return nil, err
		} else if len(bracket) > 0 {
			if fold, err := FoldBracket(bracket[1 : len(bracket)-1]); err != nil {
				return nil, err
			} else {
				r = append(r,
					Bracket{
						Left:  bracket[0].(Token),
						Right: bracket[len(bracket)-1].(Token),
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
func sliceBlock(suffix []Lex) (block, tail []Lex) {
	for i, z := range suffix {
		if len(isBracketChar(z.(Token).Char)) > 0 {
			return suffix[:i], suffix[i:]
		}
	}
	return suffix, nil
}

// isBracketChar returns a non-empty bracket string if z is Punc{bracket}.
func isBracketChar(z Char) string {
	p, ok := z.(Punc)
	if !ok {
		return ""
	}
	switch p.String {
	case "{", "}", "[", "]", "(", ")":
		return p.String
	}
	return ""
}

func sliceBracket(suffix []Lex) (bracket, tail []Lex, err error) {
	if len(suffix) == 0 {
		return nil, nil, nil
	}
	var stack []string // stack of brackets
	for i, z := range suffix {
		tok := z.(Token)
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
	return nil, nil, fmt.Errorf("bracket imbalance within %s", LexUnion(suffix...).RegionString())
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
