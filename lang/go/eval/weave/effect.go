package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func (b *WeaveController) WrapEffect(symbol Symbol) Effect {
	return WeaveEffect{Effect: symbol}
}

func (b *WeaveController) UnwrapEffect(eff Effect) Symbol {
	return eff.(WeaveEffect).Effect
}

// WeaveEffect is an Effect.
type WeaveEffect struct {
	Effect Symbol `ko:"name=effect"`
}

func (b WeaveEffect) String() string {
	return fmt.Sprintf("WEAVE-EFFECT-%s", b.Effect.String())
}
