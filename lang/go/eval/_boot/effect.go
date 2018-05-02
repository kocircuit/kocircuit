package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func (b *BootController) WrapEffect(symbol Symbol) Effect {
	return BootEffect{Effect: symbol}
}

func (b *BootController) UnwrapEffect(eff Effect) Symbol {
	return arg.(BootEffect).Effect
}

// BootEffect is an Effect.
type BootEffect struct {
	Effect Symbol `ko:"name=effect"`
}

func (b BootEffect) String() string {
	return fmt.Sprintf("BOOT-EFFECT-%s", b.Effect.String())
}
