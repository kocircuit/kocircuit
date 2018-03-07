package lex

import (
	"strconv"
)

// Token is a Lex.
type Token struct {
	Lex  `ko:"name=lex" ctx:"expand"`
	Text string `ko:"name=text"`
	Char Char   `ko:"name=char" ctx:"expand"`
}

type Char interface {
	Char() string
}

// Line is a character.
type Line struct {
	N int `ko:"name=n"`
}

func (line Line) Char() string { return "\n" }

// Punc is a character.
type Punc struct {
	String string `ko:"name=string"`
}

func (punc Punc) Char() string { return punc.String }

// Name is a character.
type Name struct {
	String string `ko:"name=string"`
}

func (name Name) Char() string { return name.String }

// Comment is a character.
type Comment struct {
	String string `ko:"name=string"`
}

func (comm Comment) Char() string { return comm.String }

type LexInteger struct {
	Int64 int64 `ko:"name=int64"`
}

func (integer LexInteger) Char() string { return strconv.Itoa(int(integer.Int64)) }

func (integer LexInteger) Negative() Char { return LexInteger{Int64: -integer.Int64} }

type LexFloat struct {
	Float64 float64 `ko:"name=float64"`
}

func (float LexFloat) Char() string { return strconv.FormatFloat(float.Float64, 'f', -1, 64) }

func (float LexFloat) Negative() Char { return LexFloat{Float64: -float.Float64} }

type LexString struct {
	String string `ko:"name=string"`
}

func (str LexString) Char() string { return str.String }

type White struct {
	String string `ko:"name=string"`
}

func (white White) Char() string { return white.String }

type LexOther struct {
	String string `ko:"name=string"`
}

func (str LexOther) Char() string { return str.String }
