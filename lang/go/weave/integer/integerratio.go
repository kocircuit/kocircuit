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
	RegisterGoPkgMacro("integer", "Ratio", new(GoIntegerRatioMacro))
}

type GoIntegerRatioMacro struct{}

func (m GoIntegerRatioMacro) MacroID() string { return m.Help() }

func (m GoIntegerRatioMacro) Label() string { return "ratio" }

func (m GoIntegerRatioMacro) MacroSheathString() *string { return PtrString("integer.Ratio") }

func (m GoIntegerRatioMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerRatioMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerRatio(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(nil, "solving integer ratio")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerRatioSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerRatioSoln) String() string { return Sprint(soln) }

func (soln *IntegerRatioSoln) Cached() *AssignCache { return nil }

func (soln *IntegerRatioSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerRatioSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerRatio(span *Span, arg GoStructure) (soln *IntegerRatioSoln, err error) {
	soln = &IntegerRatioSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer ratio argument")
	}
	v := soln.Arg.Integer.Type
	switch u := v.(type) {
	case Unknown:
		soln.Returns = u
		soln.Form = &GoUnknownForm{}
		return
	case *GoIntegerNumber:
		soln.Returns = u
		soln.Form = &GoInvariantForm{Expr: u.NumberExpr()}
		return
	case *GoBuiltin:
		switch u.Kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			soln.Returns = u
			soln.Form = &GoIdentityForm{}
			return
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			soln.Returns = u
			soln.Form = &GoIdentityForm{}
			return
		}
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
				return extendSolveIntegerRatio(span, soln, w)
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
				return extendSolveIntegerRatio(span, soln, w)
			}
		}
	}
	return nil, span.Errorf(nil, "integer ratio arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerRatioSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerRatio(span *Span, soln *IntegerRatioSoln, returns GoType) (_ *IntegerRatioSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerRatio",
		Form:    (*GoIntegerRatioForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: returns,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = returns
	return soln, nil
}

type GoIntegerRatioForm IntegerRatioSoln

func (form *GoIntegerRatioForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	root := &GoVerbatimExpr{"series"}
	indexExpr := &GoVerbatimExpr{"i"}
	ratioExpr := &GoVerbatimExpr{"ratio"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoColonAssignExpr{
				Left: root,
				Right: &GoShapeExpr{
					Shaper: form.Arg.Integer.Extractor,
					Expr:   FindSlotExpr(arg, RootSlot{}),
				},
			},
			&GoAssignExpr{
				Left:  &GoVarDeclExpr{Name: ratioExpr, Type: form.Returns},
				Right: &GoIndexExpr{Container: root, Index: ZeroExpr},
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
					&GoDyadicExpr{
						Left: ratioExpr, Op: "/=", Right: &GoIndexExpr{Container: root, Index: indexExpr},
					},
				},
			},
			&GoReturnExpr{ratioExpr},
		},
	}
}
