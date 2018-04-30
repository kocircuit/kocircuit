package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Translate", new(EvalTranslateMacro))
}

type EvalTranslateMacro struct{}

func (m EvalTranslateMacro) MacroID() string { return m.Help() }

func (m EvalTranslateMacro) Label() string { return "translate" }

func (m EvalTranslateMacro) MacroSheathString() *string { return PtrString("Translate") }

func (m EvalTranslateMacro) Help() string {
	return "Translate(translator, func, ctx, arg) -> (returns, effect)"
}

func (m EvalTranslateMacro) Doc() string {
	return `XXX

* translator is a structure of the form, where "->" indicates the return type of a variety:

	(
		Enter(step, ctx, arg) -> (returns, effect)
		Leave(step, ctx, arg) -> (returns, effect)
		Figure(step, ctx, literal) -> (returns, effect)
		SelectArg(step, ctx, arg, label, monadic) -> (returns, effect)
		Select(step, ctx, arg, path) -> (returns, effect)
		Augment(step, ctx, arg, fields) -> (returns, effect)
		Invoke(step, ctx, arg) -> (returns, effect)
		Combine(ctx, arg, stepResults) -> (returns, effect)
	)

* arg is a user object
* ctx is a user object
* returns is a user object
* effect is a user object

* step is (...)
* fields is a sequence of (name, arg) pairs
* literal is XXX
* stepResults is a sequence of (step, returned, panic, effect)
`
}

func (EvalTranslateMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	// extract translator
	translator, err := ExtractTranslator(span, a.Walk("translator"))
	if err != nil {
		return nil, nil, err
	}
	_ = translator
	// extract function
	fu, ok := a.Walk("func").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "func must be a variety, got %v", a.Walk("func"))
	}
	interpretFunc, ok := fu.Macro.(*EvalInterpretMacro)
	if !ok {
		return nil, nil, span.Errorf(nil, "func must be a user function (not a builtin), got %v", fu.Macro)
	}
	_ = interpretFunc
	// extract ctx and arg
	ctxArg := a.Walk("ctx")
	_ = ctxArg
	argArg := a.Walk("arg")
	_ = argArg
	panic("o")
}
