package syntax

import (
	"bytes"
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
)

func parseComment(suffix []Lex) (comment string, remain []Lex) {
	if len(suffix) == 0 {
		return "", suffix
	}
	tok, ok := suffix[0].(Token)
	if !ok {
		return "", suffix
	}
	c, ok := tok.Char.(Comment)
	if !ok {
		return "", suffix
	}
	return c.String, suffix[1:]
}

func parseCommentBlock(suffix []Lex) (lines int, comment string, remain []Lex) {
	var w bytes.Buffer
	remain = suffix
	linesInSeq := 2 // causes leading lines to be ignored (below)
	for len(remain) > 0 {
		tok, ok := remain[0].(Token)
		if !ok {
			break
		}
		switch t := tok.Char.(type) {
		case Comment:
			linesInSeq = 0
			fmt.Fprint(&w, t.String)
		case Line:
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
