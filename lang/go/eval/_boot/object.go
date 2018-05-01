package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) Wrap(symbol Symbol) Shape {
	return BootObject{Controller: b, Symbol: symbol}
}

func (b *BootController) UnwrapArg(arg Arg) Symbol {
	return arg.(BootObject).Symbol
}

func (b *BootController) Unwrap(shape Shape) Symbol {
	return shape.(BootObject).Symbol
}

// BootObject is a Shape.
type BootObject struct {
	Controller *BootController `ko:"name=controller"`
	Object     Symbol          `ko:"name=object"`
}

func (b BootObject) String() string {
	return fmt.Sprintf("BOOT-%s", b.String())
}

func (b BootObject) Select(bootSpan *Span, path Path) (Shape, Effect, error) {
	if residue, err := b.Booter.Select(b.Controller.BootStepCtx(bootSpan), b.Object, path); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.Wrap(residue.Effect), nil
	}
}

func (b BootObject) Link(bootSpan *Span, name string, monadic bool) (Shape, Effect, error) {
	if residue, err := b.Booter.Link(b.Controller.BootStepCtx(bootSpan), b.Object, name, monadic); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.Wrap(residue.Effect), nil
	}
}

func (b BootObject) Invoke(bootSpan *Span) (Shape, Effect, error) {
	if residue, err := b.Booter.Invoke(b.Controller.BootStepCtx(bootSpan), b.Object); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returned), b.Wrap(residue.Effect), nil
	}
}

func (b BootObject) Augment(bootSpan *Span, fields Fields) (Shape, Effect, error) {
	XXX
}
