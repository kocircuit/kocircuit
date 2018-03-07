package model

import (
	"fmt"
	"log"

	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Span struct {
	ID     ID     `ko:"name=id"`
	Parent *Span  `ko:"name=parent"`
	Sheath Sheath `ko:"name=sheath"`
}

var rootSpanID = StringID("root")

func NewSpan() *Span {
	return &Span{ID: rootSpanID}
}

func (span *Span) Splay() Tree {
	return Quote{span.SourceLine()}
}

// Sheaths are *Func, *Step, MacroSheath, Outline, Chamber, *AssignCache, *GoWeavingCtx
type Sheath interface {
	SheathID() *ID         // unique identifier
	SheathLabel() *string  // empty or a short identifier string (not necessarily uniquely-identifying)
	SheathString() *string // human-readable line ot text
}

func (span *Span) SpanID() ID {
	if span == nil {
		return rootSpanID
	} else {
		return span.ID
	}
}

func (span *Span) LabelPath() Path {
	if span == nil {
		return nil
	} else if span.Sheath != nil && span.Sheath.SheathLabel() != nil {
		return append(span.Parent.LabelPath(), *span.Sheath.SheathLabel())
	} else {
		return span.Parent.LabelPath()
	}
}

func (span *Span) Refine(sheath Sheath) *Span {
	var id ID
	if sheathID := sheath.SheathID(); sheathID != nil {
		id = Blend(span.ID, *sheathID)
	} else {
		id = span.ID
	}
	return &Span{
		ID:     id,
		Parent: span,
		Sheath: sheath,
	}
}

func RefineFunc(span *Span, f *Func) *Span {
	return span.Refine(f)
}

func RefineStep(span *Span, step *Step) *Span {
	return span.Refine(step)
}

func NearestFunc(span *Span) *Func {
	if span == nil {
		return nil
	} else if fu, _ := span.Sheath.(*Func); fu != nil {
		return fu
	} else {
		return NearestFunc(span.Parent)
	}
}

func NearestStep(span *Span) *Step {
	if span == nil {
		return nil
	} else if step, _ := span.Sheath.(*Step); step != nil {
		return step
	} else {
		return NearestStep(span.Parent)
	}
}

func NearestSyntax(span *Span) Syntax {
	if span == nil {
		return EmptySyntax{}
	} else if syntax, _ := span.Sheath.(Syntax); syntax != nil {
		return syntax
	} else {
		return NearestSyntax(span.Parent)
	}
}

func (span *Span) CommentLine() string {
	return span.SourceLine()
}

func (span *Span) SourceLine() string {
	return fmt.Sprintf("span:%s, %s", span.SpanID().String(), NearestSyntax(span).RegionString())
}

func (span *Span) String() string {
	return fmt.Sprintf("(%v)", span.SourceLine())
}

func (span *Span) ErrorfSkip(skip int, cause error, format string, arg ...interface{}) error {
	return NewSpanErrorf(skip+1, span, cause, format, arg...)
}

func (span *Span) Errorf(cause error, format string, arg ...interface{}) error {
	return span.ErrorfSkip(1, cause, format, arg...)
}

func (span *Span) Fatalf(cause error, format string, arg ...interface{}) {
	log.Fatal(span.ErrorfSkip(2, cause, format, arg...).Error())
}

func (span *Span) Panicf(cause error, format string, arg ...interface{}) {
	log.Panic(span.ErrorfSkip(2, cause, format, arg...).Error())
}

func (span *Span) Printf(format string, arg ...interface{}) {
	log.Print(span.ErrorfSkip(2, nil, format, arg...).Error())
}

func (span *Span) Print(arg ...interface{}) {
	log.Print(span.ErrorfSkip(2, nil, fmt.Sprint(arg...)).Error())
}
