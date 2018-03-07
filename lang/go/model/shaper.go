package model

import (
	"fmt"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/hash"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

type Shaper interface {
	String() string
	Shape(GoType) GoType
	ShaperID() string                // uniqely identifies shaper, for dedup on render
	CircuitEffect() *GoCircuitEffect // CircuitEffect() []GoFuncExpr
	RenderExprShaping(GoFileContext, GoExpr) string
	Shadow() Shaping
	ShaperVerify(*Span) error
}

func MustVerifyShaper(span *Span, sh Shaper) Shaper {
	if err := sh.ShaperVerify(span); err != nil {
		panic(err)
	}
	return sh
}

type ReversibleShaper interface {
	Shaper
	Reverse() Shaper
}

// Shaping specifies to and from top-level types.
type Shaping struct {
	Origin *Span  `ko:"name=origin"`
	From   GoType `ko:"name=from"` // GoX or GoAlias{GoX}
	To     GoType `ko:"name=to"`   // GoX or GoAlias{GoX}
}

func (e Shaping) ShaperVerify(span *Span) error { return nil }

func (e Shaping) Shadow() Shaping { return e }

func (e Shaping) Shape(x GoType) GoType {
	if x.TypeID() != e.From.TypeID() {
		panic("o")
	}
	return e.To
}

func (e Shaping) CircuitEffect() *GoCircuitEffect { return nil }

func (e Shaping) ShapingID() string {
	return Mix(e.From.TypeID(), e.To.TypeID())
}

func (e Shaping) Flip() Shaping {
	return Shaping{Origin: e.Origin, From: e.To, To: e.From}
}

func (e Shaping) IsIdentity() bool {
	return e.From.TypeID() == e.To.TypeID()
}

// SelectShaper...
type SelectShaper struct { // go static conversions
	Shaping `ko:"name=shaping"`
	Field   string `ko:"name=field"`
}

func (sh *SelectShaper) ShaperID() string {
	return Mix("SelectShaper", sh.ShapingID(), "field", sh.Field)
}

func (sh *SelectShaper) String() string { return Sprint(sh) }

func (sh *SelectShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	expr := &GoSelectExpr{Into: ofExpr, Field: sh.Field}
	return expr.RenderExpr(fileCtx)
}

func (SelectShaper) Reverse() Shaper { return nil }

// NilShaper...
type NilShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *NilShaper) ShaperID() string { return Mix("NilShaper", sh.ShapingID()) }

func (sh *NilShaper) String() string { return Sprint(sh) }

func (sh *NilShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	return fmt.Sprintf("(%s)(nil)", sh.To.RenderRef(fileCtx))
}

// ReShaper...
type ReShaper struct {
	Shaping `ko:"name=shaping"`
}

// IdentityShaper returns a canonical shaper that does not have a syntactic or semantic effect,
// but it serves as a placeholder within shaper chains, during compressing and collapsing for
// preserving type shaping information.
func IdentityShaper(span *Span, fromTo GoType) Shaper {
	return &ReShaper{
		Shaping: Shaping{Origin: span, From: fromTo, To: fromTo},
	}
}

func IsIdentityShaper(shaper Shaper) bool {
	return shaper.Shadow().IsIdentity() // matches IdentityShaper(...)
}

func (sh *ReShaper) ShaperID() string { return Mix("ReShaper", sh.ShapingID()) }

func (sh *ReShaper) String() string { return Sprint(sh) }

func (sh *ReShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return ofExpr.RenderExpr(fileCtx)
}

func (sh *ReShaper) Reverse() Shaper {
	return &ReShaper{Shaping: sh.Flip()}
}

// NumberShaper...
type NumberShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *NumberShaper) ShaperID() string { return Mix("NumberShaper", sh.ShapingID()) }

func (sh *NumberShaper) String() string { return Sprint(sh) }

func (sh *NumberShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	return fmt.Sprintf(
		"%s(%s)",
		sh.To.RenderRef(fileCtx),
		sh.From.(GoNumber).NumberExpr().RenderExpr(fileCtx),
	)
}

func (sh *NumberShaper) Reverse() Shaper { panic("irreversible") }

// UneraseNumberShaper...
type UneraseNumberShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *UneraseNumberShaper) ShaperID() string { return Mix("UneraseNumberShaper", sh.ShapingID()) }

