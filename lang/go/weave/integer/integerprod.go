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
	RegisterGoPkgMacro("integer", "Prod", new(GoIntegerProdMacro))
}

type GoIntegerProdMacro struct{}

func (m GoIntegerProdMacro) MacroID() string { return m.Help() }

func (m GoIntegerProdMacro) Label() string { return "prod" }

func (m GoIntegerProdMacro) MacroSheathString() *string { return PtrString("integer.Prod") }

func (m GoIntegerProdMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerProdMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerProd(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer prod")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerProdSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerProdSoln) String() string { return Sprint(soln) }

func (soln *IntegerProdSoln) Cached() *AssignCache { return nil }

func (soln *IntegerProdSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerProdSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerProd(span *Span, arg GoStructure) (soln *IntegerProdSoln, err error) {
	soln = &IntegerProdSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer prod argument")
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
				return extendSolveIntegerProd(span, soln, w)
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
				return extendSolveIntegerProd(span, soln, w)
			}
		}
	}
	return nil, span.Errorf(nil, "integer prod arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerProdSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerProd(span *Span, soln *IntegerProdSoln, returns GoType) (_ *IntegerProdSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerProd",
		Form:    (*GoIntegerProdForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: returns,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = returns
	return soln, nil
}

type GoIntegerProdForm IntegerProdSoln

func (form *GoIntegerProdForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	root := &GoShapeExpr{
		Shaper: form.Arg.Integer.Extractor,
		Expr:   FindSlotExpr(arg, RootSlot{}),
	}
	prodExpr := &GoVerbatimExpr{"prod"}
	elemExpr := &GoVerbatimExpr{"elem"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoAssignExpr{
				Left:  &GoVarDeclExpr{Name: prodExpr, Type: form.Returns},
				Right: OneExpr,
			},
			&GoForExpr{
				Range: &GoColonAssignExpr{
					Left: &GoListExpr{
						Elem: []GoExpr{UnderlineExpr, elemExpr},
					},
					Right: &GoRangeExpr{root},
				},
				Line: []GoExpr{
					&GoDyadicExpr{
						Left: prodExpr, Op: "*=", Right: elemExpr,
					},
				},
			},
			&GoReturnExpr{prodExpr},
		},
	}
}
