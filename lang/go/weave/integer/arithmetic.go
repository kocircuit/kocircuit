package integer

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

type IntegerSolnArg struct {
	Arg     GoStructure   `ko:"name=arg"`
	Integer TypeExtractor `ko:"name=integer"`
}

func (arg *IntegerSolnArg) CircuitEffect() *GoCircuitEffect {
	return arg.Integer.Extractor.CircuitEffect()
}

func SolveIntegerArg(span *Span, arg GoStructure) (soln *IntegerSolnArg, err error) {
	if monadic := StructureMonadicField(arg); monadic == nil {
		return nil, span.Errorf(nil, "integer arithmetic expects a monadic argument, got %s", Sprint(arg))
	} else {
		soln = &IntegerSolnArg{Arg: arg}
		if soln.Integer.Extractor, soln.Integer.Type, err = GoSelectSimplify(span, Path{monadic.KoName()}, arg); err != nil {
			return nil, span.Errorf(err, "integer arithmetic expects an etc argument")
		}
		return soln, nil
	}
}
