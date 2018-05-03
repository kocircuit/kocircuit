package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (b *WeaveController) Figure(weaveSpan *Span, figure Figure) (Shape, Effect, error) {
	fig := &WeaveFigure{}
	switch u := figure.(type) {
	case Bool:
		fig.Bool = &u.Value_
	case Integer:
		fig.Int64 = &u.Value_
	case Float:
		fig.Float64 = &u.Value_
	case String:
		fig.String = &u.Value_
	case *WeaveFuncMacro: // from Interpret()
		fig.Functional = &WeaveFunctional{
			Func: &WeaveFunc{Pkg: u.Func.Pkg, Name: u.Func.Name},
		}
	case *WeaveOperatorMacro: // from faculty of operators
		fig.Functional = &WeaveFunctional{
			Operator: &WeaveOperator{Pkg: u.Ideal.Pkg, Name: u.Ideal.Name},
		}
	default:
		panic("unknown figure")
	}
	if residue, err := b.Weaver.Literal(b.WeaveStepCtx(weaveSpan), fig); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returns), b.WrapEffect(residue.Effect), nil
	}
}

func (b *WeaveController) Enter(weaveSpan *Span, arg Arg) (Shape, Effect, error) {
	if residue, err := b.Weaver.Enter(b.WeaveStepCtx(weaveSpan), b.UnwrapArg(arg)); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returns), b.WrapEffect(residue.Effect), nil
	}
}

func (b *WeaveController) Leave(weaveSpan *Span, shape Shape) (Return, Effect, error) {
	if residue, err := b.Weaver.Leave(b.WeaveStepCtx(weaveSpan), b.Unwrap(shape)); err != nil {
		return nil, nil, err
	} else {
		return b.Wrap(residue.Returns), b.WrapEffect(residue.Effect), nil
	}
}
