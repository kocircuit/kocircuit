package boot

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

type BootReserveMacro struct {
	BootPlaceholderMacro `ko:"name=placeholder"`
	Ideal                Ideal `ko:"name=ideal"`
}

//XXX: add reserved ideals to booter

func BootReserveFaculty(span *Span, ideals Symbol) (Faculty, error) {
	if v, err := IntegrateInterface(span, ideals, typeOfIdeals); err != nil {
		return nil, span.Errorf(err, "boot parsing reserve (pkg, name) pairs")
	} else {
		return v.(Ideals), nil
	}
}

var typeOfIdeals = reflect.TypeOf(Ideals{})

type bootReserveFaculty struct {
	Ideals Ideals `ko:"name=ideals,monadic"`
}

func (b *bootReserveFaculty) Make() Faculty {
	faculty := Faculty{}
	for _, ideal := range b.Ideals {
		faculty[ideal] = &BootReserveMacro{Ideal: ideal}
	}
	return faculty
}
