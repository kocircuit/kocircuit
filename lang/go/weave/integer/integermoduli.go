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
	RegisterGoPkgMacro("integer", "Moduli", new(GoIntegerModuliMacro))
}

type GoIntegerModuliMacro struct{}

func (m GoIntegerModuliMacro) MacroID() string { return m.Help() }

func (m GoIntegerModuliMacro) Label() string { return "moduli" }

func (m GoIntegerModuliMacro) MacroSheathString() *string { return PtrString("integer.Moduli") }

func (m GoIntegerModuliMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerModuliMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerModuli(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer moduli")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerModuliSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerModuliSoln) String() string { return Sprint(soln) }

func (soln *IntegerModuliSoln) Cached() *AssignCache { return nil }

func (soln *IntegerModuliSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerModuliSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerModuli(span *Span, arg GoStructure) (soln *IntegerModuliSoln, err error) {
	soln = &IntegerModuliSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer moduli argument")
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
				return extendSolveIntegerModuli(span, soln, w)
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
				return extendSolveIntegerModuli(span, soln, w)
			}
		}
	}
	return nil, span.Errorf(nil, "integer moduli arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerModuliSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerModuli(span *Span, soln *IntegerModuliSoln, returns GoType) (_ *IntegerModuliSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerModuli",
		Form:    (*GoIntegerModuliForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: returns,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = returns
	return soln, nil
}

type GoIntegerModuliForm IntegerModuliSoln

func (form *GoIntegerModuliForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	seriesExpr := &GoVerbatimExpr{"s"}
	moduliExpr := &GoVerbatimExpr{"m"}
	indexExpr := &GoVerbatimExpr{"i"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoColonAssignExpr{ // s := extract(arg)
				Left: seriesExpr,
				Right: &GoShapeExpr{
					Shaper: form.Arg.Integer.Extractor,
					Expr:   FindSlotExpr(arg, RootSlot{}),
				},
			},
			&GoAssignExpr{ // var m integer = s[0]
				Left:  &GoVarDeclExpr{Name: moduliExpr, Type: form.Returns},
				Right: &GoIndexExpr{seriesExpr, ZeroExpr},
			},
			&GoForExpr{
				Range: &GoIncrementExpr{
					Zero: &GoColonAssignExpr{indexExpr, OneExpr},
					Invariant: &GoDyadicExpr{
						Left:  indexExpr,
						Op:    "<",
						Right: &GoCallExpr{Func: LenExpr, Arg: []GoExpr{seriesExpr}},
					},
					Increment: &GoMonadicExpr{Left: indexExpr, Op: "++"},
				},
				Line: []GoExpr{
					&GoDyadicExpr{
						Left: moduliExpr, Op: "%=", Right: &GoIndexExpr{seriesExpr, indexExpr},
					},
				},
			},
			&GoReturnExpr{moduliExpr},
		},
	}
}
