package sys

import (
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/ir"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/eval/macros"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type CompilePlay struct {
	Repo    string               `ko:"name=repo"`
	Pkg     string               `ko:"name=pkg"`
	Func    string               `ko:"name=func"`
	Faculty eval.Faculty         `ko:"name=faculty"`
	Arg     *symbol.StructSymbol `ko:"name=arg"` // arg can be nil
	Show    bool                 `ko:"name=show"`
}

func (arg *CompilePlay) Play(ctx *runtime.Context) *PlayResult {
	c := &Compile{
		RepoDir: arg.Repo,
		PkgPath: arg.Pkg,
		Show:    arg.Show,
	}
	compiled := c.Play(ctx)
	if compiled.Error != nil {
		return &PlayResult{Error: compiled.Error}
	}
	w := &Play{
		Pkg:     arg.Pkg,
		Func:    arg.Func,
		Repo:    compiled.Repo,
		Faculty: PostCompileFaculty(arg.Faculty, arg.Repo, compiled.Repo),
		Arg:     arg.Arg,
	}
	return w.Play(ctx)
}

func PostCompileFaculty(baseFaculty eval.Faculty, repoPath string, repo model.Repo) eval.Faculty {
	repoProto, repoProtoBytes, err := ir.SerializeEncodeRepo(repo)
	if err != nil {
		panic(err)
	}
	return eval.MergeFaculty(
		eval.Faculty{
			eval.Ideal{Pkg: "repo", Name: "Path"}:       &macros.EvalGoValueMacro{Value: repoPath},
			eval.Ideal{Pkg: "repo", Name: "Proto"}:      &macros.EvalGoValueMacro{Value: repoProto},
			eval.Ideal{Pkg: "repo", Name: "ProtoBytes"}: &macros.EvalGoValueMacro{Value: repoProtoBytes},
		},
		baseFaculty,
	)
}
