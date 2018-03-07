package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// ChainShaper ...
type ChainShaper struct {
	Origin *Span    `ko:"name=origin"`
	Chain  []Shaper `ko:"name=chain"`
}

func (cs *ChainShaper) Len() int { return len(cs.Chain) }

func (cs *ChainShaper) String() string { return Sprint(cs) }

func (cs *ChainShaper) Shadow() Shaping {
	return Shaping{
		Origin: cs.Origin,
		From:   cs.Chain[0].Shadow().From,
		To:     cs.Chain[len(cs.Chain)-1].Shadow().To,
	}
}

func (cs *ChainShaper) Shape(t GoType) GoType {
	for _, sh := range cs.Chain {
		t = sh.Shape(t)
	}
	return t
}

func (cs *ChainShaper) Reverse() Shaper {
	r := make([]Shaper, len(cs.Chain))
	for i, shaper := range cs.Chain {
		if q := shaper.(ReversibleShaper).Reverse(); q != nil {
			r[len(cs.Chain)-1-i] = q
		} else {
			return nil // chain not reversible
		}
	}
	return &ChainShaper{Origin: cs.Origin, Chain: r}
}

func (cs *ChainShaper) ShaperVerify(span *Span) error {
	for _, shaper := range cs.Chain {
		if err := shaper.ShaperVerify(span); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ChainShaper) CircuitEffect() *GoCircuitEffect {
	effect := &GoCircuitEffect{}
	for _, sh := range cs.Chain {
		effect = effect.Aggregate(sh.CircuitEffect())
	}
	return effect
}

func (cs *ChainShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	shapingExpr := ofExpr.RenderExpr(fileCtx)
	for _, shaper := range cs.Chain {
		shapingExpr = shaper.RenderExprShaping(fileCtx, &GoVerbatimExpr{shapingExpr})
	}
	return shapingExpr
}

func (cs *ChainShaper) ShaperID() (id string) {
	for _, s := range cs.Chain {
		id = Mix(id, s.ShaperID())
	}
	return
}

func CompressShapers(span *Span, seq ...Shaper) Shaper {
	if len(seq) == 0 {
		panic("o")
	}
	seq = GoCollapseSystem.CollapseChain(expandChains(seq))
	switch len(seq) {
	case 0:
		panic("o")
	case 1:
		return seq[0]
	default:
		return &ChainShaper{Origin: span, Chain: seq}
	}
}

func expandChains(seq []Shaper) (expanded []Shaper) {
	for _, s := range seq {
		switch u := s.(type) {
		case *ChainShaper:
			expanded = append(expanded, u.Chain...)
		default:
			expanded = append(expanded, u)
		}
	}
	return
}
