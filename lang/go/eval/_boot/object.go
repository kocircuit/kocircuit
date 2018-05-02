package boot

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func (b *BootController) Wrap(symbol Symbol) Shape {
	return BootObject{Controller: b, Object: symbol}
}

func (b *BootController) UnwrapArg(arg Arg) Symbol {
	return arg.(BootObject).Object
}

func (b *BootController) Unwrap(shape Shape) Symbol {
	return shape.(BootObject).Object
}

// BootObject is a Shape.
type BootObject struct {
	Controller *BootController `ko:"name=controller"`
	Object     Symbol          `ko:"name=object"`
}

func (b BootObject) String() string {
	return fmt.Sprintf("BOOT-OBJECT-%s", b.Object.String())
}

func (b BootObject) Select(bootSpan *Span, path Path) (Shape, Effect, error) {
	if residue, err := b.Controller.Booter.Select(
		b.Controller.BootStepCtx(bootSpan),
		b.Object,
		path,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returned), b.Controller.WrapEffect(residue.Effect), nil
	}
}

func (b BootObject) Link(bootSpan *Span, name string, monadic bool) (Shape, Effect, error) {
	if residue, err := b.Controller.Booter.Link(
		b.Controller.BootStepCtx(bootSpan),
		b.Object,
		name,
		monadic,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returned), b.Controller.WrapEffect(residue.Effect), nil
	}
}

func (b BootObject) Invoke(bootSpan *Span) (Shape, Effect, error) {
	if residue, err := b.Controller.Booter.Invoke(
		b.Controller.BootStepCtx(bootSpan),
		b.Object,
	); err != nil {
		return nil, nil, err
	} else {
		return b.Controller.Wrap(residue.Returned), b.Controller.WrapEffect(residue.Effect), nil
	}
}
