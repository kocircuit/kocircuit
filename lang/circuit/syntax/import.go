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

type Import struct {
	Path      string     `ko:"name=path"`
	As        Ref        `ko:"name=as"`
	Comment   string     `ko:"name=comment" ctx:"expand"`     // comments preceding import statement
	EndOfLine string     `ko:"name=eol_comment" ctx:"expand"` // end-of-line comment following import statement
	Lex       lex.Region `ko:"name=lex"`                      // lex tokens comprising the import statement and surrounding comments
}

func (imp Import) FilePath() string {
	return imp.Lex.FilePath()
}

func (imp Import) StartPosition() lex.Position {
	return imp.Lex.StartPosition()
}

func (imp Import) EndPosition() lex.Position {
	return imp.Lex.EndPosition()
}

func (imp Import) RegionString() string {
	return imp.Lex.RegionString()
}

func parseImport(needsLine bool, suffix []lex.Lex) (parsed []Import, remain []lex.Lex, err error) {
	var lines int
	parsing := Import{}
	if lines, parsing.Comment, remain = parseCommentBlock(suffix); needsLine && lines == 0 { // parse preceding comments
		return nil, suffix, SyntaxError{Remainder: suffix, Msg: "expecting new line"}
	}
	postComment := remain
	if parsing.Path, remain, err = parseImportPath(remain); err != nil {
		return nil, suffix, err
	}
	parsing.As, remain, _ = parseImportAs(remain)                             // optionally parse import alias
	parsing.EndOfLine, remain = parseComment(remain)                          // eat follow up comments
	parsing.Lex = lex.LexUnion(postComment[:len(postComment)-len(remain)]...) // save parsed tokens
	return []Import{parsing}, remain, nil
}

// parseImportPath parses strings of the form:
//	import "a/b/c"
func parseImportPath(suffix []lex.Lex) (parsed string, remain []lex.Lex, err error) {
	if remain, err = matchImportKeyword(suffix); err != nil {
		return "", suffix, err
	}
	if parsed, remain, err = parseString(remain); err != nil {
		return "", suffix, err
	}
	return parsed, remain, nil
}

func matchImportKeyword(suffix []lex.Lex) (remain []lex.Lex, err error) {
	if remain, err = matchKeyword("import", suffix); err == nil {
		return remain, nil
	}
	if remain, err = matchKeyword("use", suffix); err == nil {
		return remain, nil
	}
	return suffix, SyntaxError{Remainder: suffix, Msg: "not import or use clause"}
}

func parseImportAs(suffix []lex.Lex) (parsed Ref, remain []lex.Lex, err error) {
	remain = suffix
	if remain, err = matchKeyword("as", remain); err != nil {
		return Ref{}, suffix, err
	}
	if parsed, remain, err = parseName(remain); err != nil {
		return Ref{}, suffix, err
	}
	return parsed, remain, nil
}

func ImportToSyntax(m ...Import) []Syntax {
	s := make([]Syntax, len(m))
	for i := range m {
		s[i] = m[i]
	}
	return s
}
