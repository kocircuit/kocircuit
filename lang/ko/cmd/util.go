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
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
)

type pkgFunc struct {
	Pkg  string
	Func *string
}

// parsePkgOrPkgFunc parses a package or a function reference.
// Package refernces are of the form:
// 	path/to/pkg...
// 	path/to/pkg/...
// Function references are of the forms:
// 	path/to/pkg/Func
// 	Func
func parsePkgOrPkgFunc(s string) *pkgFunc {
	if pkg, ok := parsePkg(s); ok {
		return &pkgFunc{Pkg: pkg, Func: nil}
	} else {
		dotPkg, dotFu := parsePkgFunc(s)
		return &pkgFunc{Pkg: dotPkg, Func: &dotFu}
	}
}

// parsePkg parses a package path in the forms:
// 	path/to/pkg...
// 	path/to/pkg/...
func parsePkg(s string) (pkgPath string, ok bool) {
	if !strings.HasSuffix(s, "...") {
		return "", false
	} else {
		return path.Clean(s[:len(s)-len("...")]), true
	}
}

// parsePkgDotFunc parses a string of either form:
//	"path/to/pkg/Func" or "Func".
func parsePkgFunc(s string) (pkgPath, funcName string) {
	dir := path.Dir(s)
	if dir == "." {
		dir = ""
	}
	return dir, path.Base(s)
}

var (
	flagGoBinary string // go binary path
	flagGOROOT   string // GOROOT
	flagGOPATH   string // GOPATH
)

func newToolchain() *toolchain.GoToolchain {
	return &toolchain.GoToolchain{
		GOROOT: flagGOROOT,
		GOPATH: strings.Split(flagGOPATH, ":"),
		Binary: flagGoBinary,
	}
}

func initGoBasedCmd(cmd *cobra.Command) {
	gobinary, err := exec.LookPath("go")
	if err != nil {
		panic("err")
	}
	goroot, gopath := os.Getenv("GOROOT"), os.Getenv("GOPATH")
	cmd.PersistentFlags().StringVarP(&flagGoBinary, "gobinary", "", gobinary, "Path to Go binary")
	cmd.PersistentFlags().StringVarP(&flagGOROOT, "goroot", "", goroot, "GOROOT setting")
	cmd.PersistentFlags().StringVarP(&flagGOPATH, "gopath", "", gopath, "GOPATH setting")
}
