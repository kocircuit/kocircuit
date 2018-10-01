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

// ParseFileString parses the tokens of a file.
func ParseFileString(fileName, text string) (File, error) {
	lex, err := lex.LexifyString(fileName, text)
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

func parseFileBody(suffix []lex.Lex) (file File, err error) {
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

func parseImportOrDesign(needsLine bool, suffix []lex.Lex) (parsed []Syntax, remain []lex.Lex, err error) {
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
		Msg:       "syntax error",
		Cause:     []SyntaxError{err1.(SyntaxError), err2.(SyntaxError)},
	}
}
