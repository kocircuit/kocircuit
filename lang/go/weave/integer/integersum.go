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
	RegisterGoPkgMacro("integer", "Sum", new(GoIntegerSumMacro))
}

type GoIntegerSumMacro struct{}

func (m GoIntegerSumMacro) MacroID() string { return m.Help() }

func (m GoIntegerSumMacro) Label() string { return "sum" }

func (m GoIntegerSumMacro) MacroSheathString() *string { return PtrString("integer.Sum") }

func (m GoIntegerSumMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerSumMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerSum(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer sum")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerSumSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerSumSoln) String() string { return Sprint(soln) }

func (soln *IntegerSumSoln) Cached() *AssignCache { return nil }

func (soln *IntegerSumSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerSumSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerSum(span *Span, arg GoStructure) (soln *IntegerSumSoln, err error) {
	soln = &IntegerSumSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer sum argument")
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
				return extendSolveIntegerSum(span, soln, w)
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
				return extendSolveIntegerSum(span, soln, w)
			}
		}
	}
	return nil, span.Errorf(nil, "integer sum arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerSumSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerSum(span *Span, soln *IntegerSumSoln, returns GoType) (_ *IntegerSumSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerSum",
		Form:    (*GoIntegerSumForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: returns,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = returns
	return soln, nil
}

type GoIntegerSumForm IntegerSumSoln

func (form *GoIntegerSumForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	root := &GoShapeExpr{
		Shaper: form.Arg.Integer.Extractor,
		Expr:   FindSlotExpr(arg, RootSlot{}),
	}
	sumExpr := &GoVerbatimExpr{"sum"}
	elemExpr := &GoVerbatimExpr{"elem"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoVarDeclExpr{Name: sumExpr, Type: form.Returns},
			&GoForExpr{
				Range: &GoColonAssignExpr{
					Left: &GoListExpr{
						Elem: []GoExpr{UnderlineExpr, elemExpr},
					},
					Right: &GoRangeExpr{root},
				},
				Line: []GoExpr{
					&GoDyadicExpr{
						Left: sumExpr, Op: "+=", Right: elemExpr,
					},
				},
			},
			&GoReturnExpr{sumExpr},
		},
	}
}
