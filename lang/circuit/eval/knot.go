package eval

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Arg interface {
	String() string
	// Shape
}

type Return interface {
	String() string
	// Shape
}

type Field struct {
	Name   string `ko:"name=name"` // step label or arg name
	Shape  Shape  `ko:"name=shape"`
	Effect Effect `ko:"name=effect"`
	Frame  *Span  `ko:"name=frame"`
}

func (f Field) String() string { return Sprint(f) }

type Knot []Field

func (v Knot) String() string { return Sprint(v) }

func (v Knot) IsEmpty() bool { return len(v) == 0 }

// Select implements Shape.Select.
func (v Knot) Select(span *Span, path Path) (Shape, Effect, error) {
	if len(path) == 0 {
		return v, nil, nil
	}
	projection := v.RestrictTo(path[0])
	switch len(projection) {
	case 0:
		return Empty{}, nil, nil
	case 1:
		return projection[0].Shape.Select(span, path[1:])
	}
	if len(path) > 1 {
		return nil, nil, span.Errorf(nil, "selecting into a sequence")
	}
	return projection, nil, nil
}

// Augment implements Shape.Augment.
func (v Knot) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a knot")
}

// Invoke implements Shape.Invoke.
func (v Knot) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a knot")
}

func (v Knot) Fields() []Field { return v }

func (v Knot) Names() []string {
	n := map[string]bool{}
	r := []string{}
	for _, f := range v {
		if !n[f.Name] {
			n[f.Name] = true
			r = append(r, f.Name)
		}
	}
	return r
}

func (v Knot) FieldGroup() [][]Field {
	r := [][]Field{}
	for _, n := range v.Names() {
		r = append(r, v.RestrictTo(n))
	}
	return r
}

func (v Knot) RestrictTo(name string) Knot {
	r := Knot{}
	for _, f := range v {
		if f.Name == name {
			r = append(r, f)
		}
	}
	return r
}

func (v Knot) StringField(label string) (string, error) {
	g := v.RestrictTo(label)
	if len(g) != 1 {
		return "", fmt.Errorf("not a singleton (got %d) field", len(g))
	}
	s, ok := g[0].Shape.(String)
	if !ok {
		return "", fmt.Errorf("not a string field (type is %T)", g[0].Shape)
	}
	return s.Value_, nil
}
