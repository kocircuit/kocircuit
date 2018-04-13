package cmd

import (
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"

	"github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
)

func parsePkgFunc(s string) (pkgPath, funcName string) {
	return path.Dir(s), path.Base(s)
}

var (
	flagGoBinary string // go binary path
	flagGOROOT   string // GOROOT
	flagGOPATH   string // GOPATH
	flagKOGO     string // generate Go code at this path
)

func newToolchain() *toolchain.GoToolchain {
	return &toolchain.GoToolchain{
		GOROOT: flagGOROOT,
		GOPATH: flagGOPATH,
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
	cmd.PersistentFlags().StringVarP(&flagKOGO, "kogo", "", "kogo", "Path for generated Go code, relative to GOPATH/src")
}
