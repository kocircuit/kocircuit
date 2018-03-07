package weave

import (
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type EffectorSlotForm interface {
	GoEffect // Cached + CircuitEffect + ProgramEffect
	GoSlotForm
}

func SlotFormMacroEffect(arg GoStructure, form EffectorSlotForm) *GoMacroEffect {
	return &GoMacroEffect{
		Arg:           arg.(GoStructure),
		SlotForm:      form,
		Cached:        form.Cached(),
		CircuitEffect: form.CircuitEffect(),
		ProgramEffect: form.ProgramEffect(),
	}
}
