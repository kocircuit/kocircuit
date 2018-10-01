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
	"strings"

	"github.com/kocircuit/kocircuit/lang/circuit/lex"
)

func parseBracket(suffix []lex.Lex) (bra Bracket, remain []lex.Lex, err error) {
	if len(suffix) < 1 {
		return Bracket{}, suffix, ErrUnexpectedEnd
	}
	bra, ok := suffix[0].(Bracket)
	if !ok {
		return Bracket{}, suffix, SyntaxError{Remainder: suffix, Msg: "not a bracket"}
	}
	return bra, suffix[1:], nil
}

func matchKeyword(key string, suffix []lex.Lex) (remain []lex.Lex, err error) {
	if len(suffix) < 1 {
		return suffix, ErrUnexpectedEnd
	}
	tok, ok := suffix[0].(lex.Token)
	if !ok {
		return suffix, SyntaxError{Remainder: suffix, Msg: "not a token"}
	}
	nameKey := lex.Name{String: key}
	if tok.Char != nameKey {
		return suffix, SyntaxError{Remainder: suffix, Msg: fmt.Sprintf("not %q", key)}
	}
	return suffix[1:], nil
}

func matchPunc(key string, suffix []lex.Lex) (remain []lex.Lex, err error) {
	if len(suffix) < 1 {
		return suffix, ErrUnexpectedEnd
	}
	if !isPunc(key, suffix[0]) {
		return suffix, SyntaxError{Remainder: suffix, Msg: fmt.Sprintf("not a %q", key)}
	}
	return suffix[1:], nil
}

func isWhitespace(z lex.Lex) (string, bool) {
	tok, ok := z.(lex.Token)
	if !ok {
		return "", false
	}
	white, ok := tok.Char.(lex.White)
	if !ok {
		return "", false
	}
	return white.String, true
}

func isPunc(key string, z lex.Lex) bool {
	tok, ok := z.(lex.Token)
	if !ok {
		return false
	}
	puncKey := lex.Punc{String: key}
	return tok.Char == puncKey
}

func parseString(suffix []lex.Lex) (parsed string, remain []lex.Lex, err error) {
	if len(suffix) < 1 {
		return "", suffix, ErrUnexpectedEnd
	}
	tok, ok := suffix[0].(lex.Token)
	if !ok {
		return "", suffix, SyntaxError{Remainder: suffix, Msg: "not a token"}
	}
	s, ok := tok.Char.(lex.LexString)
	if !ok {
		return "", suffix, SyntaxError{Remainder: suffix, Msg: "not a string"}
	}
	return s.String, suffix[1:], nil
}

// Literal is Syntax.
type Literal struct {
	Value   lex.Char `ko:"name=value"` // String, Int64, or Float64
	lex.Lex `ko:"name=lex"`
}

func parseLiteral(suffix []lex.Lex) (literal Literal, remain []lex.Lex, err error) {
	negative := false
	remain = suffix
	if remain, err = matchPunc("-", remain); err == nil {
		negative = true
	}
	if len(remain) < 1 {
		return Literal{}, suffix, ErrUnexpectedEnd
	}
	tok, ok := remain[0].(lex.Token)
	if !ok {
		return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "not a token"}
	}
	switch char := tok.Char.(type) {
	case lex.LexInteger:
		if negative {
			return Literal{Value: char.Negative(), Lex: tok}, remain[1:], nil
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	case lex.LexFloat:
		if negative {
			return Literal{Value: char.Negative(), Lex: tok}, remain[1:], nil
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	case lex.LexString:
		if negative {
			return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "negating a string literal"}
		} else {
			return Literal{Value: char, Lex: tok}, remain[1:], nil
		}
	}
	return Literal{}, suffix, SyntaxError{Remainder: remain, Msg: "not a literal"}
}

func parseName(suffix []lex.Lex) (parsed Ref, remain []lex.Lex, err error) {
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

func isName(z lex.Lex) (name string, ok bool) {
	tok, ok := z.(lex.Token)
	if !ok {
		return "", false
	}
	n, ok := tok.Char.(lex.Name)
	if !ok {
		return "", false
	}
	return n.String, true
}

type Ref struct {
	lex.Lex `ko:"name=lex"`
	Path    []string `ko:"name=path"`
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

func ParseRef(suffix []lex.Lex) (ref Ref, remain []lex.Lex, err error) {
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
			Lex:  lex.RegionUnion(ref.Lex, remain[0], remain[1]),
			Path: append(ref.Path, name),
		}
		remain = remain[2:]
	}
	return ref, remain, nil
}
