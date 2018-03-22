package sys

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGateAt("ko", "DecodePlay", &DecodePlay{})
}

type DecodePlay struct {
	RepoProtoBytes []byte  `ko:"name=repoProtoBytes"`
	Pkg            string  `ko:"name=pkg"`
	Func           string  `ko:"name=func"`
	Faculty        Faculty `ko:"name=faculty"`
	Idiom          Repo    `ko:"name=idiom"`
}

func (arg *DecodePlay) Play(ctx *runtime.Context) *PlayResult {
	repo, err := DecodeRepo(arg.RepoProtoBytes)
	if err != nil {
		panic(&DecodePlayError{Decode: err.Error()})
	}
	w := &Play{
		Pkg:     arg.Pkg,
		Func:    arg.Func,
		Repo:    repo,
		Faculty: PostCompileFaculty(arg.Faculty, "", repo),
		Idiom:   arg.Idiom,
	}
	return w.Play(ctx)
}

type DecodePlayError struct {
	Decode string `ko:"name=decode"`
}
