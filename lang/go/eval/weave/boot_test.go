package weave

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
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
		G(u, w) {
			p: (a: u, b: w)
			q: (a: w, b: u)
			return: (
				s: (p.a, q.b)
				t: (p.b, q.a)
			)
		}
		Main(x) {
			r: Weave(
				weaver: (
					reserve: (pkg: "", name: "")
					reserve: (pkg: "", name: "Reserved")
					reserve: (pkg: "reserve_pkg1", name: "Reserved1")
				)
				func: G
				ctx: (weaveUserCtx: true)
				arg: (weaveUserArg: true)
			)
			return: Show(parsed: (result: r))
		}
		`,
		Arg: struct {
			Ko_x byte `ko:"name=x"`
		}{
			Ko_x: 7,
		},
		Result: struct{}{},
	},
}
