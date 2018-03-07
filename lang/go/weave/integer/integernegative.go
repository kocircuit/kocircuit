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
	RegisterGoPkgMacro("integer", "Negative", new(GoIntegerNegativeMacro))
}

type GoIntegerNegativeMacro struct{}

func (m GoIntegerNegativeMacro) MacroID() string { return m.Help() }

func (m GoIntegerNegativeMacro) Label() string { return "negative" }

func (m GoIntegerNegativeMacro) MacroSheathString() *string { return PtrString("integer.Negative") }

func (m GoIntegerNegativeMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerNegativeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerNegative(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer negative")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerNegativeSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerNegativeSoln) String() string { return Sprint(soln) }

func (soln *IntegerNegativeSoln) Cached() *AssignCache { return nil }

func (soln *IntegerNegativeSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerNegativeSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerNegative(span *Span, arg GoStructure) (soln *IntegerNegativeSoln, err error) {
	soln = &IntegerNegativeSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer negative argument")
	}
	v := soln.Arg.Integer.Type
	switch u := v.(type) {
	case Unknown:
		soln.Returns = u
		soln.Form = &GoUnknownForm{}
		return
	case *GoIntegerNumber:
		soln.Returns = u
		soln.Form = &GoInvariantForm{Expr: u.Negative().NumberExpr()}
		return
	case *GoBuiltin:
		switch u.Kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			soln.Returns = u
			soln.Form = &GoNegativeForm{Shaper: soln.Arg.Integer.Extractor}
			return
		}
	case *GoSlice:
		switch w := u.Elem.(type) {
		case Unknown:
			soln.Returns = NewGoSlice(w)
			soln.Form = &GoUnknownForm{}
			return
		case *GoBuiltin:
			switch w.Kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return extendSolveIntegerNegative(span, soln, u)
			}
		}
	case *GoArray:
		switch w := u.Elem.(type) {
		case Unknown:
			soln.Returns = NewGoArray(u.Len, w)
			soln.Form = &GoUnknownForm{}
			return
		case *GoBuiltin:
			switch w.Kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return extendSolveIntegerNegative(span, soln, u)
			}
		}
	}
	return nil, span.Errorf(nil, "integer negative arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerNegativeSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerNegative(span *Span, soln *IntegerNegativeSoln, returns GoType) (_ *IntegerNegativeSoln, err error) {
	var subForm GoSlotForm
	switch returns.(type) {
	case *GoSlice:
		subForm = (*GoIntegerNegativeSliceForm)(soln)
	case *GoArray:
		subForm = (*GoIntegerNegativeArrayForm)(soln)
	default:
		panic("o")
	}
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerNegative",
		Form:    subForm,
		Arg:     soln.Arg.Arg,
		Returns: returns,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = returns
	return soln, nil
}

//===
type GoIntegerNegativeArrayForm IntegerNegativeSoln

func (form *GoIntegerNegativeArrayForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return (*IntegerNegativeSoln)(form).initFormExpr(form.initExpr(), arg...)
}

func (form *GoIntegerNegativeArrayForm) initExpr() GoExpr {
	return &GoColonAssignExpr{ // r := [N]integer{}
		Left:  (*IntegerNegativeSoln)(form).resultExpr(),
		Right: &GoZeroExpr{form.Returns},
	}
}

//===
type GoIntegerNegativeSliceForm IntegerNegativeSoln

func (form *GoIntegerNegativeSliceForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return (*IntegerNegativeSoln)(form).initFormExpr(form.initExpr(), arg...)
}

func (form *GoIntegerNegativeSliceForm) initExpr() GoExpr {
	return &GoColonAssignExpr{ // r := make([]integer, 0, len(argSlice))
		Left: (*IntegerNegativeSoln)(form).resultExpr(),
		Right: &GoCallExpr{
			Func: MakeExpr,
			Arg: []GoExpr{
				&GoTypeRefExpr{form.Returns},
				&GoCallExpr{Func: LenExpr, Arg: []GoExpr{(*IntegerNegativeSoln)(form).sliceExpr()}},
			},
		},
	}
}

//===
func (form *IntegerNegativeSoln) resultExpr() GoExpr { return &GoVerbatimExpr{"r"} }

func (form *IntegerNegativeSoln) sliceExpr() GoExpr { return &GoVerbatimExpr{"s"} }

func (form *IntegerNegativeSoln) initFormExpr(initExpr GoExpr, arg ...*GoSlotExpr) GoExpr {
	root := &GoShapeExpr{
		Shaper: form.Arg.Integer.Extractor,
		Expr:   FindSlotExpr(arg, RootSlot{}),
	}
	indexExpr, elemExpr := &GoVerbatimExpr{"i"}, &GoVerbatimExpr{"e"}
	return &GoBlockExpr{
		Line: []GoExpr{
			&GoColonAssignExpr{ // argSlice = extract(arg)
				Left:  form.sliceExpr(),
				Right: root,
			},
			initExpr,
			&GoForExpr{
				Range: &GoColonAssignExpr{
					Left: &GoListExpr{
						Elem: []GoExpr{indexExpr, elemExpr},
					},
					Right: &GoRangeExpr{form.sliceExpr()},
				},
				Line: []GoExpr{ // r[i] = -e
					&GoAssignExpr{
						Left:  &GoIndexExpr{form.resultExpr(), indexExpr},
						Right: &GoNegativeExpr{elemExpr},
					},
				},
			},
			&GoReturnExpr{form.resultExpr()},
		},
	}
}
