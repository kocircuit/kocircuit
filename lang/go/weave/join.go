package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("", new(GoJoinMacro))
}

// GoJoinMacro is a Ko operator that joins its arguments into a structure.
type GoJoinMacro struct{}

func (m GoJoinMacro) MacroID() string { return m.Help() }

func (m GoJoinMacro) Label() string { return "join" }

func (m GoJoinMacro) MacroSheathString() *string { return PtrString("Join") }

func (m GoJoinMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoJoinMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	simple, simplifier := Simplify(span, arg.(GoType))
	if empty, ok := simple.(*GoEmpty); ok { // if simplifies to empty, join materializes the simplification
		return empty,
			&GoMacroEffect{
				Arg:      arg.(GoType),
				SlotForm: &GoShapeForm{simplifier},
			}, nil
	} else if monadic, ok := FindMonadicGoSelection(span, arg.(GoType)); ok { // if argument is monadic, e.g. (1, 2, 3)
		if selector, selected, err := GoSelect(span, Path{monadic}, arg.(GoType)); err != nil {
			panic("o")
		} else {
			return selected,
				&GoMacroEffect{
					Arg:      arg.(GoType),
					SlotForm: &GoShapeForm{Shaper: selector},
				}, nil
		}
	} else { // if no monadic argument is present, e.g. (a: 1, a: 2) or ()
		return arg,
			&GoMacroEffect{
				Arg:      arg.(GoType),
				SlotForm: &GoIdentityForm{},
			}, nil
	}
}

func FindMonadicGoSelection(span *Span, ps GoType) (string, bool) {
	simple, _ := Simplify(span, ps)
	if s, ok := simple.(*GoStruct); ok {
		if f := StructureMonadicField(s); f != nil {
			return f.KoName(), true
		} else {
			return "", false
		}
	} else {
		return "", false
	}
}
