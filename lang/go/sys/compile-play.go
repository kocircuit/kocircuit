package sys

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
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
	c := &Compile{
		RepoDir: arg.Repo,
		PkgPath: arg.Pkg,
		Faculty: PreCompileFaculty(arg.Faculty),
		Idiom:   arg.Idiom,
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
		Idiom:   arg.Idiom,
	}
	return w.Play(ctx)
}

func PreCompileFaculty(baseFaculty Faculty) Faculty {
	return MergeFaculty(
		Faculty{
			Ideal{Pkg: "repo", Name: "Path"}:       &EvalPlaceholderMacro{},
			Ideal{Pkg: "repo", Name: "Proto"}:      &EvalPlaceholderMacro{},
			Ideal{Pkg: "repo", Name: "ProtoBytes"}: &EvalPlaceholderMacro{},
		},
		baseFaculty,
	)
}

func PostCompileFaculty(baseFaculty Faculty, repoPath string, repo Repo) Faculty {
	repoProto, repoProtoBytes, err := SerializeEncodeRepo(repo)
	if err != nil {
		panic(err)
	}
	return MergeFaculty(
		Faculty{
			Ideal{Pkg: "repo", Name: "Path"}:       &EvalGoValueMacro{Value: repoPath},
			Ideal{Pkg: "repo", Name: "Proto"}:      &EvalGoValueMacro{Value: repoProto},
			Ideal{Pkg: "repo", Name: "ProtoBytes"}: &EvalGoValueMacro{Value: repoProtoBytes},
		},
		baseFaculty,
	)
}
