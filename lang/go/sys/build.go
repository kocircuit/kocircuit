package sys

import (
	"fmt"
	"path"
	"time"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"

	kos "github.com/kocircuit/kocircuit/lib/os"
)

func init() {
	RegisterGoGateAt("ko", "Build", &Build{})
	RegisterGoGateAt("ko", "BuildString", &BuildString{})
	RegisterEvalGateAt("ko", "Build", &Build{})
	RegisterEvalGateAt("ko", "BuildString", &BuildString{})
}

type BuildString struct {
	KoString  string       `ko:"name=koString"`
	KoPkg     string       `ko:"name=koPkg"`
	KoFunc    string       `ko:"name=koFunc"`
	Faculty   Faculty      `ko:"name=faculty"`
	Idiom     Repo         `ko:"name=idiom"`
	Arg       *GoStruct    `ko:"name=arg"`
	Toolchain *GoToolchain `ko:"name=toolchain"`
	GoKoRoot  string       `ko:"name=goKoRoot"`
	GoKoPkg   string       `ko:"name=goKoPkg"`
}

func (b *BuildString) Weave(repo Repo) *Weave {
	return &Weave{
		Pkg:       b.KoPkg,
		Func:      b.KoFunc,
		Repo:      repo,
		Faculty:   b.Faculty,
		Idiom:     b.Idiom,
		Arg:       b.Arg,
		Toolchain: b.Toolchain,
		GoKoRoot:  b.GoKoRoot,
		GoKoPkg:   b.GoKoPkg,
	}
}

func (b *BuildString) MainGoPkg() string {
	return path.Join(b.GoKoRoot, b.GoKoPkg)
}

func (b *BuildString) Play(ctx *runtime.Context) *BuildResult {
	r := &BuildResult{}
	repo, err := CompileString(b.KoPkg, "source.ko", b.KoString)
	if err != nil {
		r.Error = err
		return r
	}
	r.Repo = repo
	weaved := b.Weave(repo).Play(ctx)
	if weaved.Error != nil {
		r.Error = weaved.Error
		return r
	}
	r.Returns = weaved.Instrument.Returns
	gb := &GoBuild{Go: b.Toolchain, PkgPath: b.MainGoPkg()}
	if !gb.Task().Play(ctx) {
		r.Error = fmt.Errorf("go build error")
		return r
	}
	return r
}

type Build struct {
	KoRepo    string       `ko:"name=koRepo"`
	KoPkg     string       `ko:"name=koPkg"`
	KoFunc    string       `ko:"name=koFunc"`
	Faculty   Faculty      `ko:"name=faculty"`
	Idiom     Repo         `ko:"name=idiom"`
	Arg       *GoStruct    `ko:"name=arg"`
	Toolchain *GoToolchain `ko:"name=toolchain"`
	GoKoRoot  string       `ko:"name=goKoRoot"`
	GoKoPkg   string       `ko:"name=goKoPkg"`
	Install   bool         `ko:"name=install"`
	EvalExpr  string       `ko:"name=evalExpr"`
}

type BuildInstall struct{}

type BuildRun struct {
	Arg []string `ko:"name=arg"`
}

type BuildResult struct {
	Returns  GoType         `ko:"name=returns"`
	Repo     Repo           `ko:"name=repo"`
	Compiled *CompileResult `ko:"name=compiled"`
	Error    error          `ko:"name=error"`
}

func (b *Build) Compile() *Compile {
	return &Compile{
		RepoDir: b.KoRepo,
		PkgPath: b.KoPkg,
		Faculty: b.Faculty,
		Idiom:   b.Idiom,
	}
}

func (b *Build) Weave(repo Repo) *Weave {
	return &Weave{
		Pkg:       b.KoPkg,
		Func:      b.KoFunc,
		Repo:      repo,
		Faculty:   b.Faculty,
		Idiom:     b.Idiom,
		Arg:       b.Arg,
		Toolchain: b.Toolchain,
		GoKoRoot:  b.GoKoRoot,
		GoKoPkg:   b.GoKoPkg,
	}
}

func (b *Build) MainGoPkg() string {
	return path.Join(b.GoKoRoot, b.GoKoPkg)
}

func (b *Build) Play(ctx *runtime.Context) *BuildResult {
	r := &BuildResult{}
	r.Compiled = b.Compile().Play(ctx)
	r.Repo = r.Compiled.Repo
	if r.Compiled.Error != nil {
		r.Error = r.Compiled.Error
		return r
	}
	t0 := time.Now()
	weaved := b.Weave(r.Compiled.Repo).Play(ctx)
	if weaved.Error != nil {
		r.Error = weaved.Error
		return r
	}
	t1 := time.Now()
	r.Returns = weaved.Instrument.Returns
	// build or install
	mainPkg := path.Dir(weaved.MainGo)
	fullBinary := ""
	if b.Install {
		gcmd := &GoInstall{Go: b.Toolchain, PkgPath: mainPkg}
		if !gcmd.Task().Play(ctx) {
			r.Error = fmt.Errorf("go install error")
		}
		ctx.Printf("installed %s", weaved.MainBinary)
		fullBinary = b.Toolchain.BinaryInstallPath(mainPkg)
	} else {
		gcmd := &GoBuild{Go: b.Toolchain, PkgPath: mainPkg}
		if !gcmd.Task().Play(ctx) {
			r.Error = fmt.Errorf("go build error")
		}
		ctx.Printf("built %s", path.Join(mainPkg, weaved.MainBinary))
		fullBinary = b.Toolchain.BinaryBuildPath(mainPkg)
	}
	t2 := time.Now()
	ctx.Printf("ko weaving %.2fs, go compilation %.2fs", t1.Sub(t0).Seconds(), t2.Sub(t1).Seconds())
	// evaluate expression using newly built binary
	if r.Error == nil {
		if b.EvalExpr != "" {
			task := &kos.GoTask{
				Name:   "eval",
				Binary: fullBinary,
				Arg:    []string{"eval", "--expr", b.EvalExpr},
			}
			if out, err := task.RunWithOutput(); err != nil {
				r.Error = err
			} else {
				fmt.Print(out)
			}
		}
	}
	return r
}
