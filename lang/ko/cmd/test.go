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

package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	"github.com/kocircuit/kocircuit/lang/go/sys"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Execute all Test functions of a Ko program",
	Long:  `Play compiles a Ko program and then executes all Test functions.`,
	Run: func(cmd *cobra.Command, args []string) {
		tools := newToolchain()
		if len(args) < 1 {
			log.Fatalf("ko test expects aat least one argument in the form \"path/to/pkg\"")
		}
		b := &sys.CompileTest{
			Repo:    tools.PkgRoots(),
			Pkgs:    args,
			Faculty: eval.EvalFaculty(),
			Show:    false, // show compiled ko functions
		}
		if result := b.Play(runtime.CompilerContext()); result.Error != nil {
			log.Fatalln(result.Error)
		}
	},
}

func init() {
	RootCmd.AddCommand(testCmd)
	initGoBasedCmd(testCmd)
}
