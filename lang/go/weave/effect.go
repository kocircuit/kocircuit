package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type GoEffect interface {
	String() string
	Cached() *AssignCache
	CircuitEffect() *GoCircuitEffect
	ProgramEffect() *GoProgramEffect
}

type NoEffect struct{}

func (NoEffect) String() string { return Sprint(NoEffect{}) }

func (NoEffect) Cached() *AssignCache { return nil }

func (NoEffect) CircuitEffect() *GoCircuitEffect { return nil }

func (NoEffect) ProgramEffect() *GoProgramEffect { return nil }

type GoEffectSlotForm struct {
	GoEffect   `ko:"name=effect"`
	GoSlotForm `ko:"name=slotForm"`
}

// GoProgramEffects float all the way up to the top of the weaving recursion.
type GoProgramEffect struct {
	Circuit   []*GoCircuit   `ko:"name=circuit"` // dependent circuits (including associated type definitions)
	Directive []*GoDirective `ko:"name=directive"`
	// weaving-internal
	Recurrence  []*GoArgRecurrence `ko:"name=recurrence"` // pass-thru subordinate recurrences
	WeavingStat *GoWeavingStat     `ko:"name=weavingStat"`
}

func (effect *GoProgramEffect) String() string { return Sprint(effect) }

type ProgramEffector interface {
	ProgramEffect() *GoProgramEffect
}

func ProgramEffectIfNotNil(effector ProgramEffector) *GoProgramEffect {
	if effector == nil {
		return nil
	}
	return effector.ProgramEffect()
}

func (effect *GoProgramEffect) SubtractRecurrences(f *Func) (
	recurrence []*GoArgRecurrence,
	remainder *GoProgramEffect,
) {
	if effect == nil {
		return nil, nil
	}
	rest := []*GoArgRecurrence{}
	for _, g := range effect.Recurrence {
		if g.Func == f {
			recurrence = append(recurrence, g)
		} else {
			rest = append(rest, g)
		}
	}
	return recurrence, &GoProgramEffect{
		Circuit:     effect.Circuit,
		Directive:   effect.Directive,
		Recurrence:  rest,
		WeavingStat: effect.WeavingStat,
	}
}

func (effect *GoProgramEffect) AggregateWeavingStat(with *GoWeavingStat) *GoProgramEffect {
	if effect == nil {
		return &GoProgramEffect{WeavingStat: with}
	}
	return &GoProgramEffect{
		Circuit:     effect.Circuit,
		Directive:   effect.Directive,
		Recurrence:  effect.Recurrence,
		WeavingStat: SumGoWeavingStat(effect.WeavingStat, with),
	}
}

func (effect *GoProgramEffect) AggregateDirective(with ...*GoDirective) *GoProgramEffect {
	if effect == nil {
		return &GoProgramEffect{Directive: with}
	}
	return &GoProgramEffect{
		Circuit:     effect.Circuit,
		Directive:   append(append([]*GoDirective{}, effect.Directive...), with...),
		Recurrence:  effect.Recurrence,
		WeavingStat: effect.WeavingStat,
	}
}

func AggregateProgramEffects(ee ...*GoProgramEffect) (agg *GoProgramEffect) {
	for _, e := range ee {
		agg = agg.Aggregate(e)
	}
	return
}

func (effect *GoProgramEffect) Aggregate(by *GoProgramEffect) *GoProgramEffect {
	if effect == nil {
		return by
	}
	if by == nil {
		return effect
	}
	return &GoProgramEffect{
		Circuit:     append(append([]*GoCircuit{}, effect.Circuit...), by.Circuit...),
		Directive:   append(append([]*GoDirective{}, effect.Directive...), by.Directive...),
		Recurrence:  append(append([]*GoArgRecurrence{}, effect.Recurrence...), by.Recurrence...),
		WeavingStat: SumGoWeavingStat(effect.WeavingStat, by.WeavingStat),
	}
}

// GoCombineEffect is emitted by the GoCombiner.
type GoCombineEffect struct {
	Valve         *GoValve         `ko:"name=valve"` // valve for top-level circuit
	ProgramEffect *GoProgramEffect `ko:"name=programEffect"`
	Cached        *AssignCache     `ko:"name=cached"`
}

func (effect *GoCombineEffect) String() string { return Sprint(effect) }

func (effect *GoCombineEffect) AggregateWeavingStat(with *GoWeavingStat) *GoCombineEffect {
	return &GoCombineEffect{
		Valve:         effect.Valve,
		Cached:        effect.Cached,
		ProgramEffect: effect.ProgramEffect.AggregateWeavingStat(with),
	}
}

