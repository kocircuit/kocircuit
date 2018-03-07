// Package eval provides evaluation of Ko circuits.
package eval

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Evaluate struct {
	Repo    Repo    `ko:"name=repo"`
	Program Program `ko:"name=program"`
}

func NewEvaluator(faculty Faculty, repo Repo) *Evaluate {
	return &Evaluate{
		Repo: repo,
		Program: Program{
			Idiom: EvalIdiomRepo,
			Repo:  repo,
			System: System{
				Faculty:  faculty,
				Boundary: EvalBoundary{},
				Combiner: EvalCombiner{},
			},
		},
	}
}

func (eval *Evaluate) Eval(span *Span, f *Func, arg Symbol) (Symbol, Effect, error) {
	if shape, effect, err := eval.Program.EvalSeq(span, f, arg); err != nil { // top-level evaluation strategy is sequential
		return nil, nil, err
	} else {
		if sym, ok := shape.(Symbol); ok {
			return sym, effect, nil
		} else {
			return nil, effect, nil
		}
	}
}

type EvalBoundary struct{}

func (EvalBoundary) Figure(span *Span, figure Figure) (Shape, Effect, error) {
	switch u := figure.(type) {
	case Bool:
		return BasicSymbol{u.Value_}, nil, nil
	case Integer:
		return BasicSymbol{u.Value_}, nil, nil
	case Float:
		return BasicSymbol{u.Value_}, nil, nil
	case String:
		return BasicSymbol{u.Value_}, nil, nil
	case Macro:
		// macro is either a macro from registry, or from Interpret()
		return MakeVarietySymbol(u, nil), nil, nil
	}
	panic("unknown figure")
}

func (EvalBoundary) Enter(span *Span, arg Arg) (Shape, Effect, error) {
	return arg.(Symbol), nil, nil
}

func (EvalBoundary) Leave(span *Span, shape Shape) (Return, Effect, error) {
	return shape, nil, nil
}

type EvalCombiner struct{}

func (EvalCombiner) Interpret(eval Evaluator, f *Func) Macro {
	return &EvalInterpretMacro{Evaluator: eval, Func: f}
}

func (EvalCombiner) Combine(
	span *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidue StepResidues,
) (Effect, error) {
	return nil, nil
}

type EvalInterpretMacro struct {
	Evaluator Evaluator `ko:"name=evaluator"`
	Func      *Func     `ko:"name=func"`
}

func (m *EvalInterpretMacro) Splay() Tree {
	return Quote{m.Help()}
}

func (m *EvalInterpretMacro) MacroID() string { return m.Help() }

func (m *EvalInterpretMacro) Label() string { return "eval" }

func (m *EvalInterpretMacro) MacroSheathString() *string { return nil }

func (m *EvalInterpretMacro) Help() string {
	return fmt.Sprintf("%s", m.Func.FullPath())
}

func (m *EvalInterpretMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return m.InvokeSeq(span, arg) // default circuit execution mode
}

func (m *EvalInterpretMacro) InvokeSeq(span *Span, arg Arg) (Return, Effect, error) {
	ss := arg.(*StructSymbol)
	return m.Evaluator.EvalSeq(span, m.Func, RewriteMonadicField(ss, m.Func.Monadic))
}

func (m *EvalInterpretMacro) InvokePar(span *Span, arg Arg) (Return, Effect, error) {
	ss := arg.(*StructSymbol)
	return m.Evaluator.EvalPar(span, m.Func, RewriteMonadicField(ss, m.Func.Monadic))
}

// RewriteMonadicField rewrites the struct so that its monadic field has Ko name asKoName.
func RewriteMonadicField(ss *StructSymbol, asKoName string) *StructSymbol {
	fields := make(FieldSymbols, len(ss.Field))
	copy(fields, ss.Field)
	for i, field := range fields {
		if field.Monadic {
			fields[i] = &FieldSymbol{
				Name:  asKoName,
				Value: field.Value,
			}
			break
		}
	}
	return MakeStructSymbol(fields)
}
