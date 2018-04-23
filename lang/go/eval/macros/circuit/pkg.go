package circuit

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalPkgMacro("circuit", "PkgPath", EvalPkgPathMacro{})
}

type EvalPkgPathMacro struct{}

func (m EvalPkgPathMacro) MacroID() string { return m.Help() }

func (m EvalPkgPathMacro) Label() string { return "pkgpath" }

func (m EvalPkgPathMacro) MacroSheathString() *string { return PtrString("circuit.PkgPath") }

func (m EvalPkgPathMacro) Help() string { return "circuit.PkgPath" }

func (m EvalPkgPathMacro) Doc() string {
	return `The builtin PkgPath function returns the Ko package path of the function invoking PkgPath.`
}

func (EvalPkgPathMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	return BasicStringSymbol(NearestFunc(span).Pkg), nil, nil
}
