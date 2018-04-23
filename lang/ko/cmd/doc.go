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
	"path"

	"github.com/spf13/cobra"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	"github.com/kocircuit/kocircuit/lang/go/sys"
)

// docCmd represents the list command
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Print the documentation of a package or a function.",
	Long: `Doc prints the documentation of a package or function.

Usage:
	ko doc github.com/kocircuit/kocircuit/lib/strings...     # doc for package strings
	ko doc github.com/kocircuit/kocircuit/lib/strings/Join   # doc for function Join in package strings
`,
	Run: func(cmd *cobra.Command, args []string) {
		tools := newToolchain()
		if len(args) != 1 {
			log.Fatalf("ko doc expects a single argument in the form \"path/to/pkg...\" or \"path/to/pkg.Func\"")
		}
		pf := parsePkgOrPkgFunc(args[0])
		b := &sys.Compile{
			RepoDir: path.Join(tools.GOPATH, "src"),
			PkgPath: pf.Pkg,
			Show:    false,
		}
		compileResult := b.Play(runtime.CompilerContext())
		if compileResult.Error != nil {
			log.Fatalln(compileResult.Error)
		}
		//
		repo := CombineRepo(EvalIdiomRepo, compileResult.Repo)
		faculty := EvalFaculty()
		//
		if pf.Func == nil { // package doc
			if doc, ok := repo.DocPackage(pf.Pkg); ok {
				fmt.Println(doc)
			} else if doc, ok = faculty.DocPackage(pf.Pkg); ok {
				fmt.Println(doc)
			} else {
				log.Fatalf("user or builtin package %s not found", pf.Pkg)
			}
		} else { // func doc
			if doc, ok := repo.DocFunc(pf.Pkg, *pf.Func); ok {
				fmt.Println(doc)
			} else if doc, ok = faculty.DocFunc(pf.Pkg, *pf.Func); ok {
				fmt.Println(doc)
			} else {
				log.Fatalf("user or builtin function %s.%s not found", pf.Pkg, *pf.Func)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(docCmd)
}
