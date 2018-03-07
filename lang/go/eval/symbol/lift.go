package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func LiftToSeries(span *Span, sym Symbol) *SeriesSymbol {
	switch u := sym.(type) {
	case EmptySymbol:
		if empty, err := MakeSeriesSymbol(span, nil); err != nil {
			panic(err)
		} else {
			return empty
		}
	case *SeriesSymbol:
		return u
	case BasicSymbol, *OpaqueSymbol, *StructSymbol, *VarietySymbol:
		return &SeriesSymbol{
			Type_: &SeriesType{Elem: u.Type()},
			Elem:  Symbols{u},
		}
	default:
		panic("o")
	}
}