func (sh *UneraseNumberShaper) String() string { return Sprint(sh) }

func (sh *UneraseNumberShaper) RenderExprShaping(fileCtx GoFileContext, _ GoExpr) string {
	return fmt.Sprintf(
		"%s(%s)",
		sh.To.RenderRef(fileCtx),
		sh.To.(GoNumber).NumberExpr().RenderExpr(fileCtx),
	)
}

func (sh *UneraseNumberShaper) Reverse() Shaper { panic("irreversible") }

// ConvertTypeShaper captures an explicit Go type conversion.
type ConvertTypeShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *ConvertTypeShaper) ShaperID() string { return Mix("ConvertTypeShaper", sh.ShapingID()) }

func (sh *ConvertTypeShaper) String() string { return Sprint(sh) }

func (sh *ConvertTypeShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return fmt.Sprintf("(%s)(%s)", sh.To.RenderRef(fileCtx), ofExpr.RenderExpr(fileCtx))
}

func (sh *ConvertTypeShaper) Reverse() Shaper {
	return &ConvertTypeShaper{Shaping: sh.Flip()}
}

// OptShaper captures: *type -> Value(type)
// The to-type of the shaping (and the value shaper) must have nil as zero value.
type OptShaper struct {
	Shaping  `ko:"name=shaping"`
	IfNotNil Shaper `ko:"name=if_not_nil" ctx:"expand"`
}

func (sh *OptShaper) ShaperVerify(span *Span) error {
	MustBeAssignIso(span, sh.Shadow().From, sh.IfNotNil.Shadow().From)
	MustBeAssignIso(span, sh.Shadow().To, sh.IfNotNil.Shadow().To)
	return nil
}

func (sh *OptShaper) ShaperID() string {
	return Mix("OptShaper", sh.ShapingID(), "IfNotNil", sh.IfNotNil.ShaperID())
}

func (sh *OptShaper) String() string { return Sprint(sh) }

func (sh *OptShaper) CircuitEffect() *GoCircuitEffect {
	return sh.IfNotNil.CircuitEffect().AggregateDuctFunc(sh)
}

func (sh *OptShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *OptShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *OptShaper) funcExpr() GoFuncExpr {
	fromExpr := &GoVerbatimExpr{"from"}
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line: []GoExpr{
			&GoIfThenExpr{
				If: &GoEqualityExpr{
					Left:  fromExpr,
					Right: NilExpr,
				},
				Then: []GoExpr{
					&GoReturnExpr{NilExpr},
				},
			},
			&GoReturnExpr{
				&GoShapeExpr{
					Shaper: sh.IfNotNil,
					Expr:   fromExpr,
				},
			},
		},
	}
}

func (sh *OptShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("opt_shaper_%s", sh.ShapingID()),
	}
}

func (sh *OptShaper) Reverse() Shaper {
	return &OptShaper{
		Shaping:  sh.Shaping.Flip(),
		IfNotNil: sh.IfNotNil.(ReversibleShaper).Reverse(),
	}
}

// AssertTypeShaper captures a Go type assertion.
type AssertTypeShaper struct {
	Shaping `ko:"name=shaping"`
}

func (sh *AssertTypeShaper) ShaperID() string { return Mix("AssertTypeShaper", sh.ShapingID()) }

func (sh *AssertTypeShaper) String() string { return Sprint(sh) }

func (sh *AssertTypeShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return fmt.Sprintf("%s.(%s)", ofExpr.RenderExpr(fileCtx), sh.To.RenderRef(fileCtx))
}

func (sh *AssertTypeShaper) Reverse() Shaper {
	return &ReShaper{Shaping: sh.Flip()}
}

// DerefShaper...
type DerefShaper struct {
	Shaping `ko:"name=shaping"`
	N       int `ko:"name=n" ctx:"expand"`
}

func (sh *DerefShaper) ShaperID() string {
	return Mix("DerefShaper", sh.ShapingID(), "N", MixInterface(sh.N))
}

func (sh *DerefShaper) String() string { return Sprint(sh) }

func (sh *DerefShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return fmt.Sprintf("(%s(%s))", strings.Repeat("*", sh.N), ofExpr.RenderExpr(fileCtx))
}

func (sh *DerefShaper) Reverse() Shaper {
	return &RefShaper{Shaping: sh.Flip(), N: sh.N}
}

