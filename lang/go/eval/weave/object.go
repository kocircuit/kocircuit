package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func (b *WeaveController) Wrap(symbol Symbol) Shape {
	return WeaveObject{Controller: b, Object: symbol}
}

func (b *WeaveController) UnwrapArg(arg Arg) Symbol {
	return arg.(WeaveObject).Object
}

func (b *WeaveController) Unwrap(shape Shape) Symbol {
	return shape.(WeaveObject).Object
}

// WeaveObject is a Shape.
type WeaveObject struct {
	Controller *WeaveController `ko:"name=controller"`
	Object     Symbol           `ko:"name=object"`
}

func (b WeaveObject) String() string {
	return fmt.Sprintf("WEAVE-OBJECT-%s", b.Object.String())
}

func (b WeaveObject) Select(weaveSpan *Span, path Path) (Shape, Effect, error) {
	if residue, err := b.Controller.Weaver.Select(
		b.Controller.WeaveStepCtx(weaveSpan),
		b.Object,
		path,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returns), b.Controller.WrapEffect(residue.Effect), nil
	}
}

func (b WeaveObject) Link(weaveSpan *Span, name string, monadic bool) (Shape, Effect, error) {
	if residue, err := b.Controller.Weaver.Link(
		b.Controller.WeaveStepCtx(weaveSpan),
		b.Object,
		name,
		monadic,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returns), b.Controller.WrapEffect(residue.Effect), nil
	}
}

func (b WeaveObject) Invoke(weaveSpan *Span) (Shape, Effect, error) {
	if residue, err := b.Controller.Weaver.Invoke(
		b.Controller.WeaveStepCtx(weaveSpan),
		b.Object,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returns), b.Controller.WrapEffect(residue.Effect), nil
	}
}
