package inject

import (
	"fmt"
	"reflect"

	. "github.com/kocircuit/kocircuit/truss"
)

type InjectPlay struct {
	Gate          interface{}
	InjectPkgPath []string
}

func (g *InjectPlay) Play(ctx *TrussCtx) (result interface{}, err error) {
	gatePkgPath, gateName, err := gateTypeName(g.Gate)
	injProg := &InjectionProgram{
		GatePkgPath:   gatePkgPath,
		GateTypeName:  gateName,
		InjectPkgPath: g.InjectPkgPath,
	}
	injSrc := injProg.Render()
	// determine: code gen location
	// generate injection
	// build injection binary
	XXX
	// create argument file
	// execute injection
	XXX
	// read results file
	XXX
}

func gateTypeName(v interface{}) (pkgPath, typeName string, err error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return "", "", fmt.Errorf("expecting a struct, got %T", v)
	}
	if t.Name() == "" || t.PkgPath() == "" {
		return "", "", fmt.Errorf("expecting a public named type, got %T", v)
	}
	return t.PkgPath(), t.Name(), nil
}
