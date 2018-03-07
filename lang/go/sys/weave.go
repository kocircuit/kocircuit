package sys

import (
	"fmt"
	"path"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/console"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"

	"github.com/kocircuit/kocircuit/lib/os"
)

func init() {
	RegisterGoGateAt("ko", "Weave", &Weave{})
	RegisterGoGateAt("ko", "CompileWeave", &CompileWeave{})
}

type CompileWeave struct {
	Repo      string       `ko:"name=repo"`
	Pkg       string       `ko:"name=pkg"`  // e.g. github.com/kocircuit/kocircuit/codelab
	Func      string       `ko:"name=func"` // e.g. Main
	Faculty   Faculty      `ko:"name=faculty"`
	Idiom     Repo         `ko:"name=idiom"`
	Arg       *GoStruct    `ko:"name=arg"`
	Toolchain *GoToolchain `ko:"name=toolchain"`
	GoKoRoot  string       `ko:"name=weavingRoot"` // root of generated go code
	GoKoPkg   string       `ko:"name=goKoPkg"`     // root of generated circuit relative to GoKoRoot
	Show      bool         `ko:"name=show"`
}

func (arg *CompileWeave) Play(ctx *runtime.Context) *WeaveResult {
	c := &Compile{
		RepoDir: arg.Repo,
		PkgPath: arg.Pkg,
		Faculty: arg.Faculty,
		Idiom:   arg.Idiom,
		Show:    arg.Show,
	}
	compiled := c.Play(ctx)
	if compiled.Error != nil {
		return &WeaveResult{Error: compiled.Error}
	}
	w := &Weave{
		Pkg:       arg.Pkg,
		Func:      arg.Func,
		Repo:      compiled.Repo,
		Faculty:   arg.Faculty,
		Idiom:     arg.Idiom,
		Arg:       arg.Arg,
		Toolchain: arg.Toolchain,
		GoKoRoot:  arg.GoKoRoot,
		GoKoPkg:   arg.GoKoPkg,
	}
	return w.Play(ctx)
}

type Weave struct {
	Pkg       string       `ko:"name=pkg"`  // e.g. github.com/kocircuit/kocircuit/codelab
	Func      string       `ko:"name=func"` // e.g. HelloWorld
	Repo      Repo         `ko:"name=repo"` // compiled ko repo
	Faculty   Faculty      `ko:"name=faculty"`
	Idiom     Repo         `ko:"name=idiom"`
	Arg       *GoStruct    `ko:"name=arg"`
	Toolchain *GoToolchain `ko:"name=toolchain"`
	GoKoRoot  string       `ko:"name=weavingRoot"` // root of generated go code
	GoKoPkg   string       `ko:"name=goKoPkg"`     // root of generated circuit relative to GoKoRoot
}

type WeaveResult struct {
	Weave      *Weave        `ko:"name=weave"`
	Instrument *GoInstrument `ko:"name=instrument"`
	Error      error         `ko:"name=error"`
	RenderCtx  *RenderCtx    `ko:"name=renderCtx"`
	MainGo     string        `ko:"name=mainGo"` // path to main.go
	MainBinary string        `ko:"name=mainBinary"`
}

func (w *Weave) Play(ctx *runtime.Context) *WeaveResult {
	r := &WeaveResult{Weave: w}
	fu := w.Repo[w.Pkg][w.Func]
	if fu == nil {
		r.Error = fmt.Errorf("cannot find main circuit %s", path.Join(w.Pkg, w.Func))
		return r
	}

	// weave
	warg := w.Arg
	if warg == nil {
		warg = NewGoStruct()
	}
	weaveCtx := NewGoWeaveCtx("WEAVE", w.Repo, w.Faculty, w.Idiom)
	weaveID := Mix(w.Pkg, w.Func, warg.TypeID())
	topPkg := fmt.Sprintf("%s_%s", w.Func, weaveID)
	span := RefineWeaveCtx(NewSpan(), weaveCtx)
	span = RefineChamber(span, topPkg)
	if r.Instrument, r.Error = weaveCtx.WeaveInstrument(span, fu, warg); r.Error != nil {
		return r
	}
	stats := r.Instrument.Stats()
	nm := ANSI.Ref(fmt.Sprintf("%s.%s", w.Pkg, w.Func))
	ctx.Printf(
		"woven circuit=%s subcircuits=%d steps=%d types=%d recursions=%d iterations=%d",
		nm, stats.TotalFunc, stats.TotalStep, stats.TotalType,
		r.Instrument.ProgramEffect.WeavingStat.RecursionCount,
		r.Instrument.ProgramEffect.WeavingStat.IterationCount,
	)
	ctx.Printf("weaving returns=%s", Sprint(r.Instrument.Returns))

	r.RenderCtx = &RenderCtx{
		Root: path.Join(w.GoKoRoot, w.GoKoPkg),
	}

	// determine GOROOT-relative path to main/program package
	mainGroupPath := GoGroupPath{Group: KoPkgGroup, Path: ""}
	mainGoKoPkg := r.RenderCtx.UnifiedPath(mainGroupPath)
	// make directory if not there
	mkdir := &os.GoMkdir{Path: w.Toolchain.PkgDir(mainGoKoPkg)}
	if r.Error = mkdir.Play(ctx); r.Error != nil {
		return r
	}
	// clean generated go code
	wipe := &os.GoWipeDir{Dir: w.Toolchain.PkgDir(mainGoKoPkg)}
	if !wipe.Play(ctx) {
		r.Error = fmt.Errorf("wiping directory %q", wipe.Dir)
		return r
	}
	clean := &GoClean{Go: w.Toolchain, PkgPath: mainGoKoPkg}
	if !clean.Task().Play(ctx) {
		r.Error = fmt.Errorf("go clean error")
		return r
	}

	// materialize go implementation to disk at GOPATH/src/GoKoRoot/GoKoPkg/Func_TypeID
	mainPkgCtx := r.RenderCtx.PkgContext(mainGroupPath)
	r.MainBinary = path.Base(mainGoKoPkg)
	mainGo := fmt.Sprintf("%s.go", topPkg)
	r.MainGo = path.Join(mainGoKoPkg, mainGo)
	allFiles := r.RenderCtx.Shred(r.Instrument.Circuit, r.Instrument.Directive, mainGroupPath)
	mainFile := RenderGoPkgFile(
		mainPkgCtx,
		"main",
		mainGo,
		&GoMainExpr{
			Directive: FilterDirective(r.Instrument.Directive, mainGroupPath),
			Valve:     r.Instrument.Valve,
		},
	)
	allFiles = append(allFiles, mainFile)
	ctx.Printf("generated main %s\n", r.MainGo)

	r.Error = SourceRepo(allFiles).Materialize(w.Toolchain.PkgRoot())
	return r
}
