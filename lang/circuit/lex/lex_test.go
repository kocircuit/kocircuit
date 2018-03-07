package lex

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestLexer(t *testing.T) {
	for i, test := range testLex {
		if !test.enabled {
			continue
		}
		got, err := LexifyString("", test.src)
		if err != nil {
			t.Errorf("test %d: lexing (%v)", i, err)
		}
		if !IsSubset(test.expecting, got) {
			t.Errorf("test %d: expecting %v, got %v", i, Sprint(test.expecting), Sprint(got))
		}
	}
}

var testLex = []struct {
	enabled   bool
	src       string
	expecting []Region
}{
	{
		true,
		"abc\n\n// end comment",
		[]Region{
			Token{
				Char: Name{"abc"}, Text: "abc",
				Lex: &StartEndRegion{Start: Position{"", 0, 1, 1}, End: Position{"", 3, 1, 4}},
			},
			Token{
				Char: Line{}, Text: "\n",
				Lex: &StartEndRegion{Start: Position{"", 3, 1, 4}, End: Position{"", 4, 2, 1}},
			},
			Token{
				Char: Line{}, Text: "\n",
				Lex: &StartEndRegion{Start: Position{"", 4, 2, 1}, End: Position{"", 5, 3, 1}},
			},
			Token{
				Char: Comment{" end comment"}, Text: "// end comment",
				Lex: &StartEndRegion{Start: Position{"", 5, 3, 1}, End: Position{"", 19, 3, 15}},
			},
		},
	},
	{
		true,
		"\"abc\" /* eh */",
		[]Region{
			Token{
				Char: LexString{"abc"}, Text: "\"abc\"",
				Lex: &StartEndRegion{Start: Position{"", 0, 1, 1}, End: Position{"", 5, 1, 6}},
			},
			Token{
				Char: Comment{" eh "}, Text: "/* eh */",
				Lex: &StartEndRegion{Start: Position{"", 6, 1, 7}, End: Position{"", 14, 1, 15}},
			},
		},
	},
	{
		true,
		"abc {12,34}",
		[]Region{
			Token{
				Char: Name{"abc"}, Text: "abc",
				Lex: &StartEndRegion{Start: Position{"", 0, 1, 1}, End: Position{"", 3, 1, 4}},
			},
			Token{
				Char: Punc{"{"}, Text: "{",
				Lex: &StartEndRegion{Start: Position{"", 4, 1, 5}, End: Position{"", 5, 1, 6}},
			},
			Token{
				Char: LexInteger{Int64: 12}, Text: "12",
				Lex: &StartEndRegion{Start: Position{"", 5, 1, 6}, End: Position{"", 7, 1, 8}},
			},
			Token{
				Char: Line{}, Text: ",",
				Lex: &StartEndRegion{Start: Position{"", 7, 1, 8}, End: Position{"", 8, 1, 9}},
			},
			Token{
				Char: LexInteger{Int64: 34}, Text: "34",
				Lex: &StartEndRegion{Start: Position{"", 8, 1, 9}, End: Position{"", 10, 1, 11}},
			},
			Token{
				Char: Punc{"}"}, Text: "}",
				Lex: &StartEndRegion{Start: Position{"", 10, 1, 11}, End: Position{"", 11, 1, 12}},
			},
		},
	},
	{
		true,
		"...:.",
		[]Region{
			Token{
				Char: Punc{"."}, Text: ".",
				Lex: &StartEndRegion{Start: Position{"", 0, 1, 1}, End: Position{"", 1, 1, 2}},
			},
			Token{
				Char: Punc{"."}, Text: ".",
				Lex: &StartEndRegion{Start: Position{"", 1, 1, 2}, End: Position{"", 2, 1, 3}},
			},
			Token{
				Char: Punc{"."}, Text: ".",
				Lex: &StartEndRegion{Start: Position{"", 2, 1, 3}, End: Position{"", 3, 1, 4}},
			},
			Token{
				Char: Punc{":"}, Text: ":",
				Lex: &StartEndRegion{Start: Position{"", 3, 1, 4}, End: Position{"", 4, 1, 5}},
			},
			Token{
				Char: Punc{"."}, Text: ".",
				Lex: &StartEndRegion{Start: Position{"", 4, 1, 5}, End: Position{"", 5, 1, 6}},
			},
		},
	},
}
