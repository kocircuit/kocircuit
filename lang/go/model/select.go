package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func GoSelectVarietal(span *Span, path Path, arg GoStructure) (Shaper, GoVarietal, error) {
	extractor, extracted, err := GoSelectSimplify(span, path, arg)
	if err != nil {
		return nil, nil, span.Errorf(err, "expectin a %v argument", path)
	}
	switch u := extracted.(type) {
	case GoVarietal:
		return extractor, u, nil
	}
	return nil, nil, span.Errorf(nil, "argument must be a variety, got %s", Sprint(extracted))
}

func GoSelectSimplify(span *Span, path Path, into GoType) (selector Shaper, result GoType, err error) {
	if selector, result, err = GoSelect(span, path, into); err != nil {
		return nil, nil, err
	}
	simple, simplifier := Simplify(span, result)
	return CompressShapers(span, selector, simplifier), simple, nil
}

func GoSelect(span *Span, path Path, into GoType) (selector Shaper, result GoType, err error) {
	selector, result = IdentityShaper(span, into), into
	for _, field := range path {
		s, r, err := SelectFieldShaper(span, field, result)
		if err != nil {
			return nil, nil, err
		}
		selector = CompressShapers(span, selector, s)
		result = r
	}
	return selector, result, nil
}

func SelectFieldShaper(span *Span, field string, into GoType) (selector Shaper, selected GoType, err error) {
	simplified, simplifier := Simplify(span, into)
	switch u := simplified.(type) {
	case *GoEmpty:
		empty := NewGoEmpty(span)
		return &IrreversibleEraseShaper{
			Shaping: Shaping{Origin: span, From: into, To: empty},
		}, empty, nil
	case *GoPtr:
		if selector, selected, err = SelectFieldShaper(span, field, u.Elem); err != nil {
			return nil, nil, err
		}
		selectedThruPtr := NewGoPtr(selected)
		selectOpt := MustVerifyShaper(
			span,
			&OptShaper{
				Shaping: Shaping{Origin: span, From: simplified, To: selectedThruPtr},
				IfNotNil: CompressShapers(
					span,
					&DerefShaper{
						Shaping: Shaping{Origin: span, From: simplified, To: u.Elem}, N: 1,
					},
					selector,
					&RefShaper{
						Shaping: Shaping{Origin: span, From: selected, To: selectedThruPtr}, N: 1,
					},
				),
			},
		)
		return CompressShapers(span, simplifier, selectOpt), selectedThruPtr, nil
	case *GoStruct:
		goFieldName, selected := u.SelectKoField(field)
		var selector Shaper
		if selected == nil {
			selected = NewGoEmpty(span)
			selector = &IrreversibleEraseShaper{
				Shaping: Shaping{Origin: span, From: u, To: selected},
			}
		} else {
			selector = &SelectShaper{
				Shaping: Shaping{Origin: span, From: u, To: selected},
				Field:   goFieldName,
			}
		}
		return CompressShapers(span, simplifier, selector), selected, nil
	case Unknown:
		selected := NewGoUnknown(span)
		unknownSelect := &SelectShaper{
			Shaping: Shaping{Origin: span, From: u, To: selected},
			Field:   "!unknown",
		}
		return CompressShapers(span, simplifier, unknownSelect), selected, nil
	case *GoMap:
		return nil, nil, span.Errorf(nil, "selecting into maps, here %s, is disabled", Sprint(simplified))
	}
	return nil, nil, span.Errorf(nil, "cannot select %q into %s", field, Sprint(simplified))
}
