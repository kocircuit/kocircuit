package eval

import (
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/go/gate"
)

var registry = NewEvalRegistry()

type EvalRegistry struct {
	Registry *eval.Registry `ko:"name=registry"`
}

func NewEvalRegistry() *EvalRegistry {
	return &EvalRegistry{
		Registry: eval.NewRegistry(evalGateMacro{}),
	}
}

type evalGateMacro struct{}

func (evalGateMacro) GateMacro(g gate.Gate) eval.Macro {
	return &EvalCallMacro{Gate: g}
}

func (r *EvalRegistry) RegisterEvalGate(stub interface{}) {
	r.Registry.RegisterGate(stub)
}

func (r *EvalRegistry) RegisterEvalGateAt(pkg, name string, stub interface{}) {
	r.Registry.RegisterGateAt(pkg, name, stub)
}

func (r *EvalRegistry) RegisterEvalMacro(name string, macro eval.Macro) {
	r.Registry.RegisterMacro(name, macro)
}

func (r *EvalRegistry) RegisterEvalPkgMacro(pkg, name string, macro eval.Macro) {
	r.Registry.RegisterPkgMacro(pkg, name, macro)
}

func (r *EvalRegistry) Snapshot() eval.Faculty {
	return r.Registry.Snapshot()
}

func EvalFaculty() eval.Faculty {
	return registry.Snapshot()
}

func RegisterEvalGate(stub interface{}) {
	registry.RegisterEvalGate(stub)
}

func RegisterEvalGateAt(pkg, name string, stub interface{}) {
	registry.RegisterEvalGateAt(pkg, name, stub)
}

func RegisterEvalMacro(name string, macro eval.Macro) {
	registry.RegisterEvalMacro(name, macro)
}

func RegisterEvalPkgMacro(pkg, name string, macro eval.Macro) {
	registry.RegisterEvalPkgMacro(pkg, name, macro)
}
