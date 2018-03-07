package syntax

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
)

type Import struct {
	Path      string `ko:"name=path"`
	As        Ref    `ko:"name=as"`
	Comment   string `ko:"name=comment" ctx:"expand"`     // comments preceding import statement
	EndOfLine string `ko:"name=eol_comment" ctx:"expand"` // end-of-line comment following import statement
	Lex       Region `ko:"name=lex"`                      // lex tokens comprising the import statement and surrounding comments
}

func (imp Import) FilePath() string {
	return imp.Lex.FilePath()
}

func (imp Import) StartPosition() Position {
	return imp.Lex.StartPosition()
}

func (imp Import) EndPosition() Position {
	return imp.Lex.EndPosition()
}

func (imp Import) RegionString() string {
	return imp.Lex.RegionString()
}

func parseImport(needsLine bool, suffix []Lex) (parsed []Import, remain []Lex, err error) {
	var lines int
	parsing := Import{}
	if lines, parsing.Comment, remain = parseCommentBlock(suffix); needsLine && lines == 0 { // parse preceding comments
		return nil, suffix, SyntaxError{Remainder: suffix, Msg: "expecting new line"}
	}
	postComment := remain
	if parsing.Path, remain, err = parseImportPath(remain); err != nil {
		return nil, suffix, err
	}
	parsing.As, remain, _ = parseImportAs(remain)                         // optionally parse import alias
	parsing.EndOfLine, remain = parseComment(remain)                      // eat follow up comments
	parsing.Lex = LexUnion(postComment[:len(postComment)-len(remain)]...) // save parsed tokens
	return []Import{parsing}, remain, nil
}

// parseImportPath parses strings of the form:
//	import "a/b/c"
func parseImportPath(suffix []Lex) (parsed string, remain []Lex, err error) {
	if remain, err = matchImportKeyword(suffix); err != nil {
		return "", suffix, err
	}
	if parsed, remain, err = parseString(remain); err != nil {
		return "", suffix, err
	}
	return parsed, remain, nil
}

func matchImportKeyword(suffix []Lex) (remain []Lex, err error) {
	if remain, err = matchKeyword("import", suffix); err == nil {
		return remain, nil
	}
	if remain, err = matchKeyword("use", suffix); err == nil {
		return remain, nil
	}
	return suffix, SyntaxError{Remainder: suffix, Msg: "not import or use clause"}
}

func parseImportAs(suffix []Lex) (parsed Ref, remain []Lex, err error) {
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
