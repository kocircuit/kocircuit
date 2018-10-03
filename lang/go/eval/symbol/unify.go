//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package symbol

import (
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func Unify(span *model.Span, x, y Type) (Type, error) {
	ctx := &typingCtx{Span: span}
	return ctx.Unify(x, y)
}

func UnifyTypes(span *model.Span, tt []Type) (Type, error) {
	ctx := &typingCtx{Span: span}
	return ctx.UnifyTypes(tt)
}

func (ctx *typingCtx) UnifyTypes(tt []Type) (unified Type, err error) {
	if len(tt) == 0 {
		return EmptyType{}, nil
	}
	unified = tt[0]
	for i := 1; i < len(tt); i++ {
		if unified, err = ctx.Unify(unified, tt[i]); err != nil {
			return nil, err
		}
	}
	return
}

// Unify(x, y) = Unify(y, x)
// Unify(x, Unify(y, z)) = Unify(Unify(x, y), z)
func (ctx *typingCtx) Unify(x, y Type) (Type, error) {
	switch {
	case IsEmptyType(x) && IsEmptyType(y):
		return EmptyType{}, nil
	case !IsEmptyType(x) && IsEmptyType(y):
		return Optionally(x), nil
	case IsEmptyType(x) && !IsEmptyType(y):
		return Optionally(y), nil
	}
	switch xt := x.(type) {
	case *OptionalType:
		if elem, err := ctx.Unify(xt.Elem, y); err != nil {
			return nil, err
		} else {
			return Optionally(elem), nil
		}
	case *SeriesType:
		switch yt := y.(type) {
		case *OptionalType:
			return ctx.Unify(y, x) // symmetry
		case *SeriesType:
			return ctx.UnifySeries(xt, yt)
		case BasicType, *OpaqueType, *MapType, *StructType, VarietyType, NamedType:
			if elem, err := ctx.Unify(xt.Elem, y); err != nil {
				return nil, err
			} else {
				return &SeriesType{elem}, nil
			}
		}
	case BasicType:
		switch yt := y.(type) {
		case *OptionalType, *SeriesType:
			return ctx.Unify(y, x) // symmetry
		case BasicType:
			return ctx.UnifyBasic(xt, yt)
		}
	case *OpaqueType:
		switch yt := y.(type) {
		case *OptionalType, *SeriesType, BasicType:
			return ctx.Unify(y, x) // symmetry
		case *OpaqueType:
			return ctx.UnifyOpaque(xt, yt)
		}
	case *MapType:
		switch yt := y.(type) {
		case *OptionalType, *SeriesType, BasicType, *OpaqueType:
			return ctx.Unify(y, x) // symmetry
		case *MapType:
			return ctx.UnifyMap(xt, yt)
		}
	case *StructType:
		switch yt := y.(type) {
		case *OptionalType, *SeriesType, BasicType, *OpaqueType, *MapType:
			return ctx.Unify(y, x) // symmetry
		case *StructType:
			return ctx.UnifyStruct(xt, yt)
		}
	case VarietyType:
		switch y.(type) {
		case *OptionalType, *SeriesType, BasicType, *OpaqueType, *MapType, *StructType:
			return ctx.Unify(y, x) // symmetry
		case VarietyType:
			return VarietyType{}, nil
		}
	case NamedType:
		switch yt := y.(type) {
		case *OptionalType, *SeriesType, BasicType, *OpaqueType, *MapType, *StructType, VarietyType:
			return ctx.Unify(y, x) // symmetry
		case NamedType:
			return ctx.UnifyNamed(xt, yt)
		}
	}
	return nil, ctx.Errorf(nil, "%s and %s cannot be unified", tree.Sprint(x), tree.Sprint(y))
}

func (ctx *typingCtx) UnifyBasic(x, y BasicType) (Type, error) {
	if unified, ok := unifyBasic(x, y); ok {
		return unified, nil
	} else {
		return nil, ctx.Errorf(nil, "basic types %s and %s cannot be unified", tree.Sprint(x), tree.Sprint(y))
	}
}

func (ctx *typingCtx) UnifyOpaque(x, y *OpaqueType) (Type, error) {
	if x.Type == y.Type {
		return x, nil
	} else {
		return nil, ctx.Errorf(nil, "opaque types %s and %s cannot be unified", tree.Sprint(x), tree.Sprint(y))
	}
}

func (ctx *typingCtx) UnifySeries(x, y *SeriesType) (*SeriesType, error) {
	if xyElem, err := ctx.Refine("()").Unify(x.Elem, y.Elem); err != nil {
		return nil, ctx.Errorf(nil, "cannot unify sequences %s and %s", tree.Sprint(x), tree.Sprint(y))
	} else {
		return &SeriesType{Elem: xyElem}, nil
	}
}

func (ctx *typingCtx) UnifyNamed(x, y NamedType) (Type, error) {
	if x.Type == y.Type {
		return x, nil
	} else {
		return nil, ctx.Errorf(nil, "named types %s and %s cannot be unified", tree.Sprint(x), tree.Sprint(y))
	}
}

func (ctx *typingCtx) UnifyMap(x, y *MapType) (Type, error) {
	if unified, err := ctx.Unify(x.Value, y.Value); err == nil {
		return &MapType{Value: unified}, nil
	} else {
		return nil, ctx.Errorf(nil, "map types %s and %s cannot be unified", tree.Sprint(x), tree.Sprint(y))
	}
}