// RefShaper...
type RefShaper struct {
	Shaping `ko:"name=shaping"`
	N       int `ko:"name=n" ctx:"expand"` // number of times to reference
}

func (sh *RefShaper) ShaperID() string {
	return Mix("RefShaper", sh.ShapingID(), "N", MixInterface(sh.N))
}

func (sh *RefShaper) String() string { return Sprint(sh) }

func (sh *RefShaper) CircuitEffect() *GoCircuitEffect {
	return &GoCircuitEffect{
		DuctFunc: []GoFuncExpr{sh},
	}
}

func (sh *RefShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *RefShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *RefShaper) funcExpr() GoFuncExpr {
	line := []GoExpr{}
	for i := 0; i < sh.N; i++ { // i =previndex
		line = append(line,
			&GoColonAssignExpr{
				Left:  &GoVerbatimExpr{fmt.Sprintf("from%d", i+1)},
				Right: &GoVerbatimExpr{fmt.Sprintf("&from%d", i)},
			},
		)
	}
	line = append(line, &GoReturnExpr{
		&GoVerbatimExpr{fmt.Sprintf("from%d", sh.N)},
	})
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    &GoVerbatimExpr{"from0"},
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line:       line,
	}
}

func (sh *RefShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("ref_shaper_%s", sh.ShapingID()),
	}
}

func (sh *RefShaper) Reverse() Shaper {
	return &DerefShaper{Shaping: sh.Flip(), N: sh.N}
}

// ZoomShaper captures: [][]...[]█ -> █
type ZoomShaper struct {
	Shaping `ko:"name=shaping"`
	N       int `ko:"name=n"` // number of slice wrappers
}

func (sh *ZoomShaper) ShaperID() string {
	return Mix("ZoomShaper", sh.ShapingID(), "N", MixInterface(sh.N))
}

func (sh *ZoomShaper) String() string { return Sprint(sh) }

func (sh *ZoomShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return strings.Join(
		[]string{
			ofExpr.RenderExpr(fileCtx),
			strings.Repeat("[0]", sh.N),
		},
		"",
	)
}

func (sh *ZoomShaper) Reverse() Shaper {
	return &UnzoomShaper{
		Shaping: sh.Flip(), N: sh.N,
	}
}

// UnzoomShaper captures: █ -> [][]...[]█
type UnzoomShaper struct {
	Shaping `ko:"name=shaping"`
	N       int `ko:"name=n"` // number of slice wrappers
}

func (sh *UnzoomShaper) ShaperID() string {
	return Mix("UnzoomShaper", sh.ShapingID(), "N", MixInterface(sh.N))
}

func (sh *UnzoomShaper) String() string { return Sprint(sh) }

func (sh *UnzoomShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	return strings.Join(
		[]string{
			strings.Repeat("[]", sh.N),
			unslice(sh.To, sh.N).RenderRef(fileCtx),
			strings.Repeat("{", sh.N),
			ofExpr.RenderExpr(fileCtx),
			strings.Repeat("}", sh.N),
		},
		"",
	)
}

func unslice(t GoType, n int) GoType {
	for i := 0; i < n; i++ {
		t = t.(*GoSlice).Elem
	}
	return t
}

func (sh *UnzoomShaper) Reverse() Shaper {
	return &ZoomShaper{
		Shaping: sh.Flip(), N: sh.N,
	}
}

// BatchShaper shapes arrays or slices into arrays or slices: Array-Array, Array-Slice, Slice-Slice
//	[N]P -> [N]Q
//	[N]P -> []Q
//	[]P -> []Q
type BatchShaper struct {
	Shaping `ko:"name=shaping"`
	Elem    Shaper `ko:"name=elem" ctx:"expand"`
}

func (sh *BatchShaper) ShaperID() string {
	return Mix("BatchShaper", sh.ShapingID(), "Elem", sh.Elem.ShaperID())
}

func (sh *BatchShaper) String() string { return Sprint(sh) }

func (sh *BatchShaper) CircuitEffect() *GoCircuitEffect {
	return sh.Elem.CircuitEffect().AggregateDuctFunc(sh)
}

func (sh *BatchShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *BatchShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("batch_shaper_%s", sh.ShapingID()),
	}
}

