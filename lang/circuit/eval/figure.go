package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Figure interface{}

type Empty struct{}

func (e Empty) String() string { return Sprint(e) }

func (e Empty) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to empty")
}

func (e Empty) Select(span *Span, path Path) (Shape, Effect, error) {
	return e, nil, nil
}

func (e Empty) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting an empty")
}

func (e Empty) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking an empty")
}

type Integer struct{ Value_ int64 }

func (v Integer) String() string { return Sprint(v) }

func (v Integer) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to integer")
}

func (v Integer) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into integer", path)
	}
	return v, nil, nil
}

func (v Integer) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting an integer")
}

func (v Integer) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking an integer")
}

type Float struct{ Value_ float64 }

func (v Float) String() string { return Sprint(v) }

func (v Float) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to float")
}

func (v Float) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into float", path)
	}
	return v, nil, nil
}

func (v Float) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a float")
}

func (v Float) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a float")
}

type Bool struct{ Value_ bool }

func (v Bool) String() string { return Sprint(v) }

func (v Bool) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to bool")
}

func (v Bool) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into bool", path)
	}
	return v, nil, nil
}

func (v Bool) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a bool")
}

func (v Bool) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a bool")
}

type String struct{ Value_ string }

func (v String) String() string { return Sprint(v) }

func (v String) Link(span *Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to string")
}

func (v String) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into string", path)
	}
	return v, nil, nil
}

func (v String) Augment(span *Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a string")
}

func (v String) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a string")
}
