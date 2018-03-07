package syntax

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
)

var testBracket = []struct{ In, Out []Lex }{
	{
		In: []Lex{
			Token{Char: LexString{"a"}},
			Token{Char: Punc{"{"}},
			Token{Char: LexInteger{Int64: 1}},
			Token{Char: LexInteger{Int64: 2}},
			Token{Char: LexInteger{Int64: 3}},
			Token{Char: Punc{"}"}},
			Token{Char: LexString{"b"}},
		},
		Out: []Lex{
			Token{Char: LexString{"a"}},
			Bracket{
				Left:  Token{Char: Punc{"{"}},
				Right: Token{Char: Punc{"}"}},
				Body: []Lex{
					Token{Char: LexInteger{Int64: 1}},
					Token{Char: LexInteger{Int64: 2}},
					Token{Char: LexInteger{Int64: 3}},
				},
			},
			Token{Char: LexString{"b"}},
		},
	},
	{
		In: []Lex{
			Token{Char: LexString{"a"}},
			Token{Char: Punc{"{"}},
			Token{Char: LexInteger{Int64: 1}},
			Token{Char: Punc{"["}},
			Token{Char: LexInteger{Int64: 2}},
			Token{Char: Punc{"]"}},
			Token{Char: LexInteger{Int64: 3}},
			Token{Char: Punc{"}"}},
			Token{Char: LexString{"b"}},
		},
		Out: []Lex{
			Token{Char: LexString{"a"}},
			Bracket{
				Left: Token{Char: Punc{"{"}},
				Body: []Lex{
					Token{Char: LexInteger{Int64: 1}},
					Bracket{
						Left:  Token{Char: Punc{"["}},
						Body:  []Lex{Token{Char: LexInteger{Int64: 2}}},
						Right: Token{Char: Punc{"]"}},
					},
					Token{Char: LexInteger{Int64: 3}},
				},
				Right: Token{Char: Punc{"}"}},
			},
			Token{Char: LexString{"b"}},
		},
	},
}

func TestBracket(t *testing.T) {
	for i, test := range testBracket {
		if got, err := FoldBracket(test.In); err != nil {
			t.Errorf("test %d: folding (%v)", i, err)
		} else if !IsSubset(test.Out, got) {
			t.Errorf("test %d: expecting %v, got %v", i, test.Out, got)
		}
	}
}
