package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Boot", new(EvalBootMacro))
}

type EvalBootMacro struct{}

func (m EvalBootMacro) MacroID() string { return m.Help() }

func (m EvalBootMacro) Label() string { return "boot" }

func (m EvalBootMacro) MacroSheathString() *string { return PtrString("Boot") }

func (m EvalBootMacro) Help() string {
	return "Boot(booter, func, ctx, arg) -> (returns, effect)"
}

func (m EvalBootMacro) Doc() string {
	return `Boot plays the user function func, using a user-supplied evaluation logic ...

booter is a structure of the form below, where "->" indicates the return type of a function:

	(
		reserve: (pkg, name) // pkg and name are strings
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

* step is (...)
* fields is a sequence of (name, arg) pairs
* literal is XXX
* stepResults is a sequence of (step, returned, panic, effect)
`
}

func (EvalBootMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	// parse booter
	booter, err := ParseBooter(span, a.Walk("booter"))
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
	// boot
	boot := &Boot{
		Idiom:  interpretFunc.Evaluator.EvalIdiom(),
		Repo:   interpretFunc.Evaluator.EvalRepo(),
		Booter: booter,
		Func:   interpretFunc.Func,
		Ctx:    ctxArg,
		Arg:    argArg,
	}
	return boot.Play(span)
}
