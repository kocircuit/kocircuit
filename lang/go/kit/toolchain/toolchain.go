package toolchain

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	koos "github.com/kocircuit/kocircuit/lib/os"
)

type GoToolchain struct {
	GOROOT string   `ko:"name=GOROOT"`
	GOPATH []string `ko:"name=GOPATH"`
	KOPATH []string `ko:"name=KOPATH"`
	Binary string   `ko:"name=binary"`
}

func NewGoToolchain() *GoToolchain {
	binary, err := exec.LookPath("go")
	if err != nil {
		log.Fatalf("go binary not found (%v)", err)
	}
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Fatal("GOPATH is not set")
	}
	var koPath []string
	if x := os.Getenv("KOPATH"); x != "" {
		koPath = strings.Split(x, ":")
	}
	return &GoToolchain{
		GOROOT: os.Getenv("GOROOT"),
		GOPATH: strings.Split(goPath, ":"),
		KOPATH: koPath,
		Binary: binary,
	}
}

// PkgDir returns the directory containing the given package path.
// If will take the first element of the KOPATH, GOPATH that contains such a package.
func (x *GoToolchain) PkgDir(pkgPath string) string {
	roots := x.PkgRoots()
	for _, p := range roots {
		result := path.Join(p, pkgPath)
		if _, err := os.Stat(result); err == nil {
			return result
		}
	}
	// Not found, return path in first KOPATH element (if any)
	if len(x.KOPATH) > 0 {
		return path.Join(x.KOPATH[0], pkgPath)
	}
	// Not found, return path in first GOPATH element
	return path.Join(x.GOPATH[0], pkgPath)
}

// PkgRoots returns all source folders that may contain packages.
func (x *GoToolchain) PkgRoots() []string {
	result := make([]string, len(x.KOPATH)+len(x.GOPATH))
	for i, p := range x.KOPATH {
		result[i] = path.Join(p, "src")
	}
	for i, p := range x.GOPATH {
		result[i] = path.Join(p, "src")
	}
	return result
}

func (x *GoToolchain) BinaryInstallPath(pkg string) string {
	return path.Join(x.GOPATH[0], "bin", path.Base(pkg))
}

func (x *GoToolchain) BinaryBuildPath(pkg string) string {
	return path.Join(x.GOPATH[0], "src", pkg, path.Base(pkg))
}

func (x *GoToolchain) Env() []string {
	return []string{
		fmt.Sprintf("GOROOT=%s", x.GOROOT),
		fmt.Sprintf("GOPATH=%s", strings.Join(x.GOPATH, ":")),
		fmt.Sprintf("KOPATH=%s", strings.Join(x.KOPATH, ":")),
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
		Dir:    PtrString(gc.Go.GOPATH[0]),
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
