package translate

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Translation struct {
	Repo    Repo    `ko:"name=repo"`
	Program Program `ko:"name=program"`
}

func NewTranslation(faculty Faculty, repo Repo) *Translation {
	return &Translation{
		Repo: repo,
		Program: Program{
			Idiom: EvalIdiomRepo,
			Repo:  repo,
			System: System{
				Faculty:  faculty,
				Boundary: TranslationBoundary{},
				Combiner: TranslationCombiner{},
			},
		},
	}
}

func (eval *Translation) Translate(span *Span, f *Func, arg XXXSymbol) (returned XXXSymbol, eff XXXEffect, err error) {
	// catch unrecovered evaluator panics
	defer func() {
		if r := recover(); r != nil {
			evalPanic := r.(*EvalPanic)
			returned, eff, err = nil, nil, evalPanic.Origin.Errorf(nil, "unrecovered panic: %v", evalPanic.Panic)
			return
		}
	}()
	// top-level evaluation strategy is sequential
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
