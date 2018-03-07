package os

import (
	"os/exec"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterEvalGate(new(GoTask))
	RegisterGoGate(new(GoTask))
}

type GoTask struct {
	Name   string    `ko:"name=name"`
	Binary string    `ko:"name=binary"`
	Arg    []string  `ko:"name=arg"`
	Env    []string  `ko:"name=env"`
	Dir    *string   `ko:"name=dir"`
	After  []*GoTask `ko:"name=after"`
}

func (t *GoTask) Play(ctx *runtime.Context) bool {
	if out, err := t.RunWithOutput(); err != nil {
		ctx.Printf("task %s exited unsuccessfully (%v) with output:\n%s", Sprint(t), err, out)
		return false
	} else {
		return true
	}
}

func (t *GoTask) RunWithOutput() (string, error) {
	cmd := &exec.Cmd{
		Path: t.Binary,
		Args: append([]string{t.Binary}, t.Arg...),
		Dir:  OptString(t.Dir, ""),
		Env:  t.Env,
	}
	std, err := cmd.CombinedOutput()
	return string(std), err
}
