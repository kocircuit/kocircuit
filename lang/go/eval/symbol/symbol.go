package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Symbol interface {
	Shape
	Splayer
	Type() Type
	Hash() string
	Equal(Symbol) bool
}

var (
	symbolPtr    *Symbol
	typeOfSymbol = reflect.TypeOf((*Symbol)(nil)).Elem()
)

type Symbols []Symbol

func (syms Symbols) Types() Types {
	types := make(Types, len(syms))
	for i, sym := range syms {
		types[i] = sym.Type()
	}
	return types
}

type Type interface {
	Splayer
	IsType()
}

type Types []Type
