package sys

import (
	"fmt"
	"path"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type Play struct {
	Pkg     string        `ko:"name=pkg"`  // e.g. github.com/kocircuit/kocircuit/codelab
	Func    string        `ko:"name=func"` // e.g. HelloWorld
	Repo    Repo          `ko:"name=repo"` // compiled ko repo
	Faculty Faculty       `ko:"name=faculty"`
	Arg     *StructSymbol `ko:"name=arg"` // arg can be nil
}

func (w *Play) Play(ctx *runtime.Context) *PlayResult {
	pfe := &PlayFuncEval{
		Func: w.Repo[w.Pkg][w.Func],
		Eval: NewEvaluator(w.Faculty, w.Repo),
		Arg:  w.Arg,
	}
	if pfe.Func == nil {
		return &PlayResult{
			Play:  w,
			Error: fmt.Errorf("cannot find main circuit %s", path.Join(w.Pkg, w.Func)),
		}
	}
	return pfe.Play(ctx)
}

type PlayFuncEval struct {
	Func *Func         `ko:"name=func"`
	Eval *Evaluate     `ko:"name=eval"`
	Arg  *StructSymbol `ko:"name=arg"` // arg can be nil
}

func (w *PlayFuncEval) Play(ctx *runtime.Context) *PlayResult {
	r := &PlayResult{PlayFuncEval: w}
	span := NewSpan()
	var arg Symbol
	if w.Arg == nil {
		arg = MakeStructSymbol(nil)
	} else {
		arg = w.Arg
	}
	r.Returned, _, _, r.Error = w.Eval.Eval(span, w.Func, arg)
	return r
}

type PlayResult struct {
	Play         *Play         `ko:"name=play"`
	PlayFuncEval *PlayFuncEval `ko:"name=playFuncEval"`
	Returned     Symbol        `ko:"name=returned"`
	Error        error         `ko:"name=error"`
}
