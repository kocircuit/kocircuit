package model

import (
	"fmt"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func TestShaper(t *testing.T) {
	ctx := &RenderCtx{}
	for i, test := range testShaper {
		fileCtx := ctx.PkgContext(GoGroupPath{Group: KoPkgGroup, Path: "test"}).FileContext()
		fmt.Sprintln("test", i)
		fmt.Sprintln(test.Shaper.RenderExprShaping(fileCtx, &GoVerbatimExpr{"X"}))
	}
}

var testSpan = NewSpan()

var testShaper = []struct {
	Shaper Shaper
}{
	{Shaper: &ConvertTypeShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   GoInt64,
			To:     GoInt32,
		}},
	},
	{Shaper: &DerefShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   NewGoPtr(GoInt64),
			To:     GoInt64,
		},
		N: 1},
	},
	{Shaper: &RefShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   GoString,
			To:     NewGoPtr(GoString),
		},
		N: 1},
	},
	{Shaper: &AssertTypeShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From: NewGoAlias(
				&GoAddress{
					Name:      "I",
					GroupPath: GoGroupPath{Group: KoPkgGroup, Path: "test/util"},
				},
				NewGoInterface(UnrestrictedInterface),
			),
			To: NewGoPtr(GoString),
		}},
	},
	{Shaper: &ZoomShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   NewGoSlice(GoInt64), // []int64
			To:     GoInt64,             // int64
		}, N: 1},
	},
	{Shaper: &UnzoomShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   GoInt64,             // int64
			To:     NewGoSlice(GoInt64), // []int64
		}, N: 1},
	},
	{Shaper: &BatchShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From:   NewGoArray(5, NewGoPtr(GoInt8)),
			To:     NewGoSlice(GoInt8),
		},
		Elem: &DerefShaper{Shaping: Shaping{
			Origin: testSpan,
			From:   NewGoPtr(GoInt8),
			To:     GoInt8,
		}, N: 1},
	}},
	{Shaper: &StructMapShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From: NewGoStruct(
				&GoField{Name: "G1", Tag: KoTags("g1", false), Type: GoInt8},
				&GoField{Name: "G2", Tag: KoTags("g2", false), Type: GoBool},
			),
			To: NewGoMap(GoString, NewGoInterface(UnrestrictedInterface)),
		},
	}},
	{Shaper: &StructStructShaper{
		Shaping: Shaping{
			Origin: testSpan,
			From: NewGoStruct(
				&GoField{Name: "A1", Tag: KoTags("a1", false), Type: GoInt8},
				&GoField{Name: "A2", Tag: KoTags("a2", false), Type: GoBool},
			),
			To: NewGoStruct(
				&GoField{Name: "A1", Tag: KoTags("a1", false), Type: NewGoPtr(GoInt8)},
				&GoField{Name: "A2", Tag: KoTags("a2", false), Type: NewGoPtr(GoBool)},
			),
		},
		Field: []*FieldShaper{
			{From: "A1", To: "AA1", Shaper: &RefShaper{
				Shaping: Shaping{
					Origin: testSpan,
					From:   GoInt8,
					To:     NewGoPtr(GoInt8),
				},
				N: 1,
			}},
			{From: "A2", To: "AA2", Shaper: &RefShaper{
				Shaping: Shaping{
					Origin: testSpan,
					From:   GoBool,
					To:     NewGoPtr(GoBool),
				},
				N: 1,
			}},
		},
	}},
}
