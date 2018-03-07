package integer

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterGoPkgMacro("integer", "Equal", new(GoIntegerEqualMacro))
}

type GoIntegerEqualMacro struct{}

func (m GoIntegerEqualMacro) MacroID() string { return m.Help() }

func (m GoIntegerEqualMacro) Label() string { return "equal" }

func (m GoIntegerEqualMacro) MacroSheathString() *string { return PtrString("integer.Equal") }

func (m GoIntegerEqualMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerEqualMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerEqual(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer equal")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerEqualSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerEqualSoln) String() string { return Sprint(soln) }

func (soln *IntegerEqualSoln) Cached() *AssignCache { return nil }

func (soln *IntegerEqualSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerEqualSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerEqual(span *Span, arg GoStructure) (soln *IntegerEqualSoln, err error) {
	soln = &IntegerEqualSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer equal argument")
	}
	v := soln.Arg.Integer.Type
	switch u := v.(type) {
	case Unknown:
		soln.Returns = u
		soln.Form = &GoUnknownForm{}
		return
	case *GoSlice:
		switch w := u.Elem.(type) {
		case Unknown:
			soln.Returns = w
			soln.Form = &GoUnknownForm{}
			return
		case *GoBuiltin:
			switch w.Kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return extendSolveIntegerEqual(span, soln)
			}
		}
	case *GoArray:
		switch w := u.Elem.(type) {
		case Unknown:
			soln.Returns = w
			soln.Form = &GoUnknownForm{}
			return
		case *GoBuiltin:
			switch w.Kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				return extendSolveIntegerEqual(span, soln)
			}
		}
	}
	return nil, span.Errorf(nil, "integer equal arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerEqualSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerEqual(span *Span, soln *IntegerEqualSoln) (_ *IntegerEqualSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerEqual",
		Form:    (*GoIntegerEqualForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: GoBool,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = GoBool
	return soln, nil
}

type GoIntegerEqualForm IntegerEqualSoln

func (form *GoIntegerEqualForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	root := &GoVerbatimExpr{"series"}
	indexExpr := &GoVerbatimExpr{"i"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoColonAssignExpr{
				Left: root,
				Right: &GoShapeExpr{
					Shaper: form.Arg.Integer.Extractor,
					Expr:   FindSlotExpr(arg, RootSlot{}),
				},
			},
			&GoForExpr{
				Range: &GoIncrementExpr{
					Zero: &GoColonAssignExpr{indexExpr, OneExpr},
					Invariant: &GoDyadicExpr{
						Left:  indexExpr,
						Op:    "<",
						Right: &GoCallExpr{Func: LenExpr, Arg: []GoExpr{root}},
					},
					Increment: &GoMonadicExpr{Left: indexExpr, Op: "++"},
				},
				Line: []GoExpr{
					&GoIfThenExpr{
						If: &GoInequalityExpr{
							Left: &GoIndexExpr{
								Container: root,
								Index:     &GoDyadicExpr{Left: indexExpr, Op: "-", Right: OneExpr},
							},
							Right: &GoIndexExpr{Container: root, Index: indexExpr},
						},
						Then: []GoExpr{
							&GoReturnExpr{FalseExpr},
						},
					},
				},
			},
			&GoReturnExpr{TrueExpr},
		},
	}
}
