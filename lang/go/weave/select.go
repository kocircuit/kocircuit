package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func (local *GoLocal) Select(span *Span, path Path) (Shape, Effect, error) {
	selection, _, err := GoSelect(span, path, local.Image())
	if err != nil {
		return nil, nil, err
	}
	return local.Extend(span, selection), nil, nil
}