func (sh *BatchShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *BatchShaper) funcExpr() GoFuncExpr {
	toExpr, fromExpr := &GoVerbatimExpr{"to"}, &GoVerbatimExpr{"from"}
	indexExpr, elemExpr := &GoVerbatimExpr{"i"}, &GoVerbatimExpr{"elem"}
	var initExpr GoExpr
	if IsArray(sh.To) { // array
		initExpr = &GoColonAssignExpr{
			Left:  toExpr,
			Right: &GoZeroExpr{sh.To},
		}
	} else { // slice
		initExpr = &GoColonAssignExpr{
			Left: toExpr,
			Right: &GoCallExpr{
				Func: MakeExpr,
				Arg: []GoExpr{
					&GoTypeRefExpr{sh.To},
					&GoCallExpr{
						Func: LenExpr,
						Arg:  []GoExpr{fromExpr},
					},
				},
			},
		}
	}
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line: []GoExpr{
			initExpr,
			&GoForExpr{
				Range: &GoColonAssignExpr{
					Left:  &GoListExpr{Elem: []GoExpr{indexExpr, elemExpr}},
					Right: &GoRangeExpr{fromExpr},
				},
				Line: []GoExpr{
					&GoAssignExpr{
						Left:  UnderlineExpr,
						Right: elemExpr,
					},
					&GoAssignExpr{
						Left:  &GoIndexExpr{Container: toExpr, Index: indexExpr},
						Right: &GoShapeExpr{Shaper: sh.Elem, Expr: elemExpr},
					},
				},
			},
			&GoReturnExpr{toExpr},
		},
	}
}

// StructMapShaper shapes struct into map[string]CommonFieldType.
// The map key holds the Ko name of the corresponding field.
type StructMapShaper struct {
	Shaping `ko:"name=shaping"`
	Field   []*FieldShaper `ko:"name=field" ctx:"expand"`
}

func (sh *StructMapShaper) ShaperID() string {
	return Mix("StructMapShaper", sh.ShapingID(), "Field", MixFieldShapers(sh.Field))
}

func (sh *StructMapShaper) String() string { return Sprint(sh) }

func (sh *StructMapShaper) CircuitEffect() *GoCircuitEffect {
	effect := &GoCircuitEffect{DuctFunc: []GoFuncExpr{sh}}
	for _, f := range sh.Field {
		effect = effect.Aggregate(f.Shaper.CircuitEffect())
	}
	return effect
}

func (sh *StructMapShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *StructMapShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("struct_map_shaper_%s", sh.ShapingID()),
	}
}

func (sh *StructMapShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *StructMapShaper) funcExpr() GoFuncExpr {
	fromExpr, toExpr := &GoVerbatimExpr{"from"}, &GoVerbatimExpr{"to"}
	line := []GoExpr{
		&GoColonAssignExpr{Left: toExpr, Right: &GoZeroExpr{sh.To}},
	}
	line = append(line, renderKeyAssignments(fromExpr, toExpr, sh.Field)...)
	line = append(line, &GoReturnExpr{toExpr})
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line:       line,
	}
}

func renderKeyAssignments(fromExpr, toExpr GoExpr, field []*FieldShaper) []GoExpr {
	r := make([]GoExpr, len(field))
	for i, fsh := range field {
		r[i] = &GoAssignExpr{
			Left: &GoIndexExpr{
				Container: toExpr,
				Index:     &GoQuoteExpr{fsh.To},
			},
			Right: &GoShapeExpr{
				Shaper: fsh.Shaper,
				Expr:   &GoSelectExpr{Into: fromExpr, Field: fsh.From},
			},
		}
	}
	return r
}

// StructStructShaper shapes one struct into another.
type StructStructShaper struct {
	Shaping `ko:"name=shaping"`
	Field   []*FieldShaper `ko:"name=field"`
}

type FieldShaper struct {
	From   string `ko:"name=from"` // go name of from-field
	To     string `ko:"name=to"`   // go name of to-field
	Shaper Shaper `ko:"name=shaper"`
}

func MixFieldShapers(fs []*FieldShaper) string {
	h := []string{}
	for _, fs := range fs {
		h = append(h, Mix(fs.From, fs.To, fs.Shaper.ShaperID()))
	}
	return Mix(h...)
}

