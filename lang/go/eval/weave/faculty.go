package weave

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type WeaveReserveMacro struct {
	WeavePlaceholderMacro `ko:"name=placeholder"`
	Ideal                Ideal `ko:"name=ideal"`
}

func ParseWeaveReserve(span *Span, ideals Symbol) (Faculty, error) {
	if v, err := IntegrateInterface(span, ideals, typeOfIdeals); err != nil {
		return nil, span.Errorf(err, "weave parsing reserve (pkg, name) pairs")
	} else {
		mk := &weaveReserveFaculty{Ideals: v.(Ideals)}
		return mk.Make(), nil
	}
}

var typeOfIdeals = reflect.TypeOf(Ideals{})

type weaveReserveFaculty struct {
	Ideals Ideals `ko:"name=ideals,monadic"`
}

func (b *weaveReserveFaculty) Make() Faculty {
	faculty := Faculty{}
	for _, ideal := range b.Ideals {
		faculty[ideal] = &WeaveReserveMacro{Ideal: ideal}
	}
	return faculty
}
