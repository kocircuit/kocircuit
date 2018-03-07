package sys

import (
	"fmt"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"

	_ "github.com/kocircuit/kocircuit/lib/os"
)

func TestBuild(t *testing.T) {
	for _, test := range testBuild {
		if !test.Enabled {
			continue
		}
		b := BuildString{
			KoString:  test.File,
			KoPkg:     "test",
			KoFunc:    "Main",
			Faculty:   GoFaculty(),
			Idiom:     GoIdiomRepo,
			Arg:       test.Arg,
			Toolchain: NewGoToolchain(),
			GoKoRoot:  "kogo_test",
			GoKoPkg:   fmt.Sprintf("Test_%s", test.Name),
		}
		if result := b.Play(nil); result.Error != nil {
			t.Errorf("test %q: building (%v)", test.Name, result.Error)
			continue
		} else {
			t.Logf(result.Repo.BodyString())
			t.Logf("Result: %s\n", Sprint(result.Returns))
		}
	}
}

var testBuild = []struct {
	Enabled bool
	Name    string
	File    string    // compile the test file into a repo
	Arg     *GoStruct // arg for main
	Result  Shape     // weave main function, expecting result
}{
	{
		Enabled: false,
		Name:    "JoinAugmentGeneralize",
		File: `
		G(u, w) {
		   p: (a: u, b: w) // *struct{ a int8, b int32 }
		   q: (a: w, b: u) // *struct{ a int32, b int8 }
		   return: (
		      (p.a, q.b) // [2]int8
		      (p.b, q.a) // [2]int32
		   ) // [2][2]int32
		}
		Main(x, y) { return: G(u: x, w: y) }
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
			&GoField{Name: "Y", Type: GoInt32, Tag: KoTags("y", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "CallGate",
		File: `
		import "github.com/kocircuit/kocircuit/lib/os"
		Main() {
			return: os.GoTask(
				name: "ls_test"
				binary: "/bin/ls"
				arg: "/"
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "SpaceRecursion",
		File: `
		Return(x) {
			y: Return(x: x)
			return: 3
		}
		Main() { return: Return(x: 1) }
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "Yield",
		File: `
		Main(x, y, z) {
			y0: Yield(if: z, then: x, else: y)
			y1: Yield(if: true, then: x, else: y)
			y2: Yield(if: false, else: y)
			y3: Yield(if: false, then: x, _irrelevant_: 0)
			return: (y0, y1, y2, y3)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
			&GoField{Name: "Y", Type: GoInt16, Tag: KoTags("y", false)},
			&GoField{Name: "Z", Type: GoBool, Tag: KoTags("z", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "Loop",
		File: `
		G(pass?) { return: pass }
		H(ignore?) { return: 1 } // int64
		F(ignore?) { return: true } // bool
		Main(x) { // int8
			integers: (
				Loop(start: x, step: G) // int8
				Loop(start: x, step: G, stop: F) // int8
				Loop(start: x, step: H) // int64
				Loop(start: x, step: H, stop: F) // int64
				Loop(step: H, stop: F) // int64
				Loop(step: H) // int64
			)
			empties: (
				Loop(step: G, stop: F) // struct{}
				Loop(step: G) // struct{}
			)
			return: (integers: integers, empties: empties)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "Spin",
		File: `
		Return1(ignore?) { return: 1 }
		AlwaysStop(ignore?) { return: true }
		Main(x) {
			return: Spin(
				Loop[start: x, step: Return1, stop: AlwaysStop]
			)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "Range",
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
		Enabled: false,
		Name:    "SwitchCut",
		File: `
		Main(x, y, z) {
			return: Switch(
				case: (if: z, yield: x)
				case: (if: true, yield: x)
				case: (if: false, yield: y)
				otherwise: 77
			)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoInt8, Tag: KoTags("x", false)},
			&GoField{Name: "Y", Type: GoInt16, Tag: KoTags("y", false)},
			&GoField{Name: "Z", Type: GoBool, Tag: KoTags("z", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "NotAndOrHave",
		File: `
		Main(x) {
			_1: Not(true)
			_2: Not(false)
			_3: Not(x)
			_4: And(a: true, b: false)
			_5: And(c: x, d: false)
			_6: And(a: x, b: true)
			_7: Or(a: true, b: false)
			_8: Or(c: x, d: false)
			_9: Or(a: x, b: true)
			_0a: Have() // no argument is empty
			_0b: Have("a")
			return: (_1, _2, _3, _4, _5, _6, _7, _8, _9, _0a, _0b)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "X", Type: GoBool, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "FixServe",
		File: `
		Return(pass?) { return: pass }
		Foo(u, v) { return: (u: u, v: v) }
		Main(x) {
			fix: Fix(Return[Int64(0)])
			serve: Serve(
				func: Foo[
					u: Int64(0)
					v: String("ha")
				]
				path: "test.Foo"
			)
			return: serve(u: x, v: "ho")
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "Ko_X", Type: GoInt64, Tag: KoTags("x", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,            //XXX: disable
		Name:    "VarietyVariety", // XXX: this test hits a small bug in variety-variety assignment/generalization
		File: `
		YieldVariety(num, cond) {
			return: Yield(
				if: cond
				then: Return[true]
				else: Return[true, true]
			)
		}
		Main() {
			return: Serve(
				func: YieldVariety[num: Int64(0), cond: Bool(true)]
				path: "test.VarietyVariety"
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "Fib",
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
		Main(x) {
			return: Serve(
				func: Fib[n: Int64(0)]
				path: "test.F"
			)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "Ko_N", Type: GoInt64, Tag: KoTags("n", false)},
		),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "ServeLen",
		File: `
		ShowLen(series?) {
			return: Len(series)
		}
		Main() {
			return: Serve(
				func: ShowLen[ Repeated(Int64(0)) ]
				path: "test.Len"
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
	{
		Enabled: false,
		Name:    "NonRecursiveIdenticalInvocations",
		File: `
		import "github.com/kocircuit/kocircuit/lib/strings"
		Rewrite(s?) { return: strings.GoJoinStrings("_", s) }
		Main() {
			return: strings.GoJoinStrings(
				Rewrite("abc")
				Rewrite("abc")
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
}
