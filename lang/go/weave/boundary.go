package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoBoundary struct{}

func (GoBoundary) Figure(span *Span, figure Figure) (Shape, Effect, error) {
	switch u := figure.(type) {
	case Bool:
		return NewGoLocal(span, nil, NewGoBoolNumber(u.Value_)), nil, nil
	case Integer:
		return NewGoLocal(span, nil, NewGoIntegerNumber(u.Value_)), nil, nil
	case Float:
		return NewGoLocal(span, nil, NewGoFloatNumber(u.Value_)), nil, nil
	case String:
		return NewGoLocal(span, nil, NewGoStringNumber(u.Value_)), nil, nil
	case Macro:
		vty := NewGoVariety(span, u, nil)
		projected, projector := VarietalProject(span, vty)
		return NewGoLocal(span, nil, vty), &GoStepEffect{
			CircuitEffect: projector.CircuitEffect().AggregateDuctType(vty, projected),
		}, nil
	}
	panic("unknown figure")
}

func (GoBoundary) Enter(span *Span, arg Arg) (Shape, Effect, error) {
	real := arg.(GoStructure)
	step := &GoStep{
		Span:    span,
		Label:   NearestStep(span).Label,
		Returns: real,
		Logic:   &GoEnterLogic{},
	}
	return NewGoLocal(span, step, real), &GoStepEffect{Step: step}, nil
}

func (GoBoundary) Leave(span *Span, shape Shape) (Return, Effect, error) {
	local := shape.(*GoLocal)
	step := &GoStep{ // applies selections before returning a concrete value
		Span:   span,
		Label:  NearestStep(span).Label,
		Result: true, // side mechanism to retrieve final value
		Arrival: []*GoArrival{
			ArrivalFromLocal(span, local, RootSlot{}),
		},
		Returns: local.Image(),
		Logic:   &GoLeaveLogic{},
	}
	return step.Returns, &GoStepEffect{Step: step}, nil
}

func ArrivalFromLocal(span *Span, local *GoLocal, slot Slot) *GoArrival {
	if local.Step == nil {
		switch u := local.Type.(type) {
		case GoNumber:
			return &GoArrival{
				FromExpr: u.NumberExpr(),
				Slot:     slot,
				Shaper:   local.Shaper,
			}
		case GoVarietal:
			return &GoArrival{
				FromExpr: u.VarietyExpr(),
				Slot:     slot,
				Shaper:   local.Shaper,
			}
		default:
			panic("o")
		}
	} else {
		ch := fmt.Sprintf("from_%s_to_%s_step_slot_%s",
			local.Step.Label,
			NearestStep(span).Label,
			slot.Label(),
		)
		local.Step.Send = append(local.Step.Send, ch) // update send side
		return &GoArrival{
			FromChan: &GoReceiveExpr{FromChan: ch},
			Shaper:   local.Shaper,
			Slot:     slot,
		}
	}
}

type GoZeroForm struct {
	Type GoType `ko:"name=type"`
}

func (zero *GoZeroForm) FormExpr(...*GoSlotExpr) GoExpr {
	return &GoZeroExpr{zero.Type}
}
