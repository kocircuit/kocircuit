package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Len", new(GoLenMacro))
}

type GoLenMacro struct{}

func (m GoLenMacro) MacroID() string { return m.Help() }

func (m GoLenMacro) Label() string { return "len" }

func (m GoLenMacro) MacroSheathString() *string { return PtrString("Len") }

func (m GoLenMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoLenMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveLen(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving len")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type LenSolnArg struct {
	Arg    GoStructure   `ko:"name=arg"`
	Series TypeExtractor `ko:"name=series"`
}

func (arg *LenSolnArg) CircuitEffect() *GoCircuitEffect {
	return arg.Series.Extractor.CircuitEffect()
}

func SolveLenArg(span *Span, arg GoStructure) (soln *LenSolnArg, err error) {
	if monadic := StructureMonadicField(arg); monadic == nil {
		return nil, span.Errorf(nil, "len expects a monadic argument, got %s", Sprint(arg))
	} else {
		soln = &LenSolnArg{Arg: arg}
		if soln.Series.Extractor, soln.Series.Type, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg); err != nil {
			return nil, span.Errorf(err, "len expects a monadic argument")
		}
		return soln, nil
	}
}

type LenSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *LenSolnArg      `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
	Optional  bool             `ko:"name=optional"`
}

func (soln *LenSoln) String() string { return Sprint(soln) }

func (soln *LenSoln) Cached() *AssignCache { return nil }

func (soln *LenSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		soln.Arg.CircuitEffect(),
		soln.Extension,
	)
}

func (soln *LenSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveLen(span *Span, arg GoStructure) (soln *LenSoln, err error) {
	soln = &LenSoln{Origin: span}
	if soln.Arg, err = SolveLenArg(span, arg); err != nil {
		return nil, span.Errorf(err, "len argument")
	}
	switch u := soln.Arg.Series.Type.(type) {
	case *GoPtr:
		switch u.Elem.(type) {
		case *GoSlice, *GoArray:
			soln.Optional = true
			return extendSolveLen(span, soln)
		default:
			return nil, span.Errorf(nil, "cannot take length of %s", Sprint(u.Elem))
		}
	case *GoSlice:
		return extendSolveLen(span, soln)
	case *GoArray:
		soln.Form = &GoInvariantForm{&GoIntegerExpr{u.Len}}
		soln.Returns = NewGoIntegerNumber(int64(u.Len))
		return soln, nil
	default:
		return nil, span.Errorf(nil, "len argument not a series: %s", Sprint(u))
	}
}

func extendSolveLen(span *Span, soln *LenSoln) (_ *LenSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "len",
		Form:    (*GoLenForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: GoInt,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = GoInt
	return soln, nil
}

func (soln *LenSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

type GoLenForm LenSoln

func (form *GoLenForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	root := &GoShapeExpr{
		Shaper: form.Arg.Series.Extractor,
		Expr:   FindSlotExpr(arg, RootSlot{}),
	}
	if form.Optional {
		return form.optionalFormExpr(root)
	} else {
		return form.requiredFormExpr(root)
	}
}

func (form *GoLenForm) requiredFormExpr(root GoExpr) GoExpr {
	return &GoReturnExpr{
		&GoCallExpr{ // len(arg)
			Func: LenExpr,
			Arg:  []GoExpr{root},
		},
	}
}

func (form *GoLenForm) optionalFormExpr(root GoExpr) GoExpr {
	return &GoIfThenElseExpr{
		If: &GoInequalityExpr{Left: root, Right: NilExpr}, // arg != nil
		Then: []GoExpr{
			&GoReturnExpr{ // return len(*arg)
				&GoCallExpr{
					Func: LenExpr,
					Arg:  []GoExpr{&GoDerefExpr{root}},
				},
			},
		},
		Else: []GoExpr{&GoReturnExpr{ZeroExpr}}, // return 0
	}
}
