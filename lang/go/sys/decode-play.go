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

type Decode struct {
	Repo    []byte  `ko:"name=repo"`
	Arg     []byte  `ko:"name=arg"`
	Pkg     string  `ko:"name=pkg"`
	Func    string  `ko:"name=func"`
	Faculty Faculty `ko:"name=faculty"`
}

type DecodeResult struct {
	Faculty Faculty       `ko:"name=faculty"`
	Repo    Repo          `ko:"name=repo"`
	Func    *Func         `ko:"name=func"`
	Eval    *Evaluate     `ko:"name=eval"`
	Arg     *StructSymbol `ko:"name=arg"`
}

func (in *Decode) Play(ctx *runtime.Context) (result *DecodeResult, err error) {
	result = &DecodeResult{}
	if result.Repo, err = DecodeRepo(in.Repo); err != nil {
		return nil, fmt.Errorf("decoding repo (%v)", err)
	}
	if result.Func = result.Repo[in.Pkg][in.Func]; result.Func == nil {
		return nil, fmt.Errorf("cannot find function %s.%s", in.Pkg, in.Func)
	}
	result.Eval = NewEvaluator(in.Faculty, result.Repo)
	var sym Symbol
	if sym, err = DecodeSymbol(NewSpan(), result.Eval, in.Arg); err != nil {
		return nil, fmt.Errorf("decoding arg (%v)", err)
	}
	switch u := sym.(type) {
	case EmptySymbol:
		result.Arg = nil
	case *StructSymbol:
		result.Arg = u
	default:
		return nil, fmt.Errorf("arg must be structure or empty, got %v", u)
	}
	return result, nil
}

type DecodePlay DecodeResult

func (in *DecodePlay) Play(ctx *runtime.Context) *PlayResult {
	pfe := &PlayFuncEval{
		Func: in.Func,
		Eval: in.Eval,
		Arg:  in.Arg,
	}
	return pfe.Play(ctx)
}
