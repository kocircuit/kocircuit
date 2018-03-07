package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (local *GoLocal) Invoke(span *Span) (Shape, Effect, error) {
	span = AugmentSpanCache(span, local.Cached) // merge span and upstream caches
	switch v := local.Image().(type) {
	case GoVarietal:
		return local.invokeVarietal(span, v)
	}
	return nil, nil, span.Errorf(nil, "invoking a %v", local)
}

func (local *GoLocal) invokeVarietal(span *Span, vty GoVarietal) (Shape, Effect, error) {
	inv, err := GoInvoke(span, vty)
	if err != nil {
		return nil, nil, err
	}
	step := &GoStep{
		Span:  span,
		Label: NearestStep(span).Label,
		Arrival: []*GoArrival{
			ArrivalFromLocal(span, local, RootSlot{}),
		},
		Returns: inv.Returns,
		Logic:   inv.Logic,
		Cached:  inv.Cached(),
	}
	return local.Inherit(span, step, inv.Returns),
		&GoStepEffect{
			Step:          step,
			CircuitEffect: inv.CircuitEffect(),
			ProgramEffect: inv.ProgramEffect(),
		}, nil
}

type Invocation struct {
	Varietal       GoVarietal      `ko:"name=varietal"`
	Returns        GoType          `ko:"name=returns"`
	Logic          GoStepLogic     `ko:"name=logic"`
	Transformation *Transformation `ko:"name=transformation"`
}

func (inv *Invocation) Cached() *AssignCache {
	if inv == nil {
		return nil
	}
	return inv.Transformation.Cached
}

func (inv *Invocation) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return inv.Logic.FormExpr(arg...)
}

func (inv *Invocation) CircuitEffect() *GoCircuitEffect {
	if inv == nil {
		return nil
	}
	return inv.Transformation.CircuitEffect
}

func (inv *Invocation) ProgramEffect() *GoProgramEffect {
	if inv == nil {
		return nil
	}
	return inv.Transformation.ProgramEffect
}

func GoInvoke(span *Span, varietal GoVarietal) (inv *Invocation, err error) {
	if _, err = VerifyUniqueFieldAugmentations(span, varietal); err != nil {
		return nil, span.Errorf(err, "invoking augmentation")
	}
	inv = &Invocation{Varietal: varietal}
	if inv.Transformation, err = VarietalMacroTransforms(span, varietal); err != nil {
		return nil, err
	}
	branchReturns := make([]GoType, len(inv.Transformation.MacroTransform))
	for i := range inv.Transformation.MacroTransform {
		branchReturns[i] = inv.Transformation.MacroTransform[i].Returns
	}
	if inv.Returns, err = GeneralizeSequence(span, branchReturns...); err != nil {
		return nil, err
	}
	span = AugmentSpanCache(span, inv.Cached())
	invCtx := NewInvokingCtx(span, varietal, inv.Transformation.MacroTransform, inv.Returns)
	projected, projector := VarietalProject(span, varietal)
	// MustAssign(span, projector.Shape(vty), projected)
	inv.Logic = &GoInvokeLogic{
		DuctFunc:  invCtx.FuncExpr(),
		Projector: projector,
	}
	inv.Transformation.CircuitEffect = inv.Transformation.CircuitEffect.
		Aggregate(invCtx.CircuitEffect()).
		Aggregate(projector.CircuitEffect()).
		AggregateDuctType(projected)
	return
}

type Transformation struct {
	MacroTransform []*GoMacroTransform `ko:"name=macroTransform"`
	CircuitEffect  *GoCircuitEffect    `ko:"name=circuitEffect"`
	ProgramEffect  *GoProgramEffect    `ko:"name=programEffect"`
	Cached         *AssignCache        `ko:"name=cached"`
}

func VarietalMacroTransforms(span *Span, varietal GoVarietal) (*Transformation, error) {
	t := &Transformation{}
	cached := []*AssignCache{}
	t.CircuitEffect = t.CircuitEffect.AggregateDuctType(VarietalProjectionReal(varietal))
	line := VarietalProjectionLine(varietal)
	t.MacroTransform = make([]*GoMacroTransform, len(line))
	terminal := len(line) == 1
	for i, line := range line {
		switch m := line.Macro.(type) {
		case *GoUnknownMacro:
			t.MacroTransform[i] = &GoMacroTransform{
				Origin:   span,
				Macro:    m,
				Arg:      line.Real(),
				Returns:  NewGoUnknown(span),
				SlotForm: &GoUnknownSlotForm{},
			}
		default:
			returns, effect, err := m.Invoke(
				RefineMacro(RefineOutline(span, fmt.Sprintf("%d", i)), m),
				line.Real(),
			)
			if err != nil {
				if terminal {
					return nil, span.Errorf(err, "%s", Sprint(line.Macro))
				} else {
					return nil, err
				}
			}
			macroEffect := effect.(*GoMacroEffect)
			cached = append(cached, macroEffect.Cached)
			t.CircuitEffect = t.CircuitEffect.Aggregate(macroEffect.CircuitEffect)
			t.ProgramEffect = t.ProgramEffect.Aggregate(macroEffect.ProgramEffect)
			t.MacroTransform[i] = &GoMacroTransform{
				Origin:      span,
				Macro:       line.Macro,
				Arg:         macroEffect.Arg,
				Returns:     returns.(GoType),
				SlotForm:    macroEffect.SlotForm,
				ExpandValve: macroEffect.ExpandValve,
			}
		}
	}
	t.Cached = AssignCacheUnion(cached...)
	return t, nil
}
