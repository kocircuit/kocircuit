package boot

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Boot struct {
	Faculty Faculty // build
	Repo    Repo    // inherit
	Booter  *Booter
	Func    *Func
	Ctx     Symbol
	Arg     Symbol
}

func (boot *Boot) Play(XXX) (returned, effect Symbol, err error) {
	bootEval := &BootEval{
		Booter: booter,
		Repo:   repo,
		Program: Program{
			Idiom: EvalIdiomRepo,
			Repo:  repo,
			System: System{
				Faculty:  faculty,
				Boundary: &BootEvalBoundary{Booter: booter},
				Combiner: &BootEvalCombiner{Booter: booter},
			},
		},
	}
	return bootEval.Boot(XXX, XXX, XXX)
}

type BootEval struct {
	Booter  *Booter `ko:"name=booter"`
	Repo    Repo    `ko:"name=repo"`
	Program Program `ko:"name=program"`
}

// boot forwards eval panics from the caller evaluator environment.
func (eval *BootEval) Boot(span *Span, f *Func, arg Symbol) (returned, effect Symbol, err error) {
	// evaluation strategy is sequential
	if shape, effect, err := eval.Program.EvalSeq(span, f, arg); err != nil {
		return nil, nil, err
	} else {
		if sym, ok := shape.(Symbol); ok {
			return sym, effect, nil
		} else {
			return nil, effect, nil
		}
	}
}
