package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func (b *BootController) Figure(bootSpan *Span, figure Figure) (Shape, Effect, error) {
	fig := &BootFigure{}
	switch u := figure.(type) {
	case Bool:
		fig.Bool = &u.Value_
	case Integer:
		fig.Int64 = &u.Value_
	case Float:
		fig.Float64 = &u.Value_
	case String:
		fig.String = &u.Value_
	case *BootFuncMacro: // from Interpret()
		fig.Functional = &BootFunctional{
			Func: &BootFunc{
				Pkg:  u.Func.Pkg,
				Name: u.Func.Name,
			},
		}
	case *BootMacroMacro: // from faculty
		fig.Functional = &BootFunctional{
			Macro: PtrString(u.Macro),
		}
	default:
		panic("unknown figure")
	}
	if residue, err := b.Booter.Literal(b.BootStepCtx(bootSpan), fig); err != nil {
		return nil, nil, err
	} else {
		return residue.Returned, residue.Effect, nil
	}
}

func (b *BootController) Enter(bootSpan *Span, arg Arg) (Shape, Effect, error) {
	if residue, err := b.Booter.Enter(b.BootStepCtx(bootSpan), arg.(Symbol)); err != nil {
		return nil, nil, err
	} else {
		return residue.Returned, residue.Effect, nil
	}
}

func (b *BootController) Leave(bootSpan *Span, shape Shape) (Return, Effect, error) {
	if residue, err := b.Booter.Leave(b.BootStepCtx(bootSpan), arg.(Symbol)); err != nil {
		return nil, nil, err
	} else {
		return residue.Returned, residue.Effect, nil
	}
}
