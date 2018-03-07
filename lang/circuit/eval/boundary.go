package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type Boundary interface {
	Figure(*Span, Figure) (Shape, Effect, error)
	Enter(*Span, Arg) (Shape, Effect, error)
	Leave(*Span, Shape) (Return, Effect, error)
}

type IdentityBoundary struct{}

func (IdentityBoundary) Figure(_ *Span, figure Figure) (Shape, Effect, error) {
	switch u := figure.(type) {
	case Macro:
		return Variety{Macro: u}, nil, nil
	case Shape:
		return u, nil, nil
	}
	panic("unknown figure")
}

func (IdentityBoundary) Enter(_ *Span, arg Arg) (Shape, Effect, error) {
	if arg == nil {
		return nil, nil, nil
	}
	return arg.(Shape), nil, nil
}

func (IdentityBoundary) Leave(_ *Span, shape Shape) (Return, Effect, error) {
	return shape, nil, nil
}