type GoArgRecurrence struct {
	Func *Func       `ko:"name=func"`
	Arg  GoStructure `ko:"name=arg"`
	Over *GoWeaveCtx `ko:"name=over"`
}

func (effect *GoCombineEffect) SubtractRecurrences(f *Func) (
	recurrence []*GoArgRecurrence,
	remainder *GoCombineEffect,
) {
	recurrence, programRemainder := effect.ProgramEffect.SubtractRecurrences(f)
	return recurrence,
		&GoCombineEffect{
			Valve:         effect.Valve,
			Cached:        effect.Cached,
			ProgramEffect: programRemainder,
		}
}

// GoStepEffect is emitted by individual circuit steps.
type GoStepEffect struct {
	Step *GoStep `ko:"name=step"` // if non-nil, a step implementation
	// aggregation
	CircuitEffect *GoCircuitEffect `ko:"name=circuitEffect"` // circuit-level aggregation
	ProgramEffect *GoProgramEffect `ko:"name=programEffect"` // program-level aggregation
}

func (effect *GoStepEffect) String() string { return Sprint(effect) }

type GoMacroEffect struct {
	Arg         GoType     `ko:"name=arg"`
	SlotForm    GoSlotForm `ko:"name=slotForm"`
	ExpandValve *GoValve   `ko:"name=expandValve"` // used only by GoExpandMacro
	// aggregation
	Cached        *AssignCache     `ko:"name=cached"`
	CircuitEffect *GoCircuitEffect `ko:"name=circuitEffect"` // circuit-level aggregation
	ProgramEffect *GoProgramEffect `ko:"name=programEffect"` // program-level aggregation
}

func (effect *GoMacroEffect) String() string { return Sprint(effect) }

func (macro *GoMacroEffect) PlantValve(valve *GoValve) *GoMacroEffect {
	if macro == nil {
		return &GoMacroEffect{ExpandValve: valve}
	}
	return &GoMacroEffect{
		Arg:           macro.Arg,
		SlotForm:      macro.SlotForm,
		ExpandValve:   valve,
		Cached:        macro.Cached,
		CircuitEffect: macro.CircuitEffect,
		ProgramEffect: macro.ProgramEffect,
	}
}

func (macro *GoMacroEffect) AggregateProgramEffect(by *GoProgramEffect) *GoMacroEffect {
	if macro == nil {
		return &GoMacroEffect{ProgramEffect: by}
	}
	return &GoMacroEffect{
		Arg:           macro.Arg,
		SlotForm:      macro.SlotForm,
		ExpandValve:   macro.ExpandValve,
		Cached:        macro.Cached,
		CircuitEffect: macro.CircuitEffect,
		ProgramEffect: macro.ProgramEffect.Aggregate(by),
	}
}

func (macro *GoMacroEffect) AggregateCircuitEffect(by *GoCircuitEffect) *GoMacroEffect {
	if macro == nil {
		return &GoMacroEffect{CircuitEffect: by}
	}
	return &GoMacroEffect{
		Arg:           macro.Arg,
		SlotForm:      macro.SlotForm,
		ExpandValve:   macro.ExpandValve,
		Cached:        macro.Cached,
		CircuitEffect: macro.CircuitEffect.Aggregate(by),
		ProgramEffect: macro.ProgramEffect,
	}
}

func (macro *GoMacroEffect) AggregateAssignCache(cache *AssignCache) *GoMacroEffect {
	if macro == nil {
		return &GoMacroEffect{Cached: cache}
	}
	return &GoMacroEffect{
		Arg:           macro.Arg,
		SlotForm:      macro.SlotForm,
		ExpandValve:   macro.ExpandValve,
		Cached:        AssignCacheUnion(macro.Cached, cache),
		CircuitEffect: macro.CircuitEffect,
		ProgramEffect: macro.ProgramEffect,
	}
}

func (macro *GoMacroEffect) AggregateDirective(directive ...*GoDirective) *GoMacroEffect {
	if macro == nil {
		return &GoMacroEffect{
			ProgramEffect: &GoProgramEffect{
				Directive: directive,
			},
		}
	}
	return &GoMacroEffect{
		Arg:           macro.Arg,
		SlotForm:      macro.SlotForm,
		ExpandValve:   macro.ExpandValve,
		Cached:        macro.Cached,
		CircuitEffect: macro.CircuitEffect,
		ProgramEffect: macro.ProgramEffect.AggregateDirective(directive...),
	}
}
