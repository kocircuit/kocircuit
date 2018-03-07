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
	for len(remain) > 0 {
		tok, ok := remain[0].(Token)
		if !ok {
			break
		}
		switch t := tok.Char.(type) {
		case Comment:
			fmt.Fprintln(&w, t.String)
		case Line:
			lines++
			fmt.Fprintln(&w)
		default:
			return lines, w.String(), remain
		}
		remain = remain[1:]
	}
	return lines, w.String(), remain
}
