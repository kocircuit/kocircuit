package syntax

import (
	"fmt"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
)

func parseBracket(suffix []Lex) (bra Bracket, remain []Lex, err error) {
	if len(suffix) < 1 {
		return Bracket{}, suffix, ErrUnexpectedEnd
	}
	bra, ok := suffix[0].(Bracket)
	if !ok {
		return Bracket{}, suffix, SyntaxError{Remainder: suffix, Msg: "not a bracket"}
	}
	return bra, suffix[1:], nil
}

func matchKeyword(key string, suffix []Lex) (remain []Lex, err error) {
	if len(suffix) < 1 {
		return suffix, ErrUnexpectedEnd
	}
	tok, ok := suffix[0].(Token)
	if !ok {
		return suffix, SyntaxError{Remainder: suffix, Msg: "not a token"}
	}
	nameKey := Name{key}
	if tok.Char != nameKey {
		return suffix, SyntaxError{Remainder: suffix, Msg: fmt.Sprintf("not %q", key)}
	}
	return suffix[1:], nil
}

func matchPunc(key string, suffix []Lex) (remain []Lex, err error) {
	if len(suffix) < 1 {
		return suffix, ErrUnexpectedEnd
	}
	if !isPunc(key, suffix[0]) {
		return suffix, SyntaxError{Remainder: suffix, Msg: fmt.Sprintf("not a %q", key)}
	}
	return suffix[1:], nil
}

func isWhitespace(z Lex) (string, bool) {
	tok, ok := z.(Token)
	if !ok {
		return "", false
	}
	white, ok := tok.Char.(White)
	if !ok {
		return "", false
	}
	return white.String, true
}

func isPunc(key string, z Lex) bool {
	tok, ok := z.(Token)
	if !ok {
		return false
	}
	puncKey := Punc{key}
	return tok.Char == puncKey
}

func parseString(suffix []Lex) (parsed string, remain []Lex, err error) {
	if len(suffix) < 1 {
		return "", suffix, ErrUnexpectedEnd
	}
	tok, ok := suffix[0].(Token)
	if !ok {
		return "", suffix, SyntaxError{Remainder: suffix, Msg: "not a token"}
	}
	s, ok := tok.Char.(LexString)
	if !ok {
		return "", suffix, SyntaxError{Remainder: suffix, Msg: "not a string"}
	}
	return s.String, suffix[1:], nil
}

// Literal is Syntax.
type Literal struct {
	Value Char `ko:"name=value"` // String, Int64, or Float64
	Lex   `ko:"name=lex"`
}

func parseLiteral(suffix []Lex) (literal Literal, remain []Lex, err error) {
	negative := false
	remain = suffix
	if remain, err = matchPunc("-", remain); err == nil {
		negative = true
	}
	if len(remain) < 1 {
		return Literal{}, suffix, ErrUnexpectedEnd
	}
	tok, ok := remain[0].(Token)
	if !ok {
		return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "not a token"}
	}
	switch char := tok.Char.(type) {
	case LexInteger:
		if negative {
			return Literal{Value: char.Negative(), Lex: tok}, remain[1:], nil
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	case LexFloat:
		if negative {
			return Literal{Value: char.Negative(), Lex: tok}, remain[1:], nil
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	case LexString:
		if negative {
			return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "negating a string literal"}
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	}
	return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "not a literal"}
}

func parseName(suffix []Lex) (parsed Ref, remain []Lex, err error) {
	if len(suffix) < 1 {
		return Ref{}, suffix, ErrUnexpectedEnd
	}
	var ok bool
	var name string
	if name, ok = isName(suffix[0]); !ok {
		return Ref{}, suffix, SyntaxError{Remainder: suffix, Msg: "not a name"}
	}
	return Ref{
		Lex:  suffix[0],
		Path: []string{name},
	}, suffix[1:], nil
}

func isName(z Lex) (name string, ok bool) {
	tok, ok := z.(Token)
	if !ok {
		return "", false
	}
	n, ok := tok.Char.(Name)
	if !ok {
		return "", false
	}
	return n.String, true
}

type Ref struct {
	Lex  `ko:"name=lex"`
	Path []string `ko:"name=path"`
}

func (ref Ref) IsEmpty() bool {
	return len(ref.Path) == 0 || ref.Path[0] == ""
}

func (ref Ref) Name() string {
	return ref.Join("_")
}

func (ref Ref) Join(with string) string {
	return strings.Join(ref.Path, with)
}

func ParseRef(suffix []Lex) (ref Ref, remain []Lex, err error) {
	remain = suffix
	if ref, remain, err = parseName(remain); err != nil {
		return Ref{}, suffix, err
	}
	for len(remain) >= 2 {
		if !isPunc(".", remain[0]) {
			break
		}
		var ok bool
		var name string
		if name, ok = isName(remain[1]); !ok {
			break
		}
		ref = Ref{
			Lex:  RegionUnion(ref.Lex, remain[0], remain[1]),
			Path: append(ref.Path, name),
		}
		remain = remain[2:]
	}
	return ref, remain, nil
}
