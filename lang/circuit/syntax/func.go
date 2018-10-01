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

import "github.com/kocircuit/kocircuit/lang/circuit/lex"

type Design struct {
	Comment string   `ko:"name=comment"`
	Name    Ref      `ko:"name=name"`
	Factor  []Factor `ko:"name=factor"`
	Returns Assembly `ko:"name=returns"`
	lex.Lex `ko:"name=lex"`
}

type Factor struct {
	Comment string          `ko:"name=comment"`
	Name    Ref             `ko:"name=name"`
	Monadic bool            `ko:"name=monadic"`
	lex.Lex `ko:"name=lex"` // comment and name tokens for this factor
}

func parseDesign(needsLine bool, suffix []lex.Lex) (parsed []Design, remain []lex.Lex, err error) {
	remain = suffix
	var lines int
	parsing := Design{}
	lines, parsing.Comment, remain = parseCommentBlock(remain)
	if needsLine && lines == 0 {
		return nil, suffix, SyntaxError{Remainder: remain, Msg: "expecting new line"}
	}
	postComment := remain
	if parsing.Name, remain, err = parseName(remain); err != nil {
		return nil, suffix, err
	}
	if parsing.Factor, remain, err = parseFactorBlock(remain); err != nil {
		return nil, suffix, err
	}
	inline := Inline{}
	if parsing.Returns, inline, remain, err = parseAssembly(&asmCtx{}, remain); err != nil {
		return nil, suffix, SyntaxError{Remainder: remain, Msg: "expecting a function block"}
	} else if len(parsing.Returns.Sign.Path) > 0 || parsing.Returns.Type != "{}" {
		return nil, suffix, SyntaxError{Remainder: remain, Msg: "expecting a function block"}
	}
	// attach all inline terms (arising from series) to function body
	parsing.Returns.Term = append(parsing.Returns.Term, inline.Series...)
	parsing.Lex = lex.LexUnion(postComment[:len(postComment)-len(remain)]...)
	return append([]Design{parsing}, inline.Design...), remain, nil
}

func parseFactorBlock(suffix []lex.Lex) (factor []Factor, remain []lex.Lex, err error) {
	var bra Bracket
	if bra, remain, err = parseBracket(suffix); err != nil {
		return nil, suffix, err
	}
	if bra.Type() != "()" {
		return nil, suffix, SyntaxError{Remainder: remain, Msg: "factor block uses round brackets"}
	}
	if factor, err = parseFactorBody(bra.Body); err != nil {
		return nil, suffix, err
	}
	return factor, suffix[1:], nil
}

func parseFactorBody(suffix []lex.Lex) (factor []Factor, err error) {
	monadic := false
	for len(suffix) > 0 {
		var f Factor
		f, suffix, err = parseFactor(len(factor) > 0, suffix)
		if err != nil {
			break
		}
		if f.Monadic {
			if monadic {
				return nil, SyntaxError{Remainder: suffix, Msg: "multiple monadic factors"}
			} else {
				monadic = true
			}
		}
		factor = append(factor, f)
	}
	_, _, suffix = parseCommentBlock(suffix) // parse remaining comments (and ignore)
	if len(suffix) > 0 {                     // check that all of remain is exhausted
		return nil, SyntaxError{Remainder: suffix, Msg: "not a factor"}
	}
	return factor, nil
}

func parseFactor(needsLine bool, suffix []lex.Lex) (factor Factor, remain []lex.Lex, err error) {
	remain = suffix
	var lines int
	lines, factor.Comment, remain = parseCommentBlock(remain)
	if needsLine && lines == 0 {
		return Factor{}, suffix, SyntaxError{Remainder: suffix, Msg: "expecting new line"}
	}
	postComment := remain
	if factor.Name, remain, err = parseName(remain); err != nil {
		return Factor{}, suffix, err
	}
	if remain, err = matchPunc("?", remain); err == nil {
		factor.Monadic = true
	}
	factor.Lex = lex.LexUnion(postComment[:len(postComment)-len(remain)]...)
	return factor, remain, nil
}

func DesignToSyntax(f ...Design) []Syntax {
	s := make([]Syntax, len(f))
	for i := range f {
		s[i] = f[i]
	}
	return s
}
