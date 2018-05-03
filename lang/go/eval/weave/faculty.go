package weave

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type WeaveOperatorMacro struct {
	WeavePlaceholderMacro `ko:"name=placeholder"`
	Ideal                 Ideal `ko:"name=ideal"`
}

func ParseWeaveOperator(span *Span, ideals Symbol) (Faculty, error) {
	if v, err := IntegrateInterface(span, ideals, typeOfIdeals); err != nil {
		return nil, span.Errorf(err, "weave parsing operator (pkg, name) pairs")
	} else {
		mk := &weaveOperatorFaculty{Ideals: v.(Ideals)}
		return mk.Make(), nil
	}
}

var typeOfIdeals = reflect.TypeOf(Ideals{})

type weaveOperatorFaculty struct {
	Ideals Ideals `ko:"name=ideals,monadic"`
}

func (b *weaveOperatorFaculty) Make() Faculty {
	faculty := Faculty{}
	for _, ideal := range b.Ideals {
		faculty[ideal] = &WeaveOperatorMacro{Ideal: ideal}
	}
	return faculty
}
