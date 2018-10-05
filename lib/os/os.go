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
	"os"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Mkdir", new(goMkdir))
	go_eval.RegisterNamedEvalGate("Env", new(goEnv))
	go_eval.RegisterNamedEvalGate("TempDir", new(goTempDir))
	go_eval.RegisterNamedEvalGate("Exit", new(goExit))
}

type goMkdir struct {
	Path string `ko:"name=path,monadic"`
}

func (mkdir *goMkdir) Play(ctx *runtime.Context) error {
	return os.MkdirAll(mkdir.Path, 0755)
}

func (*goMkdir) Doc() string  { return "Mkdir(path?) creates the given directory and all its parents" }
func (*goMkdir) Help() string { return "Mkdir(path?)" }

type goEnv struct {
	Name string `ko:"name=name,monadic"`
}

func (env *goEnv) Play(ctx *runtime.Context) string {
	return os.Getenv(env.Name)
}

func (*goEnv) Doc() string {
	return "Env(name?) returns the value of the environment variable with given name"
}
func (*goEnv) Help() string { return "Env(name?)" }

type goTempDir struct{}

func (goTempDir) Play(ctx *runtime.Context) string {
	return os.TempDir()
}

func (*goTempDir) Doc() string  { return "TempDir() returns the full path of the temporary directory" }
func (*goTempDir) Help() string { return "TempDir()" }

type goExit struct {
	ExitCode int `ko:"name=code,monadic"`
}

func (g *goExit) Play(ctx *runtime.Context) error {
	os.Exit(g.ExitCode)
	return nil
}

func (*goExit) Doc() string  { return "Exit(code?) terminates the process with given exit code" }
func (*goExit) Help() string { return "Exit(code?)" }
