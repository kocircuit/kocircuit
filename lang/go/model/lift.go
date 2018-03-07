package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func LiftToSequence(span *Span, t GoType) (GoSequence, Shaper, error) {
	simple, simplifier := Simplify(span, t)
	var elem GoType
	switch v := simple.(type) {
	case Unknown:
		u := NewGoUnknown(span)
		return u, &UnknownShaper{
			Shaping: Shaping{Origin: span, From: t, To: u},
		}, nil
	case *GoArray:
		return v, simplifier, nil
	case *GoSlice:
		return v, simplifier, nil
	case *GoPtr:
		elem = v.Elem
	default:
		elem = v
	}
	lifted := NewGoArray(1, elem)
	if lifter, _, err := Assign(span, t, lifted); err != nil {
		return nil, nil, err
		panic("o")
	} else {
		return lifted, lifter, nil
	}
}
