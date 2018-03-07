package tree

type MultiSet []MultiSetElem

type MultiSetElem struct {
	Elem  interface{}
	Count int
}

func (mset MultiSet) Size() int {
	return len(mset)
}

func (mset MultiSet) Count(elem interface{}) int {
	for _, mset := range mset {
		if mset.Elem == elem {
			return mset.Count
		}
	}
	return 0
}

func (mset MultiSet) Copy() MultiSet {
	return append(MultiSet{}, mset...)
}

func (mset MultiSet) Add(elem interface{}) MultiSet {
	r := mset.Copy()
	for i, e := range r {
		if e.Elem == elem {
			r[i].Count++
			return r
		}
	}
	return append(r, MultiSetElem{Elem: elem, Count: 1})
}
