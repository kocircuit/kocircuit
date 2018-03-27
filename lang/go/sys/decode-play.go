package sys

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/ir"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type DecodePlay struct {
	RepoProtoBytes []byte  `ko:"name=repoProtoBytes"`
	ArgProtoBytes  []byte  `ko:"name=argProtoBytes"`
	Pkg            string  `ko:"name=pkg"`
	Func           string  `ko:"name=func"`
	Faculty        Faculty `ko:"name=faculty"`
}

func (arg *DecodePlay) Play(ctx *runtime.Context) *PlayResult {
	repo, err := DecodeRepo(arg.RepoProtoBytes)
	if err != nil {
		panic(&DecodePlayError{Decode: err.Error()})
	}
	pfe := &PlayFuncEval{
		Func: repo[arg.Pkg][arg.Func],
		Eval: NewEvaluator(arg.Faculty, repo),
	}
	// decode arg
	argSym, err := DecodeSymbol(NewSpan(), pfe.Eval, arg.ArgProtoBytes)
	if err != nil {
		return &PlayResult{
			Error: fmt.Errorf("decoding argument proto (%v)", err),
		}
	}
	switch u := argSym.(type) {
	case EmptySymbol:
	case *StructSymbol:
		pfe.Arg = u
	default:
		return &PlayResult{
			Error: fmt.Errorf("argument must be struct or empty, got %v", u),
		}
	}
	// play
	return pfe.Play(ctx)
}

type DecodePlayError struct {
	Decode string `ko:"name=decode"`
}
