package translate

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type TranslateEvaluator struct {
	Repo    Repo    `ko:"name=repo"`
	Program Program `ko:"name=program"`
}

func NewTranslateEvaluator(faculty Faculty, repo Repo) *TranslateEvaluator {
	return &TranslateEvaluator{
		Repo: repo,
		Program: Program{
			Idiom: EvalIdiomRepo,
			Repo:  repo,
			System: System{
				Faculty:  faculty,
				Boundary: EvalBoundary{}, //XXX
				Combiner: EvalCombiner{}, //XXX
			},
		},
	}
}
