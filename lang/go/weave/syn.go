package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type Join struct {
	Variety   *GoVariety `ko:"name=variety"`
	Evocation *Evocation `ko:"name=evocation"`
}

func (join *Join) Cached() *AssignCache {
	if join == nil {
		return nil
	}
	return join.Evocation.Cached()
}

func (join *Join) CircuitEffect() *GoCircuitEffect {
	if join == nil {
		return nil
	}
	return join.Evocation.CircuitEffect().AggregateDuctType(join.Variety)
}

func (join *Join) ProgramEffect() *GoProgramEffect {
	if join == nil {
		return nil
	}
	return join.Evocation.ProgramEffect()
}

func GoJoin(span *Span, with []*GoAugmentField) (join *Join, err error) {
	join = &Join{
		Variety: NewGoVariety(span, GoJoinMacro{}, nil),
	}
	if join.Evocation, err = GoEvoke(span, join.Variety, with); err != nil {
		return nil, err
	}
	return
}

func (join *Join) Returns() GoType { return join.Evocation.Returns() }

func (join *Join) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return join.Evocation.FormExpr(
		append(
			[]*GoSlotExpr{{
				Slot: RootSlot{},
				Expr: join.Variety.VarietyExpr(), // build variety zero value
			}},
			arg...,
		)...,
	)
}

type Evocation struct {
	Augmentation *Augmentation `ko:"name=augmentation"`
	Invocation   *Invocation   `ko:"name=invocation"`
}

func (evo *Evocation) Cached() *AssignCache {
	if evo == nil {
		return nil
	}
	return evo.Invocation.Cached()
}

func GoEvoke(span *Span, vty GoVarietal, with []*GoAugmentField) (evoke *Evocation, err error) {
	evoke = &Evocation{}
	if evoke.Augmentation, err = GoAugment(span, vty, with); err != nil {
		return nil, err
	}
	if evoke.Invocation, err = GoInvoke(span, evoke.Augmentation.Varietal.(*GoVariety)); err != nil {
		return nil, err
	}
	return
}

func (evoke *Evocation) Returns() GoType { return evoke.Invocation.Returns }

func (evoke *Evocation) CircuitEffect() *GoCircuitEffect {
	if evoke == nil {
		return nil
	}
	return AggregateCircuitEffects(
		CircuitEffectIfNotNil(evoke.Augmentation),
		CircuitEffectIfNotNil(evoke.Invocation),
	)
}

func (evoke *Evocation) ProgramEffect() *GoProgramEffect {
	if evoke == nil {
		return nil
	}
	return AggregateProgramEffects(
		ProgramEffectIfNotNil(evoke.Augmentation),
		ProgramEffectIfNotNil(evoke.Invocation),
	)
}

func (evoke *Evocation) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoSlotFormExpr{
		SlotExpr: []*GoSlotExpr{{
			Slot: RootSlot{},
			Expr: &GoSlotFormExpr{
				SlotExpr: arg,
				Form:     evoke.Augmentation,
			},
		}},
		Form: evoke.Invocation,
	}
}

type TypeExtractor struct {
	Name      string `ko:"name=name"`
	Type      GoType `ko:"name=type"`
	Extractor Shaper `ko:"name=extractor"`
}

type VarietalExtractor struct {
	Varietal  GoVarietal `ko:"name=varietal"`
	Extractor Shaper     `ko:"name=extractor"`
}
