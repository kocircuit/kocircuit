// Copyright © 2017 Aljabr, Inc.
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

package cmd

import (
	"log"
	"path"

	"github.com/spf13/cobra"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	"github.com/kocircuit/kocircuit/lang/go/sys"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

// weaveCmd represents the weave command
var weaveCmd = &cobra.Command{
	Use:   "weave",
	Short: "(EXP) Weave a Ko source program to Go source",
	Long: `Weaving parses a Ko program from sources,
models and simulates the program in memory to infer data types
and verify correctness. It generates a statically-typed
Go source program, implementing the execution of the original Ko program.`,
	Run: func(cmd *cobra.Command, args []string) {
		tools := newToolchain()
		if len(args) != 1 {
			log.Fatalf("ko weave expects a single argument in the form \"path/to/pkg/Func\"")
		}
		koPkg, koFunc := parsePkgFunc(args[0])
		b := &sys.CompileWeave{
			Repo:      path.Join(tools.GOPATH, "src"),
			Pkg:       koPkg,
			Func:      koFunc,
			Faculty:   GoFaculty(),
			Idiom:     GoIdiomRepo,
			Arg:       nil,
			Toolchain: tools,
			GoKoRoot:  flagKOGO,
			GoKoPkg:   path.Join(koPkg, koFunc),
			Show:      false, // show compiled ko functions
		}
		if result := b.Play(runtime.CompilerContext()); result.Error != nil {
			log.Fatalln(result.Error)
		}
	},
}

func init() {
	RootCmd.AddCommand(weaveCmd)
	initGoBasedCmd(weaveCmd)
}
