package symbol

import (
	"fmt"
	"strconv"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type typingCtx struct {
	Parent *typingCtx `ko:"name=parent"`
	Span   *Span      `ko:"name=span"`
	Walk   string     `ko:"name=walk"`
}

func (ctx *typingCtx) Refine(walk string) *typingCtx {
	return &typingCtx{Parent: ctx, Span: ctx.Span, Walk: walk}
}

func (ctx *typingCtx) RefineIndex(i int) *typingCtx {
	return ctx.Refine(strconv.Itoa(i))
}

func (ctx *typingCtx) Path() Path {
	if ctx == nil {
		return nil
	} else if ctx.Parent == nil {
		return Path{ctx.Walk}
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
