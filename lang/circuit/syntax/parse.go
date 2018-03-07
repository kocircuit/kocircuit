package syntax

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
)

// ParseFileString parses the tokens of a file.
func ParseFileString(fileName, text string) (File, error) {
	lex, err := LexifyString(fileName, text)
	if err != nil {
		return File{}, err
	}
	if lex, err = FoldBracket(lex); err != nil {
		return File{}, err
	}
	file, err := parseFileBody(lex)
	if err != nil {
		return File{}, err
	}
	file.Path = fileName
	return file, nil
}

type File struct {
	Path   string   `ko:"name=path"`
	Import []Import `ko:"name=import"`
	Design []Design `ko:"name=design"`
}

func parseFileBody(suffix []Lex) (file File, err error) {
	var n int
	for len(suffix) > 0 {
		var stmt []Syntax
		if stmt, suffix, err = parseImportOrDesign(n > 0, suffix); err != nil {
			break
		}
		for _, stmt := range stmt {
			switch t := stmt.(type) {
			case Import:
				file.Import = append(file.Import, t)
			case Design:
				file.Design = append(file.Design, t)
			}
		}
	}
	if _, _, suffix = parseCommentBlock(suffix); len(suffix) > 0 {
		return File{}, SyntaxError{Remainder: suffix, Msg: err.Error()}
	}
	return file, nil
}

func parseImportOrDesign(needsLine bool, suffix []Lex) (parsed []Syntax, remain []Lex, err error) {
	remain = suffix
	var err1, err2 error
	var importBag []Import
	if importBag, remain, err1 = parseImport(needsLine, suffix); err1 == nil {
		return ImportToSyntax(importBag...), remain, nil
	}
	var funcBag []Design
	if funcBag, remain, err2 = parseDesign(needsLine, suffix); err2 == nil {
		return DesignToSyntax(funcBag...), remain, nil
	}
	return nil, suffix, SyntaxError{
		Remainder: suffix,
		Msg:       "expecting import or design",
		Cause:     []SyntaxError{err1.(SyntaxError), err2.(SyntaxError)},
	}
}
