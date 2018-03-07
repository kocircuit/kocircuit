package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoCombiner struct{}

func (GoCombiner) Interpret(_ Evaluator, f *Func) Macro {
	return &GoExpandMacro{Func: f}
}

func (GoCombiner) Combine(
	span *Span,
	f *Func,
	arg Arg,
	returned Return,
	stepResidue StepResidues,
) (Effect, error) {
	circuit := &GoCircuit{
		Origin:  f,
		Comment: fmt.Sprintf("(%s) ko_circuit=%s", f.Syntax.RegionString(), f.FullPath()),
		Valve:   MakeValve(span, f, arg.(GoStructure), returned.(GoType)), // matches WeaveFixedPoint.undeterminedArg
	}
	programEffect := &GoProgramEffect{Circuit: []*GoCircuit{circuit}}
	cached := []*AssignCache{}
	for _, sr := range stepResidue {
		if stepEffect, ok := sr.Effect.(*GoStepEffect); ok {
			if stepEffect.Step != nil {
				circuit.Step = append(circuit.Step, stepEffect.Step)
				cached = append(cached, stepEffect.Step.Cached)
			}
			programEffect = programEffect.Aggregate(stepEffect.ProgramEffect)
			circuit.Effect = circuit.Effect.Aggregate(stepEffect.CircuitEffect)
		} else {
			// local.Select returns nil step effect.
		}
	}
	return &GoCombineEffect{
		Valve:         circuit.Valve,
		Cached:        CompressCacheUnion(cached...),
		ProgramEffect: programEffect,
	}, nil
}
