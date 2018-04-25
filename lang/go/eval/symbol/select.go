package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (ss *StructSymbol) GetMonadic() (Symbol, bool) {
	for _, field := range ss.Field {
		if field.Monadic {
			return field.Value, true
		}
	}
	return nil, false
}

func (ss *StructSymbol) SelectMonadic() Symbol {
	for _, field := range ss.Field {
		if field.Monadic {
			return field.Value
		}
	}
	return EmptySymbol{}
}

func (ss *StructSymbol) SelectArg(span *Span, name string, monadic bool) (Shape, Effect, error) {
	if found := ss.FindName(name); found != nil {
		return found.Value, nil, nil
	} else if monadic {
		if found := ss.FindMonadic(); found != nil {
			return found.Value, nil, nil
		}
	}
	return EmptySymbol{}, nil, nil
}

func (ss *StructSymbol) Select(span *Span, path Path) (_ Shape, _ Effect, err error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return ss.Walk(path[0]).Select(span, path[1:])
	}
}

func (ss *StructSymbol) Walk(step string) Symbol {
	if found := ss.FindName(step); found != nil {
		return found.Value
	} else {
		return EmptySymbol{}
	}
}
