package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Repeated", new(GoRepeatedMacro))
}

type GoRepeatedMacro struct{}

func (m GoRepeatedMacro) MacroID() string { return m.Help() }

func (m GoRepeatedMacro) Label() string { return "repeated" }

func (m GoRepeatedMacro) MacroSheathString() *string { return PtrString("Repeated") }

func (m GoRepeatedMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoRepeatedMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if f, err := solveRepeatedMacro(span, arg); err != nil {
		return nil, nil, err
	} else {
		return f.Returns,
			SlotFormMacroEffect(arg.(GoStructure), f),
			nil
	}
}

func solveRepeatedMacro(span *Span, arg Arg) (f *convertSlotForm, err error) {
	if monadic := StructureMonadicField(arg.(GoStructure)); monadic == nil {
		return nil, span.Errorf(nil, "repeated expects a monadic argument, got %s", Sprint(arg))
	} else {
		f = &convertSlotForm{}
		if f.Extractor, f.Number, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg.(GoStructure)); err != nil {
			return nil, span.Errorf(err, "repated expects a monadic argument")
		}
		f.Returns = NewGoSlice(f.Number)
		assn := NewAssignCtx(span)
		if f.Converter, err = assn.Assign(f.Number, f.Returns); err != nil {
			return nil, err
		}
		return
	}
}

type convertSlotForm struct {
	Number    GoType `ko:"name=number"`
	Returns   GoType `ko:"name=returns"`
	Extractor Shaper `ko:"name=extractor"`
	Converter Shaper `ko:"name=shaper"`
}

func (f *convertSlotForm) String() string { return Sprint(f) }

func (f *convertSlotForm) Cached() *AssignCache { return nil }

func (f *convertSlotForm) ProgramEffect() *GoProgramEffect { return nil }

func (f *convertSlotForm) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		f.Extractor.CircuitEffect(),
		f.Converter.CircuitEffect(),
	)
}

func (f *convertSlotForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoShapeExpr{
		Shaper: f.Converter,
		Expr: &GoShapeExpr{
			Shaper: f.Extractor,
			Expr:   FindSlotExpr(arg, RootSlot{}),
		},
	}
}
