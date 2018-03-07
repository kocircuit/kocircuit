package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

var registry = NewGoRegistry()

type GoRegistry struct {
	Registry *Registry `ko:"name=registry"`
}

func NewGoRegistry() *GoRegistry {
	return &GoRegistry{
		Registry: NewRegistry(goGateMacro{}),
	}
}

type goGateMacro struct{}

func (goGateMacro) GateMacro(g gate.Gate) Macro {
	return &GoCallMacro{Valve: MakeValveForGate(g)}
}

func (r *GoRegistry) RegisterGoGate(stub interface{}) {
	r.Registry.RegisterGate(stub)
}

func (r *GoRegistry) RegisterGoGateAt(pkg, name string, stub interface{}) {
	r.Registry.RegisterGateAt(pkg, name, stub)
}

func (r *GoRegistry) RegisterGoMacro(name string, macro Macro) {
	r.Registry.RegisterMacro(name, macro)
}

func (r *GoRegistry) RegisterGoPkgMacro(pkg, name string, macro Macro) {
	r.Registry.RegisterPkgMacro(pkg, name, macro)
}

func (r *GoRegistry) Snapshot() Faculty {
	return r.Registry.Snapshot()
}

func GoFaculty() Faculty {
	return registry.Snapshot()
}

func RegisterGoGate(stub interface{}) {
	registry.RegisterGoGate(stub)
}

func RegisterGoGateAt(pkg, name string, stub interface{}) {
	registry.RegisterGoGateAt(pkg, name, stub)
}

func RegisterGoMacro(name string, macro Macro) {
	registry.RegisterGoMacro(name, macro)
}

func RegisterGoPkgMacro(pkg, name string, macro Macro) {
	registry.RegisterGoPkgMacro(pkg, name, macro)
}
