package model

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

// ReduceZoom returns the type below a chain of slices: [][]...[]unzoomable -> unzoomable
func ReduceZoom(span *Span, s GoType) (unzoomable GoType, zoom Shaper) {
	unzoomable, n := s, 0
	for {
		u, ok := unzoomable.(*GoSlice)
		if !ok {
			break
		}
		n, unzoomable = n+1, u.Elem
	}
	return unzoomable,
		&ZoomShaper{
			Shaping: Shaping{Origin: span, From: s, To: unzoomable}, N: n,
		}
}

// ReduceAlias strips all aliases from the top of the type.
func ReduceAlias(span *Span, s GoType) (t GoType, shaper Shaper) {
	t = s
	for {
		u, ok := t.(*GoAlias)
		if !ok {
			break
		}
		t = u.For
	}
	if s == t {
		return t, IdentityShaper(span, s)
	}
	return t,
		&ConvertTypeShaper{
			Shaping: Shaping{Origin: span, From: s, To: t},
		}
}

func ReduceNeverNilPtr(span *Span, s GoType) (t GoType, shaper Shaper) {
	valvePtr, ok := s.(*GoNeverNilPtr)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	return valvePtr.Elem,
		&DerefShaper{
			Shaping: Shaping{Origin: span, From: s, To: valvePtr.Elem},
			N:       1,
		}
}

// ReducePtrPtr reduces chains **...* to a single *, otherwise the type is left unchanged.
func ReducePtrPtr(span *Span, s GoType) (t GoType, shaper Shaper) {
	t, n := skipPtr(s)
	if n < 2 { //s == t || s == *t
		return s, IdentityShaper(span, s)
	}
	reduced := NewGoPtr(t)
	return reduced,
		&DerefShaper{
			Shaping: Shaping{Origin: span, From: s, To: reduced},
			N:       n - 1,
		}
}

func ReduceAliasAfterPtr(span *Span, typ GoType) (reduced GoType, shaper Shaper) {
	reduced, shaper = typ, IdentityShaper(span, typ)
	for {
		var d Shaper
		reduced, d = StripAliasAfterPtr(span, reduced)
		shaper = CompressShapers(span, shaper, d)
		if IsIdentityShaper(d) {
			break
		}
	}
	return
}

// StripAliasAfterPtr reduces GoPtr{GoAlias{█}} to GoPtr{█}.
func StripAliasAfterPtr(span *Span, typ GoType) (stripped GoType, shaper Shaper) {
	ptr, ok := typ.(*GoPtr)
	if !ok {
		return typ, IdentityShaper(span, typ)
	}
	alias, ok := ptr.Elem.(*GoAlias)
	if !ok {
		return typ, IdentityShaper(span, typ)
	}
	stripped = NewGoPtr(alias.For)
	return stripped, &ConvertTypeShaper{
		Shaping: Shaping{Origin: span, From: typ, To: stripped},
	}
}

func skipPtr(s GoType) (t GoType, n int) {
	t = s
	for {
		u, ok := t.(*GoPtr)
		if !ok {
			break
		}
		t, n = u.Elem, n+1
	}
	return
}

func PeekAlias(t GoType) *GoAlias {
	for {
		switch u := t.(type) {
		case *GoPtr:
			t = u.Elem
		case *GoNeverNilPtr:
			t = u.Elem
		case *GoAlias:
			return u
		case GoReal:
			t = u.Real()
		default:
			return nil
		}
	}
}

// ReducePtrSlice reduces GoPtr{GoSlice{█}} to GoSlice{█}.
// Note that GoNeverNilPtr{GoSlice{█}} simplifies to GoSlice{█} through the greedy boosting in Simplify.
func ReducePtrSlice(span *Span, s GoType) (_ GoType, shaper Shaper) {
	ptr, ok := s.(*GoPtr)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	slice, ok := ptr.Elem.(*GoSlice)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	return slice,
		&OptShaper{
			Shaping: Shaping{Origin: span, From: s, To: slice},
			IfNotNil: &DerefShaper{
				Shaping: Shaping{Origin: span, From: s, To: slice}, N: 1,
			},
		}
}

// ReducePtrEmpty reduces GoPtr{GoEmpty} to GoEmpty.
// Note that GoNeverNilPtr{GoEmpty} simplifies to GoEmpty through the greedy boosting in Simplify.
func ReducePtrEmpty(span *Span, s GoType) (_ GoType, shaper Shaper) {
	ptr, ok := s.(*GoPtr)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	empty, ok := ptr.Elem.(*GoEmpty)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	return empty, &EraseShaper{
		Shaping: Shaping{Origin: span, From: s, To: empty},
	}
}

// ReduceSliceEmpty reduces GoSlice{GoEmpty} to GoEmpty.
func ReduceSliceEmpty(span *Span, s GoType) (_ GoType, shaper Shaper) {
	slice, ok := s.(*GoSlice)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	empty, ok := slice.Elem.(*GoEmpty)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	return empty, &EraseShaper{
		Shaping: Shaping{Origin: span, From: s, To: empty},
	}
}

// ReduceStructEmpty reduces GoStruct{} to GoEmpty.
func ReduceStructEmpty(span *Span, s GoType) (_ GoType, shaper Shaper) {
	strct, ok := s.(*GoStruct)
	if !ok {
		return s, IdentityShaper(span, s)
	}
	if strct.Len() == 0 {
		empty := NewGoEmpty(span)
		return empty, &EraseShaper{
			Shaping: Shaping{Origin: span, From: s, To: empty},
		}
	} else {
		return s, IdentityShaper(span, s)
	}
}

// ReduceArrayEmpty reduces GoArray{GoEmpty} to GoEmpty.
func ReduceArrayEmpty(span *Span, a GoType) (_ GoType, shaper Shaper) {
	array, ok := a.(*GoArray)
	if !ok {
		return a, IdentityShaper(span, a)
	}
	empty, ok := array.Elem.(*GoEmpty)
	if !ok {
		return a, IdentityShaper(span, a)
	}
	return empty, &EraseShaper{
		Shaping: Shaping{Origin: span, From: a, To: empty},
	}
}

// ReduceSingleton reduces GoArray{1, █} to █.
func ReduceSingleton(span *Span, a GoType) (_ GoType, shaper Shaper) {
	array, ok := a.(*GoArray)
	if !ok || array.Len != 1 {
		return a, IdentityShaper(span, a)
	}
	return array.Elem, &UnwrapSingletonShaper{
		Shaping: Shaping{Origin: span, From: a, To: array.Elem},
	}
}
