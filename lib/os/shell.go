package os

import (
	"fmt"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

const ShellBinary = "/bin/sh"

type GoWipeDir struct {
	Dir string `ko:"name=dir,monadic"`
}

func (wipe *GoWipeDir) Play(ctx *runtime.Context) bool {
	return wipe.Task().Play(ctx)
}

func (wipe *GoWipeDir) Task(after ...*GoTask) *GoTask {
	return &GoTask{
		Name:   fmt.Sprintf("wipeDir: %s", wipe.Dir),
		Binary: ShellBinary,
		Arg:    []string{"-c", fmt.Sprintf("rm -rf %s/*", wipe.Dir)},
		After:  after,
	}
}
