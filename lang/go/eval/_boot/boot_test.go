package boot

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/test"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func TestBoot(t *testing.T) {
	tests := &EvalTests{T: t, Test: bootTests}
	tests.Play(runtime.NewContext())
}

var bootTests = []*EvalTest{
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
			r: Boot(
				booter: (
					reserve: (pkg: "", name: "")
					reserve: (pkg: "", name: "Reserved")
					reserve: (pkg: "reserve_pkg1", name: "Reserved1")
				)
				func: G
				ctx: (bootUserCtx: true)
				arg: (bootUserArg: true)
			)
			return: Show(r)
		}
		`,
		Arg: struct {
			Ko_x byte `ko:"name=x"`
		}{
			Ko_x: 7,
		},
		Result: nil,
	},
}
