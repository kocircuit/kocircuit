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

type Syntax interface {
	lex.Region
}

type EmptySyntax struct {
	lex.EmptyRegion `ko:"name=emptyRegion"`
}

type Assembly struct {
	Sign    Ref     `ko:"name=ref"`
	Type    string  `ko:"name=type"` // "{}" or "[]" or "()"
	Term    []Term  `ko:"name=term"`
	Bracket lex.Lex `ko:"name=bracket"`
}

func (asm Assembly) Syntax() lex.Region {
	return lex.RegionUnion(asm.Sign, asm.Bracket)
}

func (asm Assembly) FilePath() string {
	return asm.Syntax().FilePath()
}

func (asm Assembly) StartPosition() lex.Position {
	return asm.Syntax().StartPosition()
}

func (asm Assembly) EndPosition() lex.Position {
	return asm.Syntax().EndPosition()
}

func (asm Assembly) RegionString() string {
	return asm.Syntax().RegionString()
}

type Term struct {
	Comment string `ko:"name=comment" ctx:"expand"`
	Label   Ref    `ko:"name=label"`
	Hitch   Syntax `ko:"name=hitch"` // Assembly, Ref, Literal
}

func (term Term) FilePath() string {
	return term.Hitch.FilePath()
}

func (term Term) StartPosition() lex.Position {
	return term.Hitch.StartPosition()
}

func (term Term) EndPosition() lex.Position {
	return term.Hitch.EndPosition()
}

func (term Term) RegionString() string {
	return term.Hitch.RegionString()
}

type asmCtx struct {
	SeriesDepth int `ko:"name=series_depth"`
	SeriesLabel Ref `ko:"name=series_label"`
}

func (ctx *asmCtx) Append(ref Ref) *asmCtx {
	return &asmCtx{
		SeriesDepth: 0,
		SeriesLabel: Ref{
			Lex: lex.RegionUnion(ctx.SeriesLabel.Lex, ref.Lex),
			Path: append(
				ctx.SeriesLabel.Path,
				fmt.Sprintf("%d_%s", ctx.SeriesDepth, ref.Join("_")),
			),
		},
	}
}

func (ctx *asmCtx) Inc() *asmCtx {
	return &asmCtx{SeriesLabel: ctx.SeriesLabel, SeriesDepth: ctx.SeriesDepth + 1}
}

func (ctx *asmCtx) InlineSeriesLabelRef() Ref {
	return Ref{
		Lex: ctx.SeriesLabel,
		Path: []string{
			fmt.Sprintf("0_inline_%s_%d", strings.Join(ctx.SeriesLabel.Path, "_"), ctx.SeriesDepth),
		},
	}
}

func parseAssembly(ctx *asmCtx, suffix []lex.Lex) (parsed Assembly, inline Inline, remain []lex.Lex, err error) {
	if len(suffix) == 0 {
		return Assembly{}, Inline{}, suffix, SyntaxError{Remainder: remain, Msg: "end of source"}
	}
	remain = suffix
	var emptySign bool
	if parsed.Sign, remain, err = ParseRef(remain); err != nil { // assembly name is optional
		emptySign = true
	}
	var bra Bracket
	if bra, remain, err = parseBracket(remain); err != nil {
		return Assembly{}, Inline{}, suffix, err
	}
	if bra.Type() != "{}" && bra.Type() != "[]" && bra.Type() != "()" {
		return Assembly{}, Inline{}, suffix, SyntaxError{Remainder: remain, Msg: "bracket is not curly or square"}
	}
	parsed.Type = bra.Type()
	if parsed.Term, inline, err = parseTermHitch(ctx, bra.Body); err != nil {
		return Assembly{}, Inline{}, suffix, err
	}
	parsed.Bracket = bra
	if emptySign {
		parsed.Sign = Ref{Lex: bra}
	}
	return parsed, inline, remain, nil
}

func parseTermHitch(ctx *asmCtx, suffix []lex.Lex) (term []Term, inline Inline, err error) {
	if term, inline, err = parseColonTermHitch(ctx, suffix); err == nil {
		return term, inline, nil
	}
	return parseNoColonTermHitch(ctx, suffix)
}

func parseColonTermHitch(ctx *asmCtx, suffix []lex.Lex) (term []Term, inline Inline, err error) {
	for len(suffix) > 0 {
		t, termInline := Term{}, Inline{}
		if t, termInline, suffix, err = parseColonTerm(ctx, len(term) > 0, suffix); err != nil {
			break
		}
		term, inline = append(term, t), inline.Union(termInline)
	}
	_, _, suffix = parseCommentBlock(suffix)
	if len(suffix) > 0 {
		return nil, Inline{}, SyntaxError{Remainder: suffix, Msg: "no terms ahead"}
	}
	return term, inline, nil
}

func parseColonTerm(ctx *asmCtx, needsLine bool, suffix []lex.Lex) (parsed Term, inline Inline, remain []lex.Lex, err error) {
	remain = suffix
	var lines int
	lines, parsed.Comment, remain = parseCommentBlock(remain)
	if needsLine && lines == 0 {
		return Term{}, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "term not preceded by new line"}
	}
	parsed.Label, remain, _ = parseName(remain) // name is optional
	if remain, err = matchPunc(":", remain); err != nil {
		return Term{}, Inline{}, suffix, err
	}
	if parsed.Hitch, inline, remain, err = parseHitchSeries(ctx.Append(parsed.Label), remain); err != nil {
		return Term{}, Inline{}, suffix, err
	}
	return parsed, inline, remain, nil
}

