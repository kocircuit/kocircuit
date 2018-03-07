package macros

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Sequential", new(EvalSequentialMacro))
	RegisterEvalMacro("Parallel", new(EvalParallelMacro))
}

type EvalSequentialMacro struct{}

func (m EvalSequentialMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalSequentialMacro) MacroID() string { return m.Help() }

func (m EvalSequentialMacro) Label() string { return "sequential" }

func (m EvalSequentialMacro) MacroSheathString() *string { return PtrString("Sequential") }

func (m EvalSequentialMacro) Help() string {
	return "Sequential"
}

func (EvalSequentialMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol).SelectMonadic()
	if vty, ok := a.(*VarietySymbol); !ok {
		return nil, nil, span.Errorf(nil, "sequential cannot be applied to a non-variety %v", a)
	} else {
		seqVty := MakeVarietySymbol(
			evalSequentialMacro{vty.Macro},
			vty.Arg,
		)
		return seqVty.Invoke(span)
	}
}

type evalSequentialMacro struct {
	Macro Macro `ko:"name=macro"`
}

func (m evalSequentialMacro) Splay() Tree { return Quote{m.Help()} }

func (m evalSequentialMacro) MacroID() string { return m.Help() }

func (m evalSequentialMacro) Label() string { return "evalsequential" }

func (m evalSequentialMacro) MacroSheathString() *string { return nil }

func (m evalSequentialMacro) Help() string {
	return fmt.Sprintf("Sequential[%s]", m.Macro.Help())
}

func (m evalSequentialMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	switch macro := m.Macro.(type) {
	case *EvalInterpretMacro:
		return macro.InvokeSeq(span, arg)
	default:
		return macro.Invoke(span, arg)
	}
}

type EvalParallelMacro struct{}

func (m EvalParallelMacro) Splay() Tree { return Quote{m.Help()} }

func (m EvalParallelMacro) MacroID() string { return m.Help() }

func (m EvalParallelMacro) Label() string { return "parallel" }

func (m EvalParallelMacro) MacroSheathString() *string { return PtrString("Parallel") }

func (m EvalParallelMacro) Help() string {
	return "Parallel"
}

func (EvalParallelMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	a := arg.(*StructSymbol).SelectMonadic()
	if vty, ok := a.(*VarietySymbol); !ok {
		return nil, nil, span.Errorf(nil, "parallel cannot be applied to a non-variety %v", a)
	} else {
		parVty := MakeVarietySymbol(
			evalParallelMacro{vty.Macro},
			vty.Arg,
		)
		return parVty.Invoke(span)
	}
}

type evalParallelMacro struct {
	Macro Macro `ko:"name=macro"`
}

func (m evalParallelMacro) Splay() Tree { return Quote{m.Help()} }

func (m evalParallelMacro) MacroID() string { return m.Help() }

func (m evalParallelMacro) Label() string { return "evalparallel" }

func (m evalParallelMacro) MacroSheathString() *string { return nil }

func (m evalParallelMacro) Help() string {
	return fmt.Sprintf("Parallel[%s]", m.Macro.Help())
}

func (m evalParallelMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	switch macro := m.Macro.(type) {
	case *EvalInterpretMacro:
		return macro.InvokePar(span, arg)
	default:
		return macro.Invoke(span, arg)
	}
}
