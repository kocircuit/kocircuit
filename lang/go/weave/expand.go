package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoExpandMacro struct {
	Func *Func `ko:"name=func"`
}

func (macro *GoExpandMacro) Splay() Tree { return Quote{macro.Help()} }

func (macro *GoExpandMacro) MacroID() string { return macro.Help() }

func (macro *GoExpandMacro) Label() string { return "expand" }

func (macro *GoExpandMacro) MacroSheathString() *string { return PtrString("Expand") }

func (macro *GoExpandMacro) Help() string {
	return fmt.Sprintf("Expand(%q.%s)", macro.Func.Pkg, macro.Func.Name)
}

func (macro *GoExpandMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	returns, weaveEffect, err := SpanWeave(span, macro.Func, arg.(GoStructure)) // GoCombineEffect
	if err != nil {
		return nil, nil, err
	}
	span = AugmentSpanCache(span, weaveEffect.Cached)
	callMacro := &GoCallMacro{Valve: weaveEffect.Valve}
	returns2, callEffect, err := callMacro.Invoke(span, arg) // GoMacroEffect
	if err != nil {
		return nil, nil, err
	}
	if returns != returns2 {
		panic("o")
	}
	return returns,
		callEffect.(*GoMacroEffect).
			PlantValve(weaveEffect.Valve).
			AggregateProgramEffect(weaveEffect.ProgramEffect).
			AggregateAssignCache(weaveEffect.Cached),
		nil
}
