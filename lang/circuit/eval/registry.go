package eval

import (
	"log"
	"reflect"
	"sync"

	"github.com/kocircuit/kocircuit/lang/go/gate"
)

type Registry struct {
	sync.Mutex `ko:"name=mutex"`
	GateMacro  `ko:"name=gateMacro"`
	Faculty    Faculty `ko:"name=faculty"`
}

type GateMacro interface {
	GateMacro(gate.Gate) Macro
}

func NewRegistry(gateMacro GateMacro) *Registry {
	return &Registry{
		GateMacro: gateMacro,
		Faculty:   Faculty{},
	}
}

func (r *Registry) RegisterGate(stub interface{}) {
	gate, err := gate.BindGate(reflect.TypeOf(stub))
	if err != nil {
		panic(err)
	}
	r.RegisterPkgMacro(gate.GoPkgPath(), gate.GoName(), r.GateMacro.GateMacro(gate))
}

func (r *Registry) RegisterGateAt(pkg, name string, stub interface{}) {
	gate, err := gate.BindGate(reflect.TypeOf(stub))
	if err != nil {
		panic(err)
	}
	r.RegisterPkgMacro(pkg, name, r.GateMacro.GateMacro(gate))
}

func (r *Registry) RegisterMacro(name string, macro Macro) {
	r.RegisterPkgMacro("", name, macro)
}

func (r *Registry) RegisterPkgMacro(pkg, name string, macro Macro) {
	ideal := Ideal{Pkg: pkg, Name: name}
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Faculty[ideal]; ok {
		log.Fatalf("re-registering eval macro %v", ideal)
	}
	r.Faculty[ideal] = macro
}

func (r *Registry) Snapshot() Faculty {
	r.Lock()
	defer r.Unlock()
	g := Faculty{}
	for k, v := range r.Faculty {
		g[k] = v
	}
	return g
}
