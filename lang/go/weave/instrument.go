package weave

import (
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoInstrument struct {
	Valve         *GoValve         `ko:"name=valve"`
	Returns       GoType           `ko:"name=returns"`
	Circuit       []*GoCircuit     `ko:"name=circuit"`
	Directive     []*GoDirective   `ko:"name=directive"`
	ProgramEffect *GoProgramEffect `ko:"name=programEffect"`
}

type GoPkgStats struct {
	TotalFunc   int64
	TotalStep   int64
	TotalType   int64
	StepPerFunc float64
	TypePerFunc float64
}

func (inst *GoInstrument) Stats() *GoPkgStats {
	stats := &GoPkgStats{}
	seen := map[string]bool{}
	for _, circuit := range inst.Circuit {
		if !seen[circuit.Valve.TypeID()] {
			stats.TotalFunc += 1
			stats.TotalStep += int64(len(circuit.Step))
			stats.TotalType += int64(len(circuit.DuctType()))
			seen[circuit.Valve.TypeID()] = true
		}
	}
	if stats.TotalFunc > 0 {
		stats.StepPerFunc = float64(stats.TotalStep) / float64(stats.TotalFunc)
		stats.TypePerFunc = float64(stats.TotalType) / float64(stats.TotalFunc)
	}
	return stats
}
