package sys

import (
	"fmt"
	"path"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGateAt("ko", "Play", &Play{})
	RegisterEvalGateAt("ko", "CompilePlay", &CompilePlay{})
}

type CompilePlay struct {
	Repo    string  `ko:"name=repo"`
	Pkg     string  `ko:"name=pkg"`
	Func    string  `ko:"name=func"`
	Faculty Faculty `ko:"name=faculty"`
	Idiom   Repo    `ko:"name=idiom"`
	Show    bool    `ko:"name=show"`
}

func (arg *CompilePlay) Play(ctx *runtime.Context) *PlayResult {
	preCompileFaculty := MergeFaculty(
		Faculty{
			Ideal{Pkg: "repo", Name: "Path"}:       &EvalGoValueMacro{Value: arg.Repo},
			Ideal{Pkg: "repo", Name: "Proto"}:      &EvalPlaceholderMacro{},
			Ideal{Pkg: "repo", Name: "ProtoBytes"}: &EvalPlaceholderMacro{},
		},
		arg.Faculty,
	)
	c := &Compile{
		RepoDir: arg.Repo,
		PkgPath: arg.Pkg,
		Faculty: preCompileFaculty,
		Idiom:   arg.Idiom,
		Show:    arg.Show,
	}
	compiled := c.Play(ctx)
	if compiled.Error != nil {
		return &PlayResult{Error: compiled.Error}
	}
	repoProto, repoProtoBytes, err := SerializeEncodeRepo(compiled.Repo)
	if err != nil {
		panic(err)
	}
	postCompileFaculty := MergeFaculty(
		Faculty{
			Ideal{Pkg: "repo", Name: "Path"}:       &EvalGoValueMacro{Value: arg.Repo},
			Ideal{Pkg: "repo", Name: "Proto"}:      &EvalGoValueMacro{Value: repoProto},
			Ideal{Pkg: "repo", Name: "ProtoBytes"}: &EvalGoValueMacro{Value: repoProtoBytes},
		},
		arg.Faculty,
	)
	w := &Play{
		Pkg:     arg.Pkg,
		Func:    arg.Func,
		Repo:    compiled.Repo,
		Faculty: postCompileFaculty,
		Idiom:   arg.Idiom,
	}
	return w.Play(ctx)
}

type Play struct {
	Pkg     string  `ko:"name=pkg"`  // e.g. github.com/kocircuit/kocircuit/codelab
	Func    string  `ko:"name=func"` // e.g. HelloWorld
	Repo    Repo    `ko:"name=repo"` // compiled ko repo
	Faculty Faculty `ko:"name=faculty"`
	Idiom   Repo    `ko:"name=idiom"`
}

type PlayResult struct {
	Play     *Play       `ko:"name=play"`
	Returned interface{} `ko:"name=returned"`
	Error    error       `ko:"name=error"`
}

func (w *Play) Play(ctx *runtime.Context) *PlayResult {
	r := &PlayResult{Play: w}
	fu := w.Repo[w.Pkg][w.Func]
	if fu == nil {
		r.Error = fmt.Errorf("cannot find main circuit %s", path.Join(w.Pkg, w.Func))
		return r
	}
	span := NewSpan()
	ev := NewEvaluator(w.Faculty, w.Repo)
	r.Returned, _, r.Error = ev.Eval(span, fu, MakeStructSymbol(nil))
	return r
}
