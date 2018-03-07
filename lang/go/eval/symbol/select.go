package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func (ss *StructSymbol) SelectMonadic() Symbol {
	for _, field := range ss.Field {
		if field.Monadic {
			return field.Value
		}
	}
	return EmptySymbol{}
}

func (ss *StructSymbol) Select(span *Span, path Path) (_ Shape, _ Effect, err error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return ss.Walk(path[0]).Select(span, path[1:])
	}
}

func (ss *StructSymbol) Walk(step string) Symbol {
	if found := FindFieldSymbol(step, ss.Field); found != nil {
		return found.Value
	} else {
		return EmptySymbol{}
	}
}

func FindFieldSymbol(name string, fields FieldSymbols) *FieldSymbol {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}
