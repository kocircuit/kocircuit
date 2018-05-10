package weave

import (
	"testing"

	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/test"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func TestWeave(t *testing.T) {
	tests := &EvalTests{T: t, Test: weaveTests}
	tests.Play(runtime.NewContext())
}

var weaveTests = []*EvalTest{
	{
		Enabled: true,
		File: `
		A() { return: () }
		B() { return: () }
		C() { return: () }
		D() { return: () }
		H() { return: () }
		G(u, w) {
			p: A(a: u, b: w)
			q: B(a: w, b: u)
			return: H(
				s: C(p.a, q.b)
				t: D(p.b, q.a)
			)
		}
		Forward(stepCtx, object) {
			return: (
				returns: object
				effect: (effect: stepCtx)
			)
		}
		Literal(stepCtx, figure) {
			return: (
				returns: (figure: figure)
				effect: (effect: stepCtx)
			)
		}
		Combine(summaryCtx, stepResidues) {
			return: (
				effect: (
					summaryCtx: summaryCtx
					stepResidues: stepResidues
				)
			)
		}
		Main(x) {
			r: Weave(
				weaver: (
					operator: (pkg: "", name: "")
					operator: (pkg: "", name: "OperatorA")
					operator: (pkg: "operator_pkg", name: "OperatorB")
					Enter: Forward
					Leave: Forward
					Link: Forward
					Select: Forward
					Augment: Forward
					Invoke: Forward
					Literal: Literal
					Combine: Combine
				)
				func: G
				ctx: (weaveUserCtx: true)
				arg: (weaveUserArg: true)
			)
			return: Equal(r, expected())
		}` + testExpected,
		Arg: struct {
			Ko_x byte `ko:"name=x"`
		}{
			Ko_x: 7,
		},
		Result: true,
	},
}

const testExpected = `
expected() {
	return: (
		returns: (
			figure: WeaveFigure(
				transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
			)
		)
		effect: (
			summaryCtx: (
				pkg: "test"
				func: "G"
				source: "test.ko:7:3"
				ctx: (weaveUserCtx: true)
				arg: (weaveUserArg: true)
				returns: (
					figure: WeaveFigure(
						transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
					)
				)
			)
			stepResidues: (
				(
					step: "0_leave"
					logic: "LEAVE"
					source: "test.ko:7:3"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "0_leave"
							logic: "LEAVE"
							source: "test.ko:7:3"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "return"
					logic: "INVOKE"
					source: "test.ko:10:12"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "return"
							logic: "INVOKE"
							source: "test.ko:10:12"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "16"
					logic: "AUGMENT"
					source: "test.ko:10:12"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "16"
							logic: "AUGMENT"
							source: "test.ko:10:12"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "5"
					logic: "\"test\".H"
					source: "test.ko:10:12"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "H"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "5"
							logic: "\"test\".H"
							source: "test.ko:10:12"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "10"
					logic: "INVOKE"
					source: "test.ko:11:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "C"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "10"
							logic: "INVOKE"
							source: "test.ko:11:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "15"
					logic: "INVOKE"
					source: "test.ko:12:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "D"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "15"
							logic: "INVOKE"
							source: "test.ko:12:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "9"
					logic: "AUGMENT"
					source: "test.ko:11:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "C"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "9"
							logic: "AUGMENT"
							source: "test.ko:11:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "14"
					logic: "AUGMENT"
					source: "test.ko:12:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "D"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "14"
							logic: "AUGMENT"
							source: "test.ko:12:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "6"
					logic: "\"test\".C"
					source: "test.ko:11:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "C"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "6"
							logic: "\"test\".C"
							source: "test.ko:11:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "7"
					logic: "SELECT(a)"
					source: "test.ko:11:10"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "A"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "7"
							logic: "SELECT(a)"
							source: "test.ko:11:10"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "8"
					logic: "SELECT(b)"
					source: "test.ko:11:15"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "B"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "8"
							logic: "SELECT(b)"
							source: "test.ko:11:15"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "11"
					logic: "\"test\".D"
					source: "test.ko:12:8"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "D"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "11"
							logic: "\"test\".D"
							source: "test.ko:12:8"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "12"
					logic: "SELECT(b)"
					source: "test.ko:12:10"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "A"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "12"
							logic: "SELECT(b)"
							source: "test.ko:12:10"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "13"
					logic: "SELECT(a)"
					source: "test.ko:12:15"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "B"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "13"
							logic: "SELECT(a)"
							source: "test.ko:12:15"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "p"
					logic: "INVOKE"
					source: "test.ko:8:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "A"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "p"
							logic: "INVOKE"
							source: "test.ko:8:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "q"
					logic: "INVOKE"
					source: "test.ko:9:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "B"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "q"
							logic: "INVOKE"
							source: "test.ko:9:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "2"
					logic: "AUGMENT"
					source: "test.ko:8:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "A"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "2"
							logic: "AUGMENT"
							source: "test.ko:8:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "4"
					logic: "AUGMENT"
					source: "test.ko:9:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "B"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "4"
							logic: "AUGMENT"
							source: "test.ko:9:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "1"
					logic: "\"test\".A"
					source: "test.ko:8:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "A"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "1"
							logic: "\"test\".A"
							source: "test.ko:8:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "3"
					logic: "\"test\".B"
					source: "test.ko:9:7"
					returns: (
						figure: WeaveFigure(
							transform: WeaveTransform(func: WeaveFunc(pkg: "test", name: "B"))
						)
					)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "3"
							logic: "\"test\".B"
							source: "test.ko:9:7"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "0_enter_w"
					logic: "ARG(w)"
					source: "test.ko:7:3"
					returns: (weaveUserArg: true)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "0_enter_w"
							logic: "ARG(w)"
							source: "test.ko:7:3"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "0_enter_u"
					logic: "ARG(u)"
					source: "test.ko:7:3"
					returns: (weaveUserArg: true)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "0_enter_u"
							logic: "ARG(u)"
							source: "test.ko:7:3"
							ctx: (weaveUserCtx: true)
						)
					)
				)
				(
					step: "0_enter"
					logic: "ENTER"
					source: "test.ko:7:3"
					returns: (weaveUserArg: true)
					effect: (
						effect: (
							pkg: "test"
							func: "G"
							step: "0_enter"
							logic: "ENTER"
							source: "test.ko:7:3"
							ctx: (weaveUserCtx: true)
						)
					)
				)
			)
		)
	)
}
`
