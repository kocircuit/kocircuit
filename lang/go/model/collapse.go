package model

type CollapseRule interface {
	Collapse(x, y Shaper) (rewritten []Shaper, shrinks bool)
}

type CollapseSystem []CollapseRule

func (cs CollapseSystem) Collapse(x, y Shaper) (result []Shaper, shrinks bool) {
	result = []Shaper{x, y}
	for _, rule := range cs {
		if result, shrinks = rule.Collapse(x, y); shrinks {
			return result, shrinks
		}
	}
	return
}

func (cs CollapseSystem) CollapseChain(chain []Shaper) []Shaper {
	var pass [2][]Shaper
	pass[0] = make([]Shaper, len(chain))
	copy(pass[0], chain)
	t := 0
	for {
		present, future := t%2, (t+1)%2
		collapsed := false
		pass[future] = nil
		i := 0
		for i < len(pass[present]) {
			if i+1 < len(pass[present]) {
				processed, shrunk := cs.Collapse(pass[present][i], pass[present][i+1])
				if shrunk {
					pass[future] = append(pass[future], processed...)
					collapsed = true
					i += 2
				} else {
					pass[future] = append(pass[future], pass[present][i])
					i += 1
				}
			} else {
				pass[future] = append(pass[future], pass[present][i])
				i += 1
			}
		}
		if !collapsed {
			return pass[future]
		}
		t++
	}
}

var GoCollapseSystem = CollapseSystem{
	// XXX: add one collapser that collapses inverse pairs generically by ShapeID
	CollapseReRe{},
	CollapseLeftIdentity{},
	CollapseRightIdentity{},
	CollapseRefDeref{},
	CollapseDerefRef{},
	CollapseEraseUnerase{},
	CollapseUneraseErase{},
}

type CollapseReRe struct{}

func (CollapseReRe) Collapse(x, y Shaper) ([]Shaper, bool) {
	s0, ok := x.(*ReShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	s1, ok := y.(*ReShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	return []Shaper{
		&ReShaper{
			Shaping: Shaping{Origin: s1.Origin, From: s0.Shaping.From, To: s1.To},
		},
	}, true
}

type CollapseLeftIdentity struct{}

func (CollapseLeftIdentity) Collapse(x, y Shaper) ([]Shaper, bool) {
	if !IsIdentityShaper(x) {
		return []Shaper{x, y}, false
	}
	return []Shaper{y}, true
}

type CollapseRightIdentity struct{}

func (CollapseRightIdentity) Collapse(x, y Shaper) ([]Shaper, bool) {
	if !IsIdentityShaper(y) {
		return []Shaper{x, y}, false
	}
	return []Shaper{x}, true
}

type CollapseRefDeref struct{}

func (CollapseRefDeref) Collapse(x, y Shaper) ([]Shaper, bool) {
	s0, ok := x.(*RefShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	s1, ok := y.(*DerefShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	if s0.N != s1.N {
		return []Shaper{x, y}, false
	}
	return collapseConvertShapers(x, y)
}

type CollapseDerefRef struct{}

func (CollapseDerefRef) Collapse(x, y Shaper) ([]Shaper, bool) {
	s0, ok := x.(*DerefShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	s1, ok := y.(*RefShaper)
	if !ok {
		return []Shaper{x, y}, false
	}
	if s0.N != s1.N {
		return []Shaper{x, y}, false
	}
	return collapseConvertShapers(x, y)
}

type CollapseEraseUnerase struct{}

func (CollapseEraseUnerase) Collapse(x, y Shaper) ([]Shaper, bool) {
	if _, ok := x.(*EraseShaper); !ok {
		return []Shaper{x, y}, false
	}
	if _, ok := y.(*UneraseShaper); !ok {
		return []Shaper{x, y}, false
	}
	return collapseConvertShapers(x, y)
}

type CollapseUneraseErase struct{}

func (CollapseUneraseErase) Collapse(x, y Shaper) ([]Shaper, bool) {
	if _, ok := x.(*UneraseShaper); !ok {
		return []Shaper{x, y}, false
	}
	if _, ok := y.(*EraseShaper); !ok {
		return []Shaper{x, y}, false
	}
	return collapseConvertShapers(x, y)
}

func collapseConvertShapers(x, y Shaper) ([]Shaper, bool) {
	start, stop := x.Shadow().From, y.Shadow().To
	if start.TypeID() == stop.TypeID() {
		return []Shaper{IdentityShaper(y.Shadow().Origin, start)}, true
	}
	return []Shaper{
		&ConvertTypeShaper{
			Shaping: Shaping{Origin: y.Shadow().Origin, From: start, To: stop},
		},
	}, true
}
