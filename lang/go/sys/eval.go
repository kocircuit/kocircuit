package sys

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGateAt("ko", "Eval", &Eval{})
}

type Eval struct {
	KoExpr string `ko:"name=koExpr"`
}

type EvalResult struct {
	Eval     *Eval  `ko:"name=eval"`
	Error    error  `ko:"name=error"`
	Returned Symbol `ko:"name=returned"`
}

func (e *Eval) Play(ctx *runtime.Context) *EvalResult {
	r := &EvalResult{Eval: e}
	framedExpr := fmt.Sprintf("Cell() {\nreturn: %s\n}", e.KoExpr)
	repo, err := CompileString("jail", "expr.ko", framedExpr)
	if err != nil {
		r.Error = err
		return r
	}
	ev := NewEvaluator(EvalFaculty(), repo)
	span := NewSpan()
	r.Returned, _, r.Error = ev.Eval(span, repo["jail"]["Cell"], MakeStructSymbol(nil))
	return r
}
