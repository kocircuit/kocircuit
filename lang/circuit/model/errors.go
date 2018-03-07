package model

import (
	"fmt"
	"runtime"

	. "github.com/kocircuit/kocircuit"
)

func NewSpanErrorf(skip int, span *Span, cause error, format string, arg ...interface{}) *Error {
	_, file, line, _ := runtime.Caller(1 + skip)
	return &Error{Span: span, File: file, Line: line, Msg: fmt.Sprintf(format, arg...), Cause: cause}
}

type Error struct {
	Span  *Span  `ko:"name=span"`
	File  string `ko:"name=file"` // file in compiler code, generating error
	Line  int    `ko:"name=line"` // line in compiler code, generating error
	Msg   string `ko:"name=msg"`
	Cause error  `ko:"name=cause"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s%s", e.Span.Trajectory(), e.errorNoTrajectory())
}

func (e *Error) errorNoTrajectory() string {
	header := fmt.Sprintf(
		"\n(span:%s, %s:%d, %s) %s",
		e.Span.SpanID().String(),
		SanitizeKoCompilerSourcePath(e.File),
		e.Line,
		NearestSyntax(e.Span).RegionString(),
		e.Msg,
	)
	if e.Cause == nil {
		return header
	} else {
		switch u := e.Cause.(type) {
		case *Error:
			return fmt.Sprintf("%s%s", header, u.errorNoTrajectory())
		default:
			return fmt.Sprintf("%s\n%s", header, u.Error())
		}
	}
}
