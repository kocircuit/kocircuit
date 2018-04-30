package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Panic", new(EvalPanicMacro))
	RegisterEvalMacro("Recover", new(EvalRecoverMacro))
}

type EvalPanicMacro struct{}

func (m EvalPanicMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalPanicMacro) MacroID() string { return m.Help() }

func (m EvalPanicMacro) Label() string { return "panic" }

func (m EvalPanicMacro) MacroSheathString() *string { return PtrString("Panic") }

func (m EvalPanicMacro) Help() string { return "Panic" }

func (m EvalPanicMacro) Doc() string {
	return `
The builtin function Panic throws a panic up the stack.
If the panic is not handled (using Recover, covered below),
it results in program termination with a message to the log output.
The panic message includes the stack of the location that caused
the panic, as well as a user-supplied panic value.

Panic accepts any number of named arguments (with any values).
Panic creates a structure whose fields are the named arguments 
passed to it, and attaches it as a _panic value_ to the panic itself.
(Panic values can be retrieved using Recover, as explained in the next section.)

Panic never returns into the function calling it.`
}

func (EvalPanicMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	join := EvalJoinMacro{}
	if returns, effect, err = join.Invoke(span, arg); err != nil {
		return nil, nil, err
	} else {
		panic(NewEvalPanic(span, returns.(Symbol)))
	}
}

type EvalRecoverMacro struct{}

func (m EvalRecoverMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalRecoverMacro) MacroID() string { return m.Help() }

func (m EvalRecoverMacro) Label() string { return "recover" }

func (m EvalRecoverMacro) MacroSheathString() *string { return PtrString("Recover") }

func (m EvalRecoverMacro) Help() string { return "Recover" }

func (m EvalRecoverMacro) Doc() string {
	return `
The builtin function Recover provides a mechanism for handling panics,
caused by runtime conditions.

Recover expects two arguments, invoke and panic, both of which must be varieties (functional values).

Recover starts by invoking the functional value invoke (without passing any arguments):

* If the invocation succeeds in returning a value without panicking, Recover will return that value.

* If the invocation panics, Recover will capture the panic and invoke
the functional value of the panic argument, while alsp passing the panic value
as a default (aka monadic) argument to panic. Whatever the call to panic returns,
will be returned by Recover.`
}

// Recover(invoke, panic)
func (EvalRecoverMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol)
	invokeVty, ok := a.Walk("invoke").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "recover invoke is not a variety")
	}
	panicVty, ok := a.Walk("panic").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "recover panic is not a variety")
	}
	defer func() {
		if r := recover(); r != nil {
			evalPanic := r.(*EvalPanic)
			panicFields := Fields{{Name: "", Shape: evalPanic.Panic, Effect: nil, Frame: span}}
			returns, effect, err = panicVty.Evoke(span, panicFields)
			return
		}
	}()
	return invokeVty.Invoke(span)
}
