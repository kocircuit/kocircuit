package shell

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	koos "github.com/kocircuit/kocircuit/lib/os"
)

func init() {
	RegisterEvalGate(new(GoShellCopy))
}

type GoShellCopy struct {
	Source string `ko:"name=source"` // source path (file or directory)
	Sink   string `ko:"name=sink"`   // destination path (not destination directory)
}

func (g *GoShellCopy) Play(ctx *runtime.Context) bool {
	return g.Task().Play(ctx)
}

func (g *GoShellCopy) Task(after ...*koos.GoTask) *koos.GoTask {
	return &koos.GoTask{
		Name:   fmt.Sprintf("shell.Copy(source: %q, sink: %q)", g.Source, g.Sink),
		Binary: koos.ShellBinary,
		Arg:    []string{"-c", fmt.Sprintf("cp %q %q", g.Source, g.Sink)},
		After:  after,
	}
}
