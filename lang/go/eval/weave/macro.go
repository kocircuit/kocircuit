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
	return `Weave plays the user function func, using a user-supplied evaluation logic ...

weaver is a structure of the form below, where "->" indicates the return type of a function:

	(
		reserve: (pkg, name)
		Enter: (step, ctx, arg) -> (returns, effect)
		Leave: (step, ctx, arg) -> (returns, effect)
		Figure: (step, ctx, literal) -> (returns, effect)
		Link: (step, ctx, arg, name, monadic) -> (returns, effect)
		Select: (step, ctx, arg, path) -> (returns, effect)
		Augment: (step, ctx, arg, fields) -> (returns, effect)
		Invoke: (step, ctx, arg) -> (returns, effect)
		Combine: (ctx, arg, stepResults) -> (returns, effect)
	)

The named arguments (above) have the following types:

	arg is a user object
	ctx is a user object
	returns is a user object
	effect is a user object
	pkg is a string
	name is a string

* step is (...)
* fields is a sequence of (name, arg) pairs
* literal is XXX
* stepResults is a sequence of (step, returns, panic, effect)
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
