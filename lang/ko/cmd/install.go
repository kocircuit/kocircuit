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

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "(EXP) Build and install (via go install) a Ko program",
	Long: `Installing is identical to building, while additionally
the resulting binary is installed via "go install".

For example,
	ko install github.com/kocircuit/kocircuit/codelab/fib/CodelabFib
`,
	Run: func(cmd *cobra.Command, args []string) {
		tools := newToolchain()
		if len(args) != 1 {
			log.Fatalf("ko install expects a single argument in the form \"path/to/pkg/Func\"")
		}
		koPkg, koFunc := parsePkgFunc(args[0])
		b := &sys.Build{
			KoRepo:    path.Join(tools.GOPATH, "src"),
			KoPkg:     koPkg,
			KoFunc:    koFunc,
			Faculty:   GoFaculty(),
			Idiom:     GoIdiomRepo,
			Arg:       nil,
			Toolchain: tools,
			GoKoRoot:  flagKOGO,
			GoKoPkg:   path.Join(koPkg, koFunc),
			Install:   true,
		}
		if result := b.Play(runtime.CompilerContext()); result.Error != nil {
			log.Fatalln(result.Error)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
	initGoBasedCmd(installCmd)
	// Cobra supports local flags which will only run when this command is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
