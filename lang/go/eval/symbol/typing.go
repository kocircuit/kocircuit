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

package symbol

import (
	"fmt"
	"strconv"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

type typingCtx struct {
	Parent *typingCtx  `ko:"name=parent"`
	Span   *model.Span `ko:"name=span"`
	Walk   string      `ko:"name=walk"`
}

func (ctx *typingCtx) Refine(walk string) *typingCtx {
	return &typingCtx{Parent: ctx, Span: ctx.Span, Walk: walk}
}

func (ctx *typingCtx) RefineIndex(i int) *typingCtx {
	return ctx.Refine(strconv.Itoa(i))
}

func (ctx *typingCtx) Path() model.Path {
	if ctx == nil {
		return nil
	} else if ctx.Parent == nil {
		return model.Path{ctx.Walk}
	} else {
		return append(ctx.Parent.Path(), ctx.Walk)
	}
}

func (ctx *typingCtx) Errorf(cause error, format string, arg ...interface{}) error {
	return ctx.Span.ErrorfSkip(
		2, cause,
		fmt.Sprintf("%v: %s", ctx.Path(), fmt.Sprintf(format, arg...)),
	)
}
