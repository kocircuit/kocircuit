package runtime

import (
	"context"
	"fmt"
	"log"
	"runtime"

	. "github.com/kocircuit/kocircuit"
)

func NewContext() *Context {
	return &Context{
		Source:  "",
		Context: context.Background(),
	}
}

func CompilerContext() *Context {
	return &Context{
		Source:  fmt.Sprintf("ko/%s", KoVersion),
		Context: context.Background(),
	}
}

// Context is runtime context passed to Go operator implementations.
type Context struct {
	Parent  *Context        `ko:"name=parent"`
	Source  string          `ko:"name=source"` // source location of this invocation
	Context context.Context `ko:"name=context"`
	Kill    <-chan struct{} `ko:"name=kill"` // closure on this channel is a kill signal
}

func errorfSkip(skip int, ctx string, format string, arg ...interface{}) string {
	_, file, line, _ := runtime.Caller(1 + skip)
	return fmt.Sprintf(
		"(%s:%d) (%s) %s",
		SanitizeKoCompilerSourcePath(file),
		line,
		ctx,
		fmt.Sprintf(format, arg...),
	)
}

func (ctx *Context) Printf(format string, arg ...interface{}) {
	if ctx == nil {
		log.Print(errorfSkip(1, "no-ctx", format, arg...))
	} else {
		log.Print(errorfSkip(1, ctx.Source, format, arg...))
	}
}

func (ctx *Context) Fatalf(format string, arg ...interface{}) {
	if ctx == nil {
		log.Fatal(errorfSkip(1, "no-ctx", format, arg...))
	} else {
		log.Fatal(errorfSkip(1, ctx.Source, format, arg...))
	}
}

func (ctx *Context) Panicf(format string, arg ...interface{}) {
	if ctx == nil {
		log.Panic(errorfSkip(1, "no-ctx", format, arg...))
	} else {
		log.Panic(errorfSkip(1, ctx.Source, format, arg...))
	}
}

// Fault represents a panic occurring while calling a subfunction from the returning function.
type Fault struct {
	Context *Context    `ko:"name=context"` // context passed to callee (which emitted error or panic)
	Panic   interface{} `ko:"name=panic"`   // set if call panicked
	GoStack []byte      `ko:"name=goStack"`
}

func (p *Fault) Error() string {
	return fmt.Sprintf("fault (%v): %s", p.Panic, string(p.GoStack))
}

// Recoverer captures the runtime result of a circuit step.
type Recoverer interface {
	Recover() (recovered interface{})
	Stack() []byte
	Context() *Context
}
