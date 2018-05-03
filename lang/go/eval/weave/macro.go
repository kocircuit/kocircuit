package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Weave", new(EvalWeaveMacro))
}

type EvalWeaveMacro struct{}

func (m EvalWeaveMacro) MacroID() string { return m.Help() }

func (m EvalWeaveMacro) Label() string { return "weave" }

func (m EvalWeaveMacro) MacroSheathString() *string { return PtrString("Weave") }

func (m EvalWeaveMacro) Help() string {
	return "Weave(weaver, func, ctx, arg) -> (returns, effect)"
}

func (m EvalWeaveMacro) Doc() string {
	return `Weave executes the circuit func on argument arg,
using a circuit evaluation logic given by argument weaver,
as well as an evaluation context ctx.

weaver is a structure of the form below,
where "->" indicates the return type of a function:

	(
		reserve: sequence of (pkg, name) pairs
		Enter: function (stepCtx, object) -> (returns, effect)
		Leave: function (stepCtx, object) -> (returns, effect)
		Link: function (stepCtx, object, name, monadic) -> (returns, effect)
		Select: function (stepCtx, object, path) -> (returns, effect)
		Augment: function (stepCtx, object, fields) -> (returns, effect)
		Invoke: function (stepCtx, object) -> (returns, effect)
		Literal: function (stepCtx, figure) -> (returns, effect)
		Combine: function (summaryCtx, stepResidues) -> (returns, effect)
	)

The named arguments (above) have the following types:

	object is a any type
	ctx is any type
	returns is any type
	effect is any type
	pkg is a string
	name is a string

stepCtx is a structure of the form (pkg, func, step, logic, source, ctx).
summaryCtx is a structure of the form (pkg, func, source, ctx, arg, returns).
fields is a sequence of structures of the form (name, monadic, objects).
`
}

func (EvalWeaveMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	// parse weaver
	weaver, err := ParseWeaver(span, a.Walk("weaver"))
	if err != nil {
		return nil, nil, err
	}
	// parse function
	fu, ok := a.Walk("func").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "func must be a variety, got %v", a.Walk("func"))
	}
	interpretFunc, ok := fu.Macro.(*EvalInterpretMacro)
	if !ok {
		return nil, nil, span.Errorf(nil, "func must be a user function (not a builtin), got %v", fu.Macro)
	}
	// parse ctx and arg
	ctxArg, argArg := a.Walk("ctx"), a.Walk("arg")
	// weave
	weave := &Weave{
		Idiom:  interpretFunc.Evaluator.EvalIdiom(),
		Repo:   interpretFunc.Evaluator.EvalRepo(),
		Weaver: weaver,
		Func:   interpretFunc.Func,
		Ctx:    ctxArg,
		Arg:    argArg,
	}
	return weave.Play(span)
}
