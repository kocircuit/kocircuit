package test

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func TestEvalIntegrations(t *testing.T) {
	RegisterEvalGateAt("test", "MapGate", new(testMapGate))
	tests := &EvalTests{T: t, Test: evalIntegrationTests}
	tests.Play(runtime.NewContext())
}

var evalIntegrationTests = []*EvalTest{
	{ // maps
		Enabled: true,
		File: `
		import "test"
		Main(m) { 
			return: test.MapGate(map: m)
		}
		`,
		Arg: struct {
			Ko_m map[string]int32 `ko:"name=m"`
		}{
			Ko_m: map[string]int32{"a": 1, "b": 2},
		},
		Result: map[string]int64{"a": 1, "b": 2},
	},
}

type testMapGate struct {
	Map map[string]int64 `ko:"name=map"`
}

func (g *testMapGate) Play(ctx *runtime.Context) map[string]int64 {
	return g.Map
}
