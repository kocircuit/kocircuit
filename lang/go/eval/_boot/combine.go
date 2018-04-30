package boot

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) Combine(
	span *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidue StepResidues,
) (Effect, error) {
	XXX //XXX
	return nil, nil
}

func (b *BootController) Interpret(_ Evaluator, f *Func) Macro {
	return &BootInterpretMacro{Func: f}
}

type BootInterpretMacroXXX struct {
	Func *Func `ko:"name=func"`
}

// InterpretFunc communicates to Variety.Disassemble the underlying function identity.
func (m *BootInterpretMacro) InterpretFunc() (pkgPath, funcName string) {
	return m.Func.Pkg, m.Func.Name
}

func (m *BootInterpretMacro) Splay() Tree {
	return Quote{m.Help()}
}

func (m *BootInterpretMacro) MacroID() string { return m.Help() }

func (m *BootInterpretMacro) Label() string { return "eval" }

func (m *BootInterpretMacro) MacroSheathString() *string { return nil }

func (m *BootInterpretMacro) Help() string {
	return fmt.Sprintf("%s", m.Func.FullPath())
}

func (m *BootInterpretMacro) Doc() string {
	return m.Func.DocLong()
}

func (m *BootInterpretMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	return m.InvokeSeq(span, arg) // default circuit execution mode
}

func (m *BootInterpretMacro) InvokeSeq(span *Span, arg Arg) (Return, Effect, error) {
	ss := arg.(*StructSymbol)
	return m.Evaluator.EvalSeq(span, m.Func, ss)
}

func (m *BootInterpretMacro) InvokePar(span *Span, arg Arg) (Return, Effect, error) {
	ss := arg.(*StructSymbol)
	return m.Evaluator.EvalPar(span, m.Func, ss)
}
