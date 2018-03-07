package model

import (
	"fmt"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/gate"
)

// GoValve captures the interface (name, arguments and return value) of a circuit in Go.
// It is used to render weaved circuit definitions in Go, as well as
// calls to imported Go types with circuit interface.
type GoValve struct {
	Origin  *Span       `ko:"name=origin"`
	Address *GoAddress  `ko:"name=address"`
	Arg     GoStructure `ko:"name=arg"`
	Returns GoType      `ko:"name=returns"`
}

func (valve *GoValve) TypeID() string {
	return valve.Real().TypeID() // also captures gate valves
}

func (valve *GoValve) Real() GoType {
	return NewGoNeverNilPtr(valve.alias())
}

func (valve *GoValve) alias() *GoAlias {
	return NewGoAlias(valve.Address, FlattenGoStructure(valve.Arg))
}

func (valve *GoValve) RenderDef(fileCtx GoFileContext) string {
	return valve.alias().RenderDef(fileCtx)
}

func (valve *GoValve) RenderTypeName(fileCtx GoFileContext) string {
	return valve.Address.RenderExpr(fileCtx)
}

func (valve *GoValve) RenderNew(fileCtx GoFileContext) string {
	return fmt.Sprintf("new(%s)", valve.RenderTypeName(fileCtx))
}

func MakeValveForGate(g gate.Gate) *GoValve {
	arg := make([]*GoField, len(g.Arg()))
	for i, f := range g.Arg() {
		arg[i] = &GoField{
			Name: f.GoName(),
			Type: GraftGoType(f.Type),
			Tag:  KoTags(f.KoName(), f.IsMonadic()),
		}
	}
	return &GoValve{
		Address: &GoAddress{
			Comment: fmt.Sprintf("%v", g),
			GroupPath: GoGroupPath{
				Group: GoHereditaryPkgGroup,
				Path:  g.GoPkgPath(),
			},
			Name: g.GoName(),
		},
		Arg:     NewGoStruct(arg...),
		Returns: GraftGoType(g.Returns()),
	}
}

func MakeValve(span *Span, f *Func, arg GoStructure, returns GoType) *GoValve {
	arg = FlattenGoStructure(arg)
	addr := &GoAddress{
		Comment:   fmt.Sprintf("source=%s", f.Syntax.RegionString()),
		Span:      span,
		GroupPath: SpanGroupPath(span),
		Name:      GoCircuitName(span, f.Name, arg), // dependence on span and arg ensures uniqueness
	}
	return &GoValve{Origin: span, Address: addr, Arg: arg, Returns: returns}
}

// Circuit names must uniquely depend on the: (1) argument structure and (2) span.
// This ensures that circuit names act as exact content hashes for the
// semantics of the corresponding go circuit implementation.
func GoCircuitName(span *Span, name string, arg GoStructure) string {
	if strings.ToLower(name[:1]) == name[:1] {
		// make the package private functions in Ko public in Go
		return fmt.Sprintf("N_%s_arg_%s_span_%s", name, arg.TypeID(), span.SpanID().String())
	}
	return fmt.Sprintf("%s_arg_%s_span_%s", name, arg.TypeID(), span.SpanID().String())
}

func FlattenGoStructure(v GoStructure) GoStructure {
	return NewGoStruct(v.StructureField()...)
}
