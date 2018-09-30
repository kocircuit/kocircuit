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
	"fmt"
	"log"

	"github.com/spf13/cobra"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	"github.com/kocircuit/kocircuit/lang/go/sys"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Execute a Ko program",
	Long:  `Play compiles a Ko program and then executes it.`,
	Run: func(cmd *cobra.Command, args []string) {
		tools := newToolchain()
		if len(args) != 1 {
			log.Fatalf("ko play expects a single argument in the form \"path/to/pkg/Func\"")
		}
		koPkg, koFunc := parsePkgFunc(args[0])
		b := &sys.CompilePlay{
			Repo:    tools.PkgRoots(),
			Pkg:     koPkg,
			Func:    koFunc,
			Faculty: EvalFaculty(),
			Show:    false, // show compiled ko functions
		}
		if result := b.Play(runtime.CompilerContext()); result.Error != nil {
			log.Fatalln(result.Error)
		} else {
			fmt.Printf("%v\n", result.Returned)
		}
	},
}

func init() {
	RootCmd.AddCommand(playCmd)
	initGoBasedCmd(playCmd)
}
