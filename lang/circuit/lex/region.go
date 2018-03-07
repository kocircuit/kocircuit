package lex

import (
	"fmt"
	"text/scanner"
)

type Region interface {
	FilePath() string
	StartPosition() Position
	EndPosition() Position
	RegionString() string
}

func RegionStart(r Region) Region {
	return &StartEndRegion{
		Start: r.StartPosition(),
		End:   r.StartPosition(),
	}
}

func RegionEnd(r Region) Region {
	return &StartEndRegion{
		Start: r.EndPosition(),
		End:   r.EndPosition(),
	}
}

func LexUnion(monotone ...Lex) Lex {
	monotone = foldNilLex(monotone)
	return &StartEndRegion{
		Start: monotone[0].StartPosition(),
		End:   monotone[len(monotone)-1].EndPosition(),
	}
}

func foldNilLex(x []Lex) (y []Lex) {
	for _, x := range x {
		if x != nil {
			y = append(y, x)
		}
	}
	return
}

func RegionUnion(monotone ...Region) Region {
	monotone = foldNilRegion(monotone)
	return &StartEndRegion{
		Start: monotone[0].StartPosition(),
		End:   monotone[len(monotone)-1].EndPosition(),
	}
}

func foldNilRegion(x []Region) (y []Region) {
	for _, x := range x {
		if x != nil {
			y = append(y, x)
		}
	}
	return
}

type EmptyRegion struct{}

func (er EmptyRegion) FilePath() string { return "" }

func (er EmptyRegion) StartPosition() Position { return Position{} }

func (er EmptyRegion) EndPosition() Position { return Position{} }

func (er EmptyRegion) RegionString() string { return "" }

type StartEndRegion struct {
	Start Position `ko:"name=start"`
	End   Position `ko:"name=end"`
}

func (sr *StartEndRegion) FilePath() string { return sr.Start.Filename }

func (sr *StartEndRegion) StartPosition() Position { return sr.Start }

func (sr *StartEndRegion) EndPosition() Position { return sr.End }

func (sr *StartEndRegion) RegionString() string {
	return fmt.Sprintf(
		"%s:%d:%d",
		// "%s:%d:%d-%d:%d",
		sr.FilePath(),
		sr.Start.Line, sr.Start.Column,
		// sr.End.Line, sr.End.Column,
	)
}

// type Position struct {
//    Filename string // filename, if any
//    Offset   int    // byte offset, starting at 0
//    Line     int    // line number, starting at 1
//    Column   int    // column number, starting at 1 (character count per line)
// }
// A source position is represented by a Position value. A position is valid if Line > 0.
type Position scanner.Position
