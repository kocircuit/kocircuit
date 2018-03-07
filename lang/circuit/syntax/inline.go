package syntax

type Inline struct {
	Design []Design `ko:"name=design"` // inline function definitions
	Series []Term   `ko:"name=series"` // inline step definitions, arising from series composition
}

func (inline Inline) Union(u Inline) Inline {
	return Inline{
		Design: append(append([]Design{}, inline.Design...), u.Design...),
		Series: append(append([]Term{}, inline.Series...), u.Series...),
	}
}
