package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

// Shapes are the figures (Integer, Float, String, Bool, Empty, Variety) and Fields.
type Shape interface {
	String() string
	Select(*Span, Path) (Shape, Effect, error)
	Link(*Span, string, bool) (Shape, Effect, error)
	Augment(*Span, Fields) (Shape, Effect, error)
	Invoke(*Span) (Shape, Effect, error)
}

type Effect interface {
	String() string
}
