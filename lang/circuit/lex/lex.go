package lex

import (
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
	"text/scanner"
)

type Lex interface {
	Region
}

func LexifyFile(fileName string) ([]Lex, error) { // []Token
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return LexifyString(fileName, string(buf))
}

func LexifyString(fileName, src string) ([]Lex, error) { // []Token
	return LexifyReader(fileName, strings.NewReader(src))
}

func LexifyReader(fileName string, r io.Reader) ([]Lex, error) { // []Token
	var tok []Lex
	lex := NewLexer(fileName, r)
	for lex.Next() {
		t := lex.Token()
		if _, isWhite := t.Char.(White); !isWhite { // skip whitespaces
			tok = append(tok, t)
		}
	}
	if lex.Error() != nil {
		return nil, lex.Error()
	}
	return tok, nil
}

type Lexer struct {
	file string
	scan scanner.Scanner
	at   Position
	tok  Token
	err  error
}

func NewLexer(fileName string, r io.Reader) *Lexer {
	l := &Lexer{file: fileName}
	l.scan.Init(r)
	l.scan.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars |
		scanner.ScanStrings | scanner.ScanRawStrings | scanner.ScanComments
	l.scan.Whitespace = 0
	l.at = Position{Filename: l.file, Offset: 0, Line: 1, Column: 1}
	return l
}

func (l *Lexer) Error() error { return l.err }

func (l *Lexer) Token() Token { return l.tok }

func (l *Lexer) tokenText() string { return l.scan.TokenText() }

func (l *Lexer) newToken(char Char) (Token, Position) {
	start := l.at
	end := Position(l.scan.Pos())
	end.Filename = l.file
	return Token{
		Char: char,
		Text: l.tokenText(),
		Lex:  &StartEndRegion{Start: start, End: end},
	}, end
}

func (l *Lexer) Next() bool {
	for {
		switch r := l.scan.Scan(); r {
		case '\t', '\r', ' ': // update position
			white := []byte{byte(r)}
			for {
				if q := l.scan.Peek(); q == '\t' || q == '\r' || q == ' ' {
					white = append(white, byte(l.scan.Scan()))
				} else {
					break
				}
			}
			l.tok, l.at = l.newToken(White{string(white)})
			return true
		case '\n', ',', ';':
			l.tok, l.at = l.newToken(Line{})
			return true
		case ':', '{', '}', '[', ']', '.', '(', ')', '-', '?':
			l.tok, l.at = l.newToken(Punc{l.tokenText()})
			return true
		case scanner.Ident:
			l.tok, l.at = l.newToken(Name{l.tokenText()})
			return true
		case scanner.Int:
			i, err := strconv.Atoi(l.tokenText())
			if err != nil {
				panic(err)
			}
			l.tok, l.at = l.newToken(LexInteger{Int64: int64(i)})
			return true
		case scanner.Float:
			f, err := strconv.ParseFloat(l.tokenText(), 64)
			if err != nil {
				panic(err)
			}
			l.tok, l.at = l.newToken(LexFloat{Float64: f})
			return true
		case scanner.RawString:
			l.tok, l.at = l.newToken(LexString{String: l.tokenText()})
			return true
		case scanner.String:
			unquoted, err := strconv.Unquote(l.tokenText())
			if err != nil {
				l.err = err
				return false
			}
			l.tok, l.at = l.newToken(LexString{String: unquoted})
			return true
		case scanner.Comment:
			l.tok, l.at = l.newToken(Comment{Uncomment(l.tokenText())})
			return true
		case scanner.EOF:
			return false
		case scanner.Char:
			l.tok, l.at = l.newToken(LexOther{String: l.tokenText()})
			return true
		default:
			l.err = fmt.Errorf("unrecognized token %q", l.tokenText())
			return false
		}
	}
}

func Uncomment(raw string) string {
	if strings.HasPrefix(raw, "//") {
		raw = raw[len("//"):]
	} else if strings.HasPrefix(raw, "/*") && strings.HasSuffix(raw, "*/") {
		raw = raw[len("/*") : len(raw)-len("*/")]
	}
	return raw
}
