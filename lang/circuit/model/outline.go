package model

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Outline string

func (outline Outline) SheathID() *ID {
	return PtrID(StringID(string(outline)))
}

func (outline Outline) SheathLabel() *string {
	return PtrString(string(outline))
}

func (outline Outline) SheathString() *string {
	return PtrString(string(outline))
}

func RefineOutline(span *Span, outline string) *Span {
	return span.Refine(Outline(outline))
}
