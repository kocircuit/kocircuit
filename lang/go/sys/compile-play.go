package sys

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type CompilePlay struct {
	Repo    []string      `ko:"name=repo"`
	Pkg     string        `ko:"name=pkg"`
	Func    string        `ko:"name=func"`
	Faculty Faculty       `ko:"name=faculty"`
	Arg     *StructSymbol `ko:"name=arg"` // arg can be nil
	Show    bool          `ko:"name=show"`
}

func (arg *CompilePlay) Play(ctx *runtime.Context) *PlayResult {
	c := &Compile{
		RepoDirs: arg.Repo,
		PkgPath:  arg.Pkg,
		Show:     arg.Show,
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

func PostCompileFaculty(baseFaculty Faculty, repoPaths []string, repo Repo) Faculty {
	repoProto, repoProtoBytes, err := SerializeEncodeRepo(repo)
	if err != nil {
		panic(err)
	}
	return MergeFaculty(
		Faculty{
			Ideal{Pkg: "repo", Name: "Path"}:       &EvalGoValueMacro{Value: repoPaths[0]},
			Ideal{Pkg: "repo", Name: "Roots"}:      &EvalGoValueMacro{Value: repoPaths},
			Ideal{Pkg: "repo", Name: "Proto"}:      &EvalGoValueMacro{Value: repoProto},
			Ideal{Pkg: "repo", Name: "ProtoBytes"}: &EvalGoValueMacro{Value: repoProtoBytes},
		},
		baseFaculty,
	)
}