func parseNoColonTermHitch(ctx *asmCtx, suffix []lex.Lex) (term []Term, inline Inline, err error) {
	for len(suffix) > 0 {
		t, termInline := Term{}, Inline{}
		if t, termInline, suffix, err = parseNoColonTerm(ctx, len(term) > 0, suffix); err != nil {
			break
		}
		ctx = ctx.Inc()
		term, inline = append(term, t), inline.Union(termInline)
	}
	_, _, suffix = parseCommentBlock(suffix)
	if len(suffix) > 0 {
		return nil, Inline{}, SyntaxError{Remainder: suffix, Msg: "no terms ahead"}
	}
	return term, inline, nil
}

const NoLabel = ""

func parseNoColonTerm(ctx *asmCtx, needsLine bool, suffix []lex.Lex) (parsed Term, inline Inline, remain []lex.Lex, err error) {
	if len(suffix) == 0 {
		return Term{}, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "unexpected end of source"}
	}
	remain = suffix
	var lines int
	lines, parsed.Comment, remain = parseCommentBlock(remain)
	if needsLine && lines == 0 {
		return Term{}, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "term not preceded by new line"}
	}
	if len(remain) == 0 {
		return Term{}, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "unexpected end of source"}
	}
	parsed.Label = Ref{
		Path: []string{NoLabel},
		Lex:  lex.RegionStart(remain[0]),
	}
	if parsed.Hitch, inline, remain, err = parseHitchSeries(ctx.Append(parsed.Label), remain); err != nil {
		return Term{}, Inline{}, suffix, err
	}
	return parsed, inline, remain, nil
}

func parseHitchSeries(ctx *asmCtx, suffix []lex.Lex) (parsed Syntax, inline Inline, remain []lex.Lex, err error) {
	if len(suffix) == 0 {
		return nil, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "unexpected end of source"}
	}
	ctx = ctx.Append(
		Ref{
			Path: []string{"series"},
			Lex:  lex.RegionStart(suffix[0]),
		},
	)
	remain = suffix
	var hitch Syntax
	var hitchInline Inline
	if hitch, hitchInline, remain, err = parseHitch(ctx, remain); err != nil {
		return nil, Inline{}, suffix, err
	}
	inline = inline.Union(hitchInline)
	series := []Term{{
		Label: ctx.InlineSeriesLabelRef(),
		Hitch: hitch,
	}}
	pastCtx, presentCtx := ctx, ctx.Inc() // iter 0
	for len(remain) > 0 {
		tailTerm, tailInline, tailRemain, err := parseAssemblyTail(pastCtx, presentCtx, remain)
		if err != nil {
			if tailTerm, tailInline, tailRemain, err = parseSelectTail(pastCtx, presentCtx, remain); err != nil {
				break
			}
		}
		series = append(series, tailTerm)
		inline, remain = inline.Union(tailInline), tailRemain
		pastCtx, presentCtx = presentCtx, presentCtx.Inc() // iter++
	}
	if len(series) == 1 {
		return hitch, inline, remain, nil
	}
	inline.Series = append(inline.Series, series...)
	return pastCtx.InlineSeriesLabelRef(), inline, remain, nil
}

func parseAssemblyTail(pastCtx, presentCtx *asmCtx, lex []lex.Lex) (parsed Term, inline Inline, remain []lex.Lex, err error) {
	remain = lex
	asm, asmInline, asmRemain, err := parseAssembly(presentCtx, remain)
	if err != nil {
		return Term{}, Inline{}, lex, err
	}
	if len(asm.Sign.Path) != 0 || asm.Type == "{}" {
		return Term{}, Inline{}, lex, fmt.Errorf("not a tail invocation () or augmentation []")
	}
	asm.Sign = pastCtx.InlineSeriesLabelRef()
	return Term{
		Label: presentCtx.InlineSeriesLabelRef(),
		Hitch: asm,
	}, asmInline, asmRemain, nil
}

func parseSelectTail(pastCtx, presentCtx *asmCtx, l []lex.Lex) (Term, Inline, []lex.Lex, error) {
	if len(l) < 2 {
		return Term{}, Inline{}, l, fmt.Errorf("not a tail selection")
	}
	if !isPunc(".", l[0]) {
		return Term{}, Inline{}, l, fmt.Errorf("tail selection starts with a \".\"")
	}
	if ref, remain, err := ParseRef(l[1:]); err != nil {
		return Term{}, Inline{}, l, fmt.Errorf("not a tail selection (%v)", err)
	} else {
		return Term{
			Label: presentCtx.InlineSeriesLabelRef(),
			Hitch: Ref{
				Lex:  lex.LexUnion(l[:len(l)-len(remain)]...),
				Path: append([]string{pastCtx.InlineSeriesLabelRef().Path[0]}, ref.Path...),
			},
		}, Inline{}, remain, nil
	}
}

func parseHitch(ctx *asmCtx, suffix []lex.Lex) (parsed Syntax, inline Inline, remain []lex.Lex, err error) {
	remain = suffix
	funcBag := []Design{}
	if funcBag, remain, err = parseDesign(false, suffix); err == nil { // if hitch is an inline function definition
		return funcBag[0].Name, // return a syntactic package-level reference to the inline function
			Inline{Design: funcBag},
			remain, nil
	}
	if parsed, remain, err = parseLiteral(suffix); err == nil {
		return
	}
	if parsed, inline, remain, err = parseAssembly(ctx, suffix); err == nil {
		return
	}
	if parsed, remain, err = ParseRef(suffix); err == nil {
		return
	}
	return nil, Inline{}, suffix, SyntaxError{Remainder: suffix, Msg: "expecting literal, reference or assembly"}
}