func (sh *StructStructShaper) ShaperID() string {
	return Mix("StructStructShaper", sh.ShapingID(), "Field", MixFieldShapers(sh.Field))
}

func (sh *StructStructShaper) String() string { return Sprint(sh) }

func (sh *StructStructShaper) CircuitEffect() *GoCircuitEffect {
	effect := &GoCircuitEffect{DuctFunc: []GoFuncExpr{sh}}
	for _, f := range sh.Field {
		effect = effect.Aggregate(f.Shaper.CircuitEffect())
	}
	return effect
}

func (sh *StructStructShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *StructStructShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("struct_struct_shaper_%s", sh.ShapingID()),
	}
}

func (sh *StructStructShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *StructStructShaper) funcExpr() GoFuncExpr {
	fromExpr := &GoVerbatimExpr{"from"}
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line: []GoExpr{
			&GoReturnExpr{
				&GoMakeStructExpr{
					For:   sh.To,
					Field: renderFieldAssignments(fromExpr, sh.Field),
				},
			},
		},
	}
}

func renderFieldAssignments(fromExpr GoExpr, field []*FieldShaper) []*GoFieldExpr {
	r := make([]*GoFieldExpr, len(field))
	for i, fsh := range field {
		r[i] = &GoFieldExpr{
			Field: &GoField{
				Name: fsh.To,
				Type: fsh.Shaper.Shadow().To,
			},
			Expr: &GoShapeExpr{
				Shaper: fsh.Shaper,
				Expr:   &GoSelectExpr{Into: fromExpr, Field: fsh.From},
			},
		}
	}
	return r
}

// MapMapShaper shapes map[FromKey]FromValue to map[ToKey]FromValue.
type MapMapShaper struct {
	Shaping `ko:"name=shaping"`
	Key     Shaper `ko:"name=key" ctx:"expand"`
	Value   Shaper `ko:"name=value" ctx:"expand"`
}

func (sh *MapMapShaper) ShaperID() string {
	return Mix("MapMapShaper", sh.ShapingID(), "Key", sh.Key.ShaperID(), "Value", sh.Value.ShaperID())
}

func (sh *MapMapShaper) String() string { return Sprint(sh) }

func (sh *MapMapShaper) CircuitEffect() *GoCircuitEffect {
	effect := &GoCircuitEffect{
		DuctFunc: []GoFuncExpr{sh},
	}
	return effect.Aggregate(sh.Key.CircuitEffect()).Aggregate(sh.Value.CircuitEffect())
}

func (sh *MapMapShaper) RenderExprShaping(fileCtx GoFileContext, ofExpr GoExpr) string {
	call := &GoCallFuncExpr{Func: sh, Arg: []GoExpr{ofExpr}}
	return call.RenderExpr(fileCtx)
}

func (sh *MapMapShaper) NameExpr() GoExpr {
	return &GoNameExpr{
		Origin: sh.Origin,
		Name:   fmt.Sprintf("map_map_shaper_%s", sh.ShapingID()),
	}
}

func (sh *MapMapShaper) RenderExpr(fileCtx GoFileContext) string {
	return sh.funcExpr().RenderExpr(fileCtx)
}

func (sh *MapMapShaper) funcExpr() GoFuncExpr {
	fromExpr, toExpr := &GoVerbatimExpr{"from"}, &GoVerbatimExpr{"to"}
	keyExpr, valueExpr := &GoVerbatimExpr{"key"}, &GoVerbatimExpr{"value"}
	return &GoShaperFuncExpr{
		// Comment:    Sprint(sh),
		FuncName:   sh.NameExpr(),
		ArgName:    fromExpr,
		ArgType:    sh.From,
		ReturnType: sh.To,
		Line: []GoExpr{
			&GoColonAssignExpr{
				Left:  toExpr,
				Right: &GoZeroExpr{sh.To},
			},
			&GoForExpr{
				Range: &GoColonAssignExpr{
					Left:  &GoListExpr{Elem: []GoExpr{keyExpr, valueExpr}},
					Right: &GoRangeExpr{fromExpr},
				},
				Line: []GoExpr{
					&GoAssignExpr{
						Left: &GoIndexExpr{
							Container: toExpr,
							Index:     keyExpr,
						},
						Right: valueExpr,
					},
				},
			},
			&GoReturnExpr{toExpr},
		},
	}
}
