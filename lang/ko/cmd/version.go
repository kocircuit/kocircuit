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

	"github.com/spf13/cobra"

	. "github.com/kocircuit/kocircuit"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version, capability and lineage information",
	Long: `Print information identifying this Ko binary and its origin.
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", KoVersion)
		fmt.Printf("GoFacultyID: %s\n", GoFaculty().ID())
		fmt.Printf("EvalFacultyID: %s\n", EvalFaculty().ID())
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
