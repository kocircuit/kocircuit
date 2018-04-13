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
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
	"github.com/kocircuit/kocircuit/lang/go/sys"
)

// evalCmd represents the eval command
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate a Ko expression",
	Long: `Evaluate a Ko expression against builtin functions.

For example,
	ko eval -e'And(true, Xor(true, false))'
`,
	Run: func(cmd *cobra.Command, args []string) {
		expr := flagEvalExpr
		if expr == "" {
			if buf, err := ioutil.ReadAll(os.Stdin); err != nil {
				log.Fatalf("reading standard input (%v)", err)
			} else {
				expr = string(buf)
			}
		}
		e := &sys.Eval{KoExpr: expr}
		if result := e.Play(runtime.CompilerContext()); result.Error != nil {
			log.Fatalln(result.Error)
		} else {
			fmt.Printf("%v\n", result.Returned)
		}
	},
}

var flagEvalExpr string

func init() {
	RootCmd.AddCommand(evalCmd)
	evalCmd.PersistentFlags().StringVarP(&flagEvalExpr, "expr", "e", "", "Ko expression to evaluate, otherwise read from STDIN.")
}
