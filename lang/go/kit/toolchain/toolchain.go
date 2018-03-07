package toolchain

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	koos "github.com/kocircuit/kocircuit/lib/os"
)

type GoToolchain struct {
	GOROOT string `ko:"name=GOROOT"`
	GOPATH string `ko:"name=GOPATH"`
	Binary string `ko:"name=binary"`
}

func NewGoToolchain() *GoToolchain {
	binary, err := exec.LookPath("go")
	if err != nil {
		panic("err")
	}
	return &GoToolchain{
		GOROOT: os.Getenv("GOROOT"),
		GOPATH: os.Getenv("GOPATH"),
		Binary: binary,
	}
}

func (x *GoToolchain) PkgDir(pkgPath string) string {
	return path.Join(x.PkgRoot(), pkgPath)
}

func (x *GoToolchain) PkgRoot() string {
	return path.Join(x.GOPATH, "src")
}

func (x *GoToolchain) BinaryInstallPath(pkg string) string {
	return path.Join(x.GOPATH, "bin", path.Base(pkg))
}

func (x *GoToolchain) BinaryBuildPath(pkg string) string {
	return path.Join(x.GOPATH, "src", pkg, path.Base(pkg))
}

func (x *GoToolchain) Env() []string {
	return []string{
		fmt.Sprintf("GOROOT=%s", x.GOROOT),
		fmt.Sprintf("GOPATH=%s", x.GOPATH),
	}
}

type GoClean struct {
	Go      *GoToolchain `ko:"name=go"`
	PkgPath string       `ko:"name=pkg"`
}

func (gc *GoClean) Task(after ...*koos.GoTask) *koos.GoTask {
	return &koos.GoTask{
		Name:   fmt.Sprintf("go_clean:%s", gc.PkgPath),
		Binary: gc.Go.Binary,
		Arg:    []string{"clean", gc.PkgPath},
		Env:    gc.Go.Env(),
		Dir:    PtrString(gc.Go.GOPATH),
		After:  after,
	}
}

type GoBuild struct {
	Go      *GoToolchain `ko:"name=go"`
	PkgPath string       `ko:"name=pkg"`
}

func (gb *GoBuild) Task(after ...*koos.GoTask) *koos.GoTask {
	return &koos.GoTask{
		Name:   fmt.Sprintf("go_build:%s", gb.PkgPath),
		Binary: gb.Go.Binary,
		Arg:    []string{"build", gb.PkgPath},
		Env:    gb.Go.Env(),
		Dir:    PtrString(gb.Go.PkgDir(gb.PkgPath)),
		After:  after,
	}
}

type GoInstall struct {
	Go      *GoToolchain `ko:"name=go"`
	PkgPath string       `ko:"name=pkg"`
}

func (gi *GoInstall) Task(after ...*koos.GoTask) *koos.GoTask {
	return &koos.GoTask{
		Name:   fmt.Sprintf("go_install:%s", gi.PkgPath),
		Binary: gi.Go.Binary,
		Arg:    []string{"install", gi.PkgPath},
		Env:    gi.Go.Env(),
		Dir:    PtrString(gi.Go.PkgDir(gi.PkgPath)),
		After:  after,
	}
}
