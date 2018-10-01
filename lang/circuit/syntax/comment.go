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
	"bytes"
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/lex"
)

func parseComment(suffix []lex.Lex) (comment string, remain []lex.Lex) {
	if len(suffix) == 0 {
		return "", suffix
	}
	tok, ok := suffix[0].(lex.Token)
	if !ok {
		return "", suffix
	}
	c, ok := tok.Char.(lex.Comment)
	if !ok {
		return "", suffix
	}
	return c.String, suffix[1:]
}

func parseCommentBlock(suffix []lex.Lex) (lines int, comment string, remain []lex.Lex) {
	var w bytes.Buffer
	remain = suffix
	linesInSeq := 2 // causes leading lines to be ignored (below)
	for len(remain) > 0 {
		tok, ok := remain[0].(lex.Token)
		if !ok {
			break
		}
		switch t := tok.Char.(type) {
		case lex.Comment:
			linesInSeq = 0
			fmt.Fprint(&w, t.String)
		case lex.Line:
			lines++
			linesInSeq++
			if linesInSeq < 2 {
				fmt.Fprintln(&w)
			}
		default:
			return lines, w.String(), remain
		}
		remain = remain[1:]
	}
	return lines, w.String(), remain
}
