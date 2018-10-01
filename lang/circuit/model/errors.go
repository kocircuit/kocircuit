//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package model

import (
	"fmt"
	"runtime"

	ko "github.com/kocircuit/kocircuit"
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
		ko.SanitizeKoCompilerSourcePath(e.File),
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
