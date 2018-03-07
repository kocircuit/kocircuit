package test

import (
	"reflect"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestEval(t *testing.T) {
	for i, test := range testEval {
		if !test.Enabled {
			continue
		}
		repo, err := CompileString("test", "test.ko", test.File)
		if err != nil {
			t.Errorf("test %d: compile (%v)", i, err)
			continue
		}
		// fmt.Println(repo["test"].BodyString())
		eval := NewEvaluator(EvalFaculty(), repo)
		span := NewSpan()
		span = RefineChamber(span, "testEval")
		span = RefineOutline(span, "Main")
		// deconstruct test arguments
		arg, err := Deconstruct(span, reflect.ValueOf(test.Arg))
		if err != nil {
			t.Errorf("test %d: argument deconstruction (%v)", i, err)
			continue
		}
		returned, _, err := eval.Eval(span, repo["test"]["Main"], arg)
		if err != nil {
			t.Errorf("test %d: eval (%v)", i, err)
			continue
		}
		// integrate test results
		result, err := Integrate(span, returned, reflect.TypeOf(test.Result))
		if err != nil {
			t.Errorf("test %d: result integration (%v)", i, err)
			continue
		}
		// test result values
		if !IsSubset(test.Result, result.Interface()) {
			t.Errorf("test %d: expecting %s, got %v", i, Sprint(test.Result), result.Interface())
		}
	}
}

var testEval = []struct {
	Enabled bool
	File    string      // compile the test file into a repo
	Arg     interface{} // arg for main
	Result  interface{} // weave main function, expecting result
}{
	{ // select, repeat
		Enabled: true,
		File: `
		G(u, w) {
			p: (a: u, b: w)
			q: (a: w, b: u)
			return: (
				s: (p.a, q.b) // []byte
				t: (p.b, q.a) // []float64
			)
		}
		Main(x, y) { return: G(u: x, w: y) }
		`,
		Arg: struct { // deconstruct/construct
			Ko_x byte    `ko:"name=x"`
			Ko_y float64 `ko:"name=y"`
		}{
			Ko_x: 7,
			Ko_y: 3.3,
		},
		Result: struct {
			Ko_s []byte    `ko:"name=s"`
			Ko_t []float64 `ko:"name=t"`
		}{
			Ko_s: []byte{7, 7},
			Ko_t: []float64{3.3, 3.3},
		},
	},
	{ // select and repeated
		Enabled: true,
		File: `
		Main(x) {
			_: (a: 1, a: 2, a: x)
			return: _.a
		}
		`,
		Arg: struct {
			Ko_x int64 `ko:"name=x"`
		}{
			Ko_x: 7,
		},
		Result: []int64{1, 2, 7},
	},
	{ // empty
		Enabled: true,
		File: `
		Empty() { return: empty() }
		empty(dontPass) { return: dontPass }
		Main() { return: Empty() }
		`,
		Arg:    struct{}{},
		Result: (*struct{})(nil),
	},
	{ // boolean gates
		Enabled: true,
		File: `
		Main() { 
			return: Xor(
				true
				And(
					Or(true, false)
					Or(Xor(true, true), false)
				)
			)
		}`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // augmentation and invocation
		Enabled: true,
		File: `
		Pass(p?) { return: p }
		Main(x) {
			callMe: Peek(Pass[x])
			return: callMe()
		}
		`,
		Arg: struct {
			Ko_x byte `ko:"name=x"`
		}{
			Ko_x: 7,
		},
		Result: byte(7),
	},
	{ // integer macros
		Enabled: true,
		File: `
		import "integer"
		Main() {
			return: And(
				integer.Equal(-1, integer.Negative(1))
				Not(integer.Less(3, 2))
				integer.Equal(
					integer.Prod(2, 3)
					integer.Sum(1, 5)
				)
				integer.Equal(integer.Moduli(7, 5), 2)
				integer.Equal(integer.Ratio(7, 2), 3)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // len, have
		Enabled: true,
		File: `
		Main() {
			return: And(
				Not(Have())
				Have("a")
				Equal(Len(5), 1)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // equal
		Enabled: true,
		File: `
		Main() {
			return: And(
				Not(Equal("a", "a", "c"))
				Equal((1), 1)
				Equal(
					(a: "a", b: 1)
					(a: "a", b: (1))
				)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // Int8, ..., Uint8, ..., Float32, ...
		Enabled: true,
		File: `
		Main() {
			return: And(
				Equal(Int8(0), Int8(256))
				Equal(Uint8(0), Uint8(256))
				Equal(Int8(255), Int8(-1))
				Equal(Float32(-3.14e7), Float32(-3.14e7))
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // yield
		Enabled: true,
		File: `
		Main() {
			return: And(
				Equal(
					Yield(if: false, else: (1))
					1
				)
				Equal(
					Len(Yield(if: true, then: (1, 2, 3)))
					3
				)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // hash
		Enabled: true,
		File: `
		import "integer"
		NotEqual(x?) { 
			return: Not(Equal(x))
		}
		Main() {
			return: And(
				NotEqual(Hash(1, 2), Hash(1, 2, 3))
				Equal(Hash(a: "a", b: 1), Hash(a: "a", b: 1))
				Equal(Hash(a: "a", b: ()), Hash(a: "a"))
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // range
		Enabled: true,
		File: `
		import "integer"
		Main() {
			return: And(
				Equal(
					Range(
						start: 0
						over: (1, 2, 3, 4, 5)
						with: sum(carry, elem) {
							return: (
								emit: integer.Negative(elem)
								carry: integer.Sum(carry, elem)
							)
						}
					)
					(image: (-1, -2, -3, -4, -5), residue: 15)
				)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // loop
		Enabled: true,
		File: `
		import "integer"
		Main() {
			return: And(
				Equal(
					Loop(
						start: 0
						step: step(carry?) { return: integer.Sum(carry, 1) }
						stop: stop(carry?) { return: Equal(carry, 33) }
					)
					33
				)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // merge
		Enabled: true,
		File: `
		Main() {
			return: And(
				Equal(
					Merge((1, 2, 3), (11, 13, 17))
					(1, 2, 3, 11, 13, 17)
				)
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // spin, wait
		Enabled: true,
		File: `
		Main() {
			return: Spin(Peek["Hello, world!"]).Wait()
		}
		`,
		Arg: struct{}{},
		Result: struct {
			Ko_returned string `ko:"name=returned"`
			Ko_success  bool   `ko:"name=success"`
		}{
			Ko_returned: "Hello, world!",
			Ko_success:  true,
		},
	},
	{ // switch, take
		Enabled: true,
		File: `
		Main() {
			s: (1, 2, 3, 4, 5)
			return: Equal(
				Switch(
					case: (if: false, yield: "abc")
					case: (if: Equal(Len(s), 5), yield: "def")
					otherwise: "ghi"
				)
				"def"
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // switch, take
		Enabled: true,
		File: `
		import "integer"
		F(x?) { return: integer.Sum(x ,2) }
		G(x?) { return: integer.Prod(x ,2) }
		Main() {
			return: Equal(
				integer.Sum(
					Parallel(F[3]) // 5
					Sequential(G[5]) // 10
				)
				15
			)
		}
		`,
		Arg:    struct{}{},
		Result: true,
	},
	{ // format macro
		Enabled: true,
		File: `
		Main(x) {
			return: Format(
				format: "abc %% def %% hij %%"
				args: (x, "Y", "Z")
				withString: Return
				withArg: Return
			)
		}
		`,
		Arg: struct {
			Ko_x string `ko:"name=x"`
		}{
			Ko_x: "X",
		},
		Result: []string{"abc ", "X", " def ", "Y", " hij ", "Z", ""},
	},
	{ // recover/panic macro
		// TODO: test recovery from parallel invocations
		Enabled: true,
		File: `
		Main(x) {
			return: Recover(
				invoke: panicOnAll[x]
				panic: Return
			)
		}
		panicOnAll(x?) {
			return: Panic(x)
		}
		`,
		Arg: struct {
			Ko_x string `ko:"name=x"`
		}{
			Ko_x: "msg",
		},
		Result: "msg",
	},
}
