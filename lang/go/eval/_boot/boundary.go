package boot

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type BootEvalBoundary struct {
	Booter *Booter `ko:"name=booter"`
}

func (BootEvalBoundary) Figure(span *Span, figure Figure) (Shape, Effect, error) {
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
	case Macro:
		// macro is either a macro from registry, or from Interpret()
		// return MakeVarietySymbol(u, nil), nil, nil
		vty := XXX
		fig.Func = &u.Value_
	default:
		panic("unknown figure")
	}
	XXX
}

func (BootEvalBoundary) Enter(_ *Span, arg Arg) (Shape, Effect, error) {
	if residue, err := booter.Enter(ctx, arg.(Symbol)); err != nil {
		return nil, nil, err
	} else {
		XXX //save residue.Effect somewhere centrally
		return residue.Returned, nil, nil
	}
}

func (BootEvalBoundary) Leave(_ *Span, shape Shape) (Return, Effect, error) {
	// return shape, nil, nil
}
