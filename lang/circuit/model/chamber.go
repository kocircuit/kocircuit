package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Chamber string

func (chamber Chamber) SheathID() *ID {
	return PtrID(StringID(string(chamber)))
}

func (chamber Chamber) SheathLabel() *string {
	return PtrString(string(chamber))
}

func (chamber Chamber) SheathString() *string {
	return PtrString(fmt.Sprintf("chamber=%s", chamber))
}

func RefineChamber(span *Span, ch string) *Span {
	return span.Refine(Chamber(ch))
}

func ChamberPath(span *Span) Path {
	if span == nil {
		return nil
	} else if chamber, ok := span.Sheath.(Chamber); ok {
		return append(ChamberPath(span.Parent), string(chamber))
	} else {
		return ChamberPath(span.Parent)
	}
}
