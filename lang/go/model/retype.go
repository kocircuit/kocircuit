package model

import (
	"sort"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type GoRetyper interface {
	Type() GoType
}

type GoTypeRetyper struct {
	Type_ GoType `ko:"name=type"`
}

func (retyper *GoTypeRetyper) Type() GoType {
	return retyper.Type_
}

type GoReserveRetyper struct {
	Solution GoRetyper `ko:"name=solution"`
}

func (retyper *GoReserveRetyper) Type() GoType {
	return retyper.Solution.Type()
}

type GoPtrRetyper struct {
	Elem GoRetyper `ko:"name=elem"`
}

func (retyper *GoPtrRetyper) Type() GoType {
	return NewGoPtr(retyper.Elem.Type())
}

type GoExpandedPtrRetyper struct {
	Elem GoRetyper `ko:"name=elem"`
}

func (retyper *GoExpandedPtrRetyper) Type() GoType {
	return NewGoPtr(retyper.Elem.Type())
}

type GoArrayRetyper struct {
	Len  int       `ko:"name=len"`
	Elem GoRetyper `ko:"name=elem"`
}

func (retyper *GoArrayRetyper) Type() GoType {
	return NewGoArray(retyper.Len, retyper.Elem.Type())
}

type GoSliceRetyper struct {
	Span *Span     `ko:"name=span"`
	Elem GoRetyper `ko:"name=elem"`
}

func (retyper *GoSliceRetyper) Type() GoType {
	return NewGoSlice(retyper.Elem.Type())
}

type GoStructRetyper struct {
	Span  *Span                   `ko:"name=span"`
	Field []*GoStructRetyperField `ko:"name=field"`
}

type GoStructRetyperField struct {
	Name    string    `ko:"name=name"` // ko name
	Retyper GoRetyper `ko:"name=retyper"`
	Monadic bool      `ko:"name=monadic"`
}

func (retyper *GoStructRetyper) Type() GoType {
	field := make([]*GoField, len(retyper.Field))
	for i, fieldRetyper := range retyper.Field {
		field[i] = BuildGoField(retyper.Span, fieldRetyper.Name, fieldRetyper.Retyper.Type(), fieldRetyper.Monadic)
	}
	return NewGoStruct(field...)
}

func SortGoRetyperField(field []*GoStructRetyperField) []*GoStructRetyperField {
	sort.Sort(SortGoStructRetyperField(field))
	return field
}

type SortGoStructRetyperField []*GoStructRetyperField

func (s SortGoStructRetyperField) Len() int { return len(s) }

func (s SortGoStructRetyperField) Less(i, j int) bool { return s[i].Name < s[j].Name }

func (s SortGoStructRetyperField) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
