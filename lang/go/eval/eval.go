// Package eval provides evaluation of Ko circuits.
package eval

import (
	"fmt"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Evaluate struct {
	Repo    model.Repo   `ko:"name=repo"`
	Program eval.Program `ko:"name=program"`
}

func NewEvaluator(faculty eval.Faculty, repo model.Repo) *Evaluate {
	return &Evaluate{
		Repo: repo,
		Program: eval.Program{
			Idiom: EvalIdiomRepo,
			Repo:  repo,
			System: eval.System{
				Faculty:  faculty,
				Boundary: EvalBoundary{},
				Combiner: EvalCombiner{},
			},
		},
	}
}

type EvalPanic struct {
	Origin *model.Span   `ko:"name=origin"`
	Panic  symbol.Symbol `ko:"name=panic"`
}

func NewEvalPanic(origin *model.Span, panik symbol.Symbol) *EvalPanic {
	return &EvalPanic{Origin: origin, Panic: panik}
}

func (eval *Evaluate) AssembleMacro(span *model.Span, pkgPath, funcName string) (eval.Macro, error) {
	if fu := eval.Repo.Lookup(pkgPath, funcName); fu == nil {
		return nil, span.Errorf(nil, "function %s.%s not found", pkgPath, funcName)
	} else {
		return EvalCombiner{}.Interpret(eval.Program, fu), nil
	}
}

func (eval *Evaluate) Eval(span *model.Span, f *model.Func, arg symbol.Symbol) (returned symbol.Symbol, panicked symbol.Symbol, eff eval.Effect, err error) {
	// catch unrecovered evaluator panics
	defer func() {
		if r := recover(); r != nil {
			evalPanic := r.(*EvalPanic)
			returned, panicked = nil, evalPanic.Panic
			eff, err = nil, evalPanic.Origin.Errorf(nil, "unrecovered panic: %v", evalPanic.Panic)
			return
		}
	}()
	// top-level evaluation strategy is sequential
	if shape, effect, err := eval.Program.EvalSeq(span, f, arg); err != nil {
		return nil, nil, nil, err
	} else {
		if sym, ok := shape.(symbol.Symbol); ok {
			return sym, nil, effect, nil
		} else {
			return nil, nil, effect, nil
		}
	}
}

type EvalBoundary struct{}

func (EvalBoundary) Figure(span *model.Span, figure eval.Figure) (eval.Shape, eval.Effect, error) {
	switch u := figure.(type) {
	case eval.Bool:
		return symbol.BasicSymbol{Value: u.Value_}, nil, nil
	case eval.Integer:
		return symbol.BasicSymbol{Value: u.Value_}, nil, nil
	case eval.Float:
		return symbol.BasicSymbol{Value: u.Value_}, nil, nil
	case eval.String:
		return symbol.BasicSymbol{Value: u.Value_}, nil, nil
	case eval.Macro:
		// macro is either a macro from registry, or from Interpret()
		return symbol.MakeVarietySymbol(u, nil), nil, nil
	}
	panic("unknown figure")
}

func (EvalBoundary) Enter(span *model.Span, arg eval.Arg) (eval.Shape, eval.Effect, error) {
	return arg.(symbol.Symbol), nil, nil
}

func (EvalBoundary) Leave(span *model.Span, shape eval.Shape) (eval.Return, eval.Effect, error) {
	return shape, nil, nil
}

type EvalCombiner struct{}

func (EvalCombiner) Interpret(eval eval.Evaluator, f *model.Func) eval.Macro {
	return &EvalInterpretMacro{Evaluator: eval, Func: f}
}

func (EvalCombiner) Combine(
	span *model.Span,
	f *model.Func,
	arg eval.Arg,
	returned eval.Return,
	stepResidue eval.StepResidues,
) (eval.Effect, error) {
	return nil, nil
}

type EvalInterpretMacro struct {
	Evaluator eval.Evaluator `ko:"name=evaluator"`
	Func      *model.Func    `ko:"name=func"`
}

// InterpretFunc communicates to Variety.Disassemble the underlying function identity.
func (m *EvalInterpretMacro) InterpretFunc() (pkgPath, funcName string) {
	return m.Func.Pkg, m.Func.Name
}

func (m *EvalInterpretMacro) Splay() tree.Tree {
	return tree.Quote{String_: m.Help()}
}

func (m *EvalInterpretMacro) MacroID() string { return m.Help() }

func (m *EvalInterpretMacro) Label() string { return "eval" }

func (m *EvalInterpretMacro) MacroSheathString() *string { return nil }

func (m *EvalInterpretMacro) Help() string {
	return fmt.Sprintf("%s", m.Func.FullPath())
}

func (m *EvalInterpretMacro) Doc() string {
	return m.Func.DocLong()
}

func (m *EvalInterpretMacro) Invoke(span *model.Span, arg eval.Arg) (eval.Return, eval.Effect, error) {
	return m.InvokeSeq(span, arg) // default circuit execution mode
}

func (m *EvalInterpretMacro) InvokeSeq(span *model.Span, arg eval.Arg) (eval.Return, eval.Effect, error) {
	ss := arg.(*symbol.StructSymbol)
	return m.Evaluator.EvalSeq(span, m.Func, ss)
}

func (m *EvalInterpretMacro) InvokePar(span *model.Span, arg eval.Arg) (eval.Return, eval.Effect, error) {
	ss := arg.(*symbol.StructSymbol)
	return m.Evaluator.EvalPar(span, m.Func, ss)
}
