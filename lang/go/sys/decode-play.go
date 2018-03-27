package sys

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	// . "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type DecodePlay struct {
	RepoProtoBytes []byte  `ko:"name=repoProtoBytes"`
	ArgProtoBytes  []byte  `ko:"name=argProtoBytes"` //XXX
	Pkg            string  `ko:"name=pkg"`
	Func           string  `ko:"name=func"`
	Faculty        Faculty `ko:"name=faculty"`
	Idiom          Repo    `ko:"name=idiom"`
}

//XXX: DecodePlayResult (re-encode returned)
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
