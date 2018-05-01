package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) Wrap(symbol Symbol) Shape {
	return BootShape{Controller: b, Symbol: symbol}
}

func (b *BootController) UnwrapArg(arg Arg) Symbol {
	return arg.(BootShape).Symbol
}

func (b *BootController) Unwrap(shape Shape) Symbol {
	return shape.(BootShape).Symbol
}

type BootShape struct {
	Controller *BootController `ko:"name=controller"`
	Symbol     Symbol          `ko:"name=symbol"`
}

func (bsym BootShape) String() string {
	return fmt.Sprintf("BOOT-%s", bsym.String())
}

func (bsym BootShape) Select(bootSpan *Span, path Path) (Shape, Effect, error) {
	XXX
}

func (bsym BootShape) Link(bootSpan *Span, name string, monadic bool) (Shape, Effect, error) {
	XXX
}

func (bsym BootShape) Augment(bootSpan *Span, fields Fields) (Shape, Effect, error) {
	XXX
}

func (bsym BootShape) Invoke(bootSpan *Span) (Shape, Effect, error) {
	XXX
}
