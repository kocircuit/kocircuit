package flow

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type Envelope interface {
	Enter(*Span) (Flow, error)                       // Enter returns the output of the Enter step of the enclosing function.
	Make(*Span, interface{}) (Flow, error)           // Make ...
	MakePkgFunc(*Span, string, string) (Flow, error) // MakePkgFunc ...
	MakeOp(*Span, []string) (Flow, error)            // MakeOp ...
}

type Flow interface {
	Link(*Span, string, bool) (Flow, error)    // Link ...
	Select(*Span, []string) (Flow, error)      // Select ...
	Augment(*Span, []GatherFlow) (Flow, error) // Augment ...
	Invoke(*Span) (Flow, error)                // Invoke ...
	Leave(*Span) (Flow, error)                 // Leave ...
}

type GatherFlow struct {
	Field string
	Flow  Flow
}

func PlaySeqFlow(span *Span, f *Func, env Envelope) (Flow, []Flow, error) {
	p := &flowPlayer{span: RefineFunc(span, f), fun: f, env: env}
	stepReturn, err := playSeq(f, p)
	if err != nil {
		return nil, nil, err
	}
	return stepReturnFlow(f, stepReturn)
}

func PlayParFlow(span *Span, f *Func, env Envelope) (Flow, []Flow, error) {
	p := &flowPlayer{span: RefineFunc(span, f), fun: f, env: env}
	stepReturn, err := playPar(f, p)
	if err != nil {
		return nil, nil, err
	}
	return stepReturnFlow(f, stepReturn)
}

func stepReturnFlow(f *Func, stepReturn map[*Step]Edge) (Flow, []Flow, error) {
	stepFlow := []Flow{}
	for _, step := range f.Step {
		stepFlow = append(stepFlow, stepReturn[step].(Flow))
	}
	return stepReturn[f.Leave].(Flow), stepFlow, nil
}
