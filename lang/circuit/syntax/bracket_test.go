package syntax

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
)

var testBracket = []struct{ In, Out []Lex }{
	{
		In: []Lex{
			Token{Char: LexString{String: "a"}},
			Token{Char: Punc{String: "{"}},
			Token{Char: LexInteger{Int64: 1}},
			Token{Char: LexInteger{Int64: 2}},
			Token{Char: LexInteger{Int64: 3}},
			Token{Char: Punc{String: "}"}},
			Token{Char: LexString{String: "b"}},
		},
		Out: []Lex{
			Token{Char: LexString{String: "a"}},
			Bracket{
				Left:  Token{Char: Punc{String: "{"}},
				Right: Token{Char: Punc{String: "}"}},
				Body: []Lex{
					Token{Char: LexInteger{Int64: 1}},
					Token{Char: LexInteger{Int64: 2}},
					Token{Char: LexInteger{Int64: 3}},
				},
			},
			Token{Char: LexString{String: "b"}},
		},
	},
	{
		In: []Lex{
			Token{Char: LexString{String: "a"}},
			Token{Char: Punc{String: "{"}},
			Token{Char: LexInteger{Int64: 1}},
			Token{Char: Punc{String: "["}},
			Token{Char: LexInteger{Int64: 2}},
			Token{Char: Punc{String: "]"}},
			Token{Char: LexInteger{Int64: 3}},
			Token{Char: Punc{String: "}"}},
			Token{Char: LexString{String: "b"}},
		},
		Out: []Lex{
			Token{Char: LexString{String: "a"}},
			Bracket{
				Left: Token{Char: Punc{String: "{"}},
				Body: []Lex{
					Token{Char: LexInteger{Int64: 1}},
					Bracket{
						Left:  Token{Char: Punc{String: "["}},
						Body:  []Lex{Token{Char: LexInteger{Int64: 2}}},
						Right: Token{Char: Punc{String: "]"}},
					},
					Token{Char: LexInteger{Int64: 3}},
				},
				Right: Token{Char: Punc{String: "}"}},
			},
			Token{Char: LexString{String: "b"}},
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
