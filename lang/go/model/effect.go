package model

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type GoCircuitEffect struct {
	DuctType []GoType     `ko:"name=duct_type" ctx:"refer"` // auxiliary type definitions
	DuctFunc []GoFuncExpr `ko:"name=duct_func"`             // helper function definitions
}

func (effect *GoCircuitEffect) String() string { return Sprint(effect) }

type CircuitEffector interface {
	CircuitEffect() *GoCircuitEffect
}

func CircuitEffectIfNotNil(effector CircuitEffector) *GoCircuitEffect {
	if effector == nil {
		return nil
	}
	return effector.CircuitEffect()
}

func AggregateCircuitEffects(ee ...*GoCircuitEffect) (agg *GoCircuitEffect) {
	for _, e := range ee {
		agg = agg.Aggregate(e) // e can be nil
	}
	return
}

func (effect *GoCircuitEffect) DuctFuncs() []GoFuncExpr {
	if effect == nil {
		return nil
	}
	return effect.DuctFunc
}

func (effect *GoCircuitEffect) Aggregate(by *GoCircuitEffect) *GoCircuitEffect {
	if effect == nil {
		return by
	}
	if by == nil {
		return effect
	}
	return &GoCircuitEffect{
		DuctType: append(append([]GoType{}, effect.DuctType...), by.DuctType...),
		DuctFunc: append(append([]GoFuncExpr{}, effect.DuctFunc...), by.DuctFunc...),
	}
}

func (effect *GoCircuitEffect) AggregateDuctFunc(helper ...GoFuncExpr) *GoCircuitEffect {
	return &GoCircuitEffect{
		DuctType: effect.CarryType(),
		DuctFunc: append(append([]GoFuncExpr{}, effect.CarryFunc()...), helper...),
	}
}

func (effect *GoCircuitEffect) AggregateDuctType(typ ...GoType) *GoCircuitEffect {
	return &GoCircuitEffect{
		DuctType: append(append([]GoType{}, effect.CarryType()...), typ...),
		DuctFunc: effect.CarryFunc(),
	}
}

func (effect *GoCircuitEffect) CarryFunc() []GoFuncExpr {
	if effect == nil {
		return nil
	}
	return effect.DuctFunc
}

func (effect *GoCircuitEffect) CarryType() []GoType {
	if effect == nil {
		return nil
	}
	return effect.DuctType
}
