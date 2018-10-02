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

// evalGateMacro is a helper to create a Macro from a Gate.
type evalGateMacro struct{}

func (evalGateMacro) GateMacro(g gate.Gate) eval.Macro {
	return &EvalCallMacro{Gate: g}
}

// RegisterEvalGate registers an evaluation gate with the name and
// package of the given stub.
func (r *EvalRegistry) RegisterEvalGate(stub interface{}) {
	r.Registry.RegisterGate(stub)
}

// RegisterNamedEvalGate registers an evaluation gate with the given name in
// the package of the given stub.
func (r *EvalRegistry) RegisterNamedEvalGate(name string, stub interface{}) {
	r.Registry.RegisterNamedGate(name, stub)
}

// RegisterEvalGateAt registers an evaluation gate with given name in the given package.
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

// RegisterEvalGate registers an evaluation gate in the default registry
// with the name and package of the given stub.
func RegisterEvalGate(stub interface{}) {
	registry.RegisterEvalGate(stub)
}

// RegisterNamedEvalGate registers an evaluation gate in the default registry
// with the given name in the package of the given stub.
func RegisterNamedEvalGate(name string, stub interface{}) {
	registry.RegisterNamedEvalGate(name, stub)
}

// RegisterEvalGateAt registers an evaluation gate in the default registry
// with given name in the given package.
func RegisterEvalGateAt(pkg, name string, stub interface{}) {
	registry.RegisterEvalGateAt(pkg, name, stub)
}

func RegisterEvalMacro(name string, macro eval.Macro) {
	registry.RegisterEvalMacro(name, macro)
}

func RegisterEvalPkgMacro(pkg, name string, macro eval.Macro) {
	registry.RegisterEvalPkgMacro(pkg, name, macro)
}
