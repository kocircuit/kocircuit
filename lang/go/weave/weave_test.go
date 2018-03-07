package weave

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/model"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func TestWeave(t *testing.T) {
	initTest()
	weaveTests.Run(t)
}

var weaveTests = WeaveTests{
	{
		Enabled: true,
		Name:    "select-join",
		File: `
		G(u, w) {
			p: (a: u, b: w)
			q: (a: w, b: u)
			return: (
				s: (p.a, q.b)
				t: (p.b, q.a)
			)
		}
		Main(x, y) { return: G(u: x, w: y) }
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
			&GoField{Name: "Y", Type: GoFloat64, Tag: KoTags("y", false)},
		),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "loop-monadics",
		File: `
		G(pass?) { return: pass }
		H(ignore?) { return: 1 }
		F(ignore?) { return: true }
		Main(x) {
			_1: Loop(start: x, step: G)
			_2: Loop(start: x, step: G, stop: F)
			_3: Loop(step: H, stop: F)
			return: (_1, _2, _3)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "tail-select",
		File: `
		Main(x) {
			_: (a: 1, a: 2)
			return: _.a
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "empties",
		File: `
		Empty() { return: empty() }
		empty(dontPass) { return: dontPass }
		Main() { return: Empty() }
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "range",
		File: `
		IgnoreCarryThruElem(elem) {
			return: (emit: elem)
		}
		Main(int8slice) {
			_1: Range(with: IgnoreCarryThruElem, start: "start") // GoEmpty{}
			_2: Range(over: int8slice, with: IgnoreCarryThruElem) // GoSlice{GoInt8}
			_3: Range(over: 64, with: IgnoreCarryThruElem) // lifted to GoSlice{GoInt64}
			return: (_1.image, _2.image, _3.image)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "Int8slice", Type: NewGoSlice(GoInt8), Tag: KoTags("int8slice", false)},
		),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "etc-selection",
		File: `
		Series(etc?) { // used to create series types, e.g. "Series(Bool(true))"
			series: (etc, etc)
			return: series
		}
		Main() {
			return: 1
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: true,
		Name:    "type-macro",
		File: `
		Main(x) {
			return: (
				optX: Optional(x)
				bool: Bool(true)
				int64: Int64(-2)
				float64: Float64(-2.1e4)
				repeatedString: Repeated(String(""))
			)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "Ko_X", Type: GoInt8, Tag: KoTags("x", false)},
		),
		Result: NewGoNeverNilPtr(
			NewGoStruct(
				TestGoField("span:q5c5b8i, test.ko:3:12", "optX", NewGoPtr(GoInt8)),
				TestGoField("span:q5c5b8i, test.ko:3:12", "bool", GoBool),
				TestGoField("span:q5c5b8i, test.ko:3:12", "int64", GoInt64),
				TestGoField("span:q5c5b8i, test.ko:3:12", "float64", GoFloat64),
				TestGoField("span:q5c5b8i, test.ko:3:12", "repeatedString", NewGoSlice(GoString)),
			),
		),
	},
	{
		Enabled: true,
		Name:    "fix-negation-number-aware-assign",
		File: `
		Return(etc?) { return: etc }
		Main() {
			fix: Fix(Return[Int8(-1)])
			return: fix(77)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: false, // test needs package integer, move to sys
		Name:    "fib",
		File: `
		import "boolean"
		import "integer"
		Fib(n) {
			return: Yield(
				if: boolean.Or(
					integer.Equal(n, 0)
					integer.Equal(n, 1)
				)
				then: Return[1]
				else: fibSum[n: n]
			)()
		}
		fibSum(n) {
			return: integer.Sum(
				Fib(n: integer.Sum(n, -1))
				Fib(n: integer.Sum(n, -2))
			)
		}
		Main() {
			return: Fix(
				Fib[n: Int64(0)]
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
}

func initTest() {
	RegisterGoGateAt("boolean", "Or", new(testBoolOr))
}

type testBoolOr struct {
	Etc []bool `ko:"name=etc,monadic"`
}

func (g *testBoolOr) Play(ctx *runtime.Context) bool {
	or := false
	for _, e := range g.Etc {
		or = or || e
	}
	return or
}
