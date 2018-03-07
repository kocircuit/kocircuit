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
	RegisterGoPkgMacro("integer", "Less", new(GoIntegerLessMacro))
}

type GoIntegerLessMacro struct{}

func (m GoIntegerLessMacro) MacroID() string { return m.Help() }

func (m GoIntegerLessMacro) Label() string { return "less" }

func (m GoIntegerLessMacro) MacroSheathString() *string { return PtrString("integer.Less") }

func (m GoIntegerLessMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoIntegerLessMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveIntegerLess(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving integer less")
	}
	return soln.Returns, SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type IntegerLessSoln struct {
	Origin    *Span            `ko:"name=origin"`
	Arg       *IntegerSolnArg  `ko:"name=arg"`
	Returns   GoType           `ko:"name=returns"`
	Form      GoSlotForm       `ko:"name=form"`
	Extension *GoCircuitEffect `ko:"name=extension"`
}

func (soln *IntegerLessSoln) String() string { return Sprint(soln) }

func (soln *IntegerLessSoln) Cached() *AssignCache { return nil }

func (soln *IntegerLessSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(soln.Arg.CircuitEffect(), soln.Extension)
}

func (soln *IntegerLessSoln) ProgramEffect() *GoProgramEffect { return nil }

func SolveIntegerLess(span *Span, arg GoStructure) (soln *IntegerLessSoln, err error) {
	soln = &IntegerLessSoln{Origin: span}
	if soln.Arg, err = SolveIntegerArg(span, arg); err != nil {
		return nil, span.Errorf(err, "integer less argument")
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
				return extendSolveIntegerLess(span, soln)
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
				return extendSolveIntegerLess(span, soln)
			}
		}
	}
	return nil, span.Errorf(nil, "integer less arguments must be integral, have %s", Sprint(v))
}

func (soln *IntegerLessSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return soln.Form.FormExpr(arg...)
}

func extendSolveIntegerLess(span *Span, soln *IntegerLessSoln) (_ *IntegerLessSoln, err error) {
	ext := &GoExtendForm{
		Origin:  span,
		Prefix:  "integerLess",
		Form:    (*GoIntegerLessForm)(soln),
		Arg:     soln.Arg.Arg,
		Returns: GoBool,
	}
	soln.Form = ext
	soln.Extension = ext.CircuitEffect()
	soln.Returns = GoBool
	return soln, nil
}

type GoIntegerLessForm IntegerLessSoln

func (form *GoIntegerLessForm) FormExpr(arg ...*GoSlotExpr) GoExpr {
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
						If: &GoDyadicExpr{
							Left: &GoIndexExpr{
								Container: root,
								Index:     &GoDyadicExpr{Left: indexExpr, Op: "-", Right: OneExpr},
							},
							Op:    ">=",
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
