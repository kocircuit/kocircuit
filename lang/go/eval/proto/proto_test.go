package proto

import (
	"testing"

	_ "github.com/kocircuit/kocircuit/lang/go/eval/proto/testdata"

	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/test"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func TestProto(t *testing.T) {
	RegisterEvalProtoFile("test.proto")
	tests := &EvalTests{T: t, Test: protoTests}
	tests.Play(runtime.NewContext())
}

var protoTests = []*EvalTest{
	{
		Enabled: true,
		File: `
		import "proto/testdata" as pb
		Main(x, y) {
			_: pb.ProtoGoEnum(
				foo: pb.EnumFOO_FOO1()
			)
			return: pb.ProtoGoTest(
				Kind: pb.EnumGoTest_KIND_FUNCTION()
			)
		}
		`,
		Arg: struct { // deconstruct/construct
			Ko_x byte    `ko:"name=x"`
			Ko_y float64 `ko:"name=y"`
		}{
			Ko_x: 7,
			Ko_y: 3.3,
		},
		Result: struct {
			Ko_Kind int32 `ko:"name=Kind"`
		}{
			Ko_Kind: 12,
		},
	},
}
