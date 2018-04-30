package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Variety is a gate shape.
type Variety struct {
	Macro Macro
	Arg   Arg
}

func (v Variety) Doc() string { return Sprint(v) }

func (v Variety) String() string { return Sprint(v) }

func (v Variety) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to variety")
}

func (v Variety) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into variety", path)
	}
	return v, nil, nil
}

func (v Variety) Augment(span *Span, arg Knot) (Shape, Effect, error) {
	aug := Knot{} // copy-and-append arg
	if v.Arg != nil {
		aug = append(aug, v.Arg.(Knot)...)
	}
	aug = append(aug, arg...)
	return Variety{Macro: v.Macro, Arg: aug}, nil, nil
}

func (v Variety) Invoke(span *Span) (Shape, Effect, error) {
	r, eff, err := v.Macro.Invoke(span, v.Arg)
	return r.(Shape), eff, err
}
