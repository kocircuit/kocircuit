package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Spin", new(GoSpinMacro))
}

type GoSpinMacro struct{}

func (m GoSpinMacro) MacroID() string { return m.Help() }

func (m GoSpinMacro) Label() string { return "spin" }

func (m GoSpinMacro) MacroSheathString() *string { return PtrString("Spin") }

func (m GoSpinMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoSpinMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if monadic := StructureMonadicField(arg.(GoStructure)); monadic == nil {
		return nil, nil, span.Errorf(nil, "spin expects a monadic argument, got %s", Sprint(arg))
	} else {
		spin := &GoSpin{Origin: span}
		selector, selected, err := GoSelect(span, Path{monadic.KoName()}, arg.(GoStructure))
		if err != nil {
			return nil, nil, span.Errorf(err, "spin expects a monadic argument")
		}
		switch vty := selected.(type) {
		case Unknown:
			return NewGoEmpty(span), &GoMacroEffect{Arg: arg.(GoStructure)}, nil
		case *GoVariety:
			spin.Select = selector
			if spin.Invoke, err = GoInvoke(span, vty); err != nil {
				return nil, nil, span.Errorf(err, "spin invoking payload")
			}
			return NewGoEmpty(span),
				&GoMacroEffect{
					Arg:           arg.(GoStructure),
					SlotForm:      spin,
					CircuitEffect: spin.CircuitEffect(),
					ProgramEffect: spin.ProgramEffect(),
					Cached:        spin.Invoke.Cached(),
				},
				nil
		default:
			return nil, nil, span.Errorf(nil, "spin argument must be a variety, got %s", Sprint(selected))
		}
	}
}

type GoSpin struct {
	Origin *Span       `ko:"name=origin"`
	Select Shaper      `ko:"name=select"`
	Invoke *Invocation `ko:"name=invoke"`
}

func (spin *GoSpin) CircuitEffect() *GoCircuitEffect {
	return spin.Invoke.CircuitEffect().
		Aggregate(spin.Select.CircuitEffect()).
		AggregateDuctFunc(spin.FuncExpr())
}

func (spin *GoSpin) ProgramEffect() *GoProgramEffect {
	return spin.Invoke.ProgramEffect()
}

func (spin *GoSpin) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: spin.Origin,
		Name:   fmt.Sprintf("spin_%s", spin.Origin.SpanID().String()),
	}
}

func (spin *GoSpin) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: spin.FuncExpr(),
		Arg: []GoExpr{
			&GoVerbatimExpr{"step_ctx"},
			&GoShapeExpr{
				Shaper: spin.Select,
				Expr:   FindSlotExpr(arg, RootSlot{}),
			},
		},
	}
}

// FuncExpr returns the spin_ function definition:
// func spin_xyz567(step_ctx *runtime.Context, payload Variety) Empty {
// 	go func() {
//		INVOKE(step_ctx, payload)
//	}()
//	return struct{}{/*empty*/}
// }
func (spin *GoSpin) FuncExpr() GoFuncExpr {
	payloadExpr := &GoVerbatimExpr{"payload"}
	return &GoDuctFuncExpr{
		Comment:    fmt.Sprintf("(%s)", spin.Origin.SourceLine()),
		FuncName:   spin.FuncName(),
		ArgName:    payloadExpr,
		ArgType:    spin.Invoke.Varietal,
		ReturnType: NewGoEmpty(nil),
		Line: []GoExpr{
			&GoGoFuncExpr{
				Line: []GoExpr{
					&GoSlotFormExpr{
						SlotExpr: []*GoSlotExpr{{Slot: RootSlot{}, Expr: payloadExpr}},
						Form:     spin.Invoke.Logic,
					},
				},
			},
			&GoReturnExpr{&GoZeroExpr{NewGoEmpty(nil)}},
		},
	}
}
