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

package os

import (
	"os/exec"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Task", new(GoTask))
}

// GoTask implements the running on an external command.
type GoTask struct {
	Name   string    `ko:"name=name"`
	Binary string    `ko:"name=binary"`
	Arg    []string  `ko:"name=arg"`
	Env    []string  `ko:"name=env"`
	Dir    *string   `ko:"name=dir"`
	After  []*GoTask `ko:"name=after"`
}

// Play invokes the external command.
func (t *GoTask) Play(ctx *runtime.Context) bool {
	if out, err := t.RunWithOutput(); err != nil {
		ctx.Printf("task %s exited unsuccessfully (%v) with output:\n%s", tree.Sprint(t), err, out)
		return false
	}
	return true
}

// RunWithOutput invokes the external command, waits until it is terminated
// and returns its output.
func (t *GoTask) RunWithOutput() (string, error) {
	cmd := &exec.Cmd{
		Path: t.Binary,
		Args: append([]string{t.Binary}, t.Arg...),
		Dir:  util.OptString(t.Dir, ""),
		Env:  t.Env,
	}
	std, err := cmd.CombinedOutput()
	return string(std), err
}

func (*GoTask) Doc() string {
	return "Task(name, binary, arg, env, dir, after) runs an external process"
}
func (*GoTask) Help() string { return "Task(name, binary, arg, env, dir, after)" }
