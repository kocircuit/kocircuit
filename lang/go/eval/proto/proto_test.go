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
		Main() {
			_: pb.ProtoGoEnum(
				foo: pb.EnumFOO_FOO1()
			)
			return: pb.ProtoGoTest(
				Kind: pb.EnumGoTest_KIND_FUNCTION()
			)
		}
		`,
		Arg: nil,
		Result: struct {
			Ko_Kind int32 `ko:"name=Kind"`
		}{
			Ko_Kind: 12,
		},
	},
	{
		Enabled: true,
		File: `
		import "proto/testdata" as pb
		Main() {
			p1: pb.ProtoOldMessage(
				num: Int32(13)
				nested: pb.ProtoOldMessage(
					num: Int32(13)
				)
			)
			blob: pb.MarshalOldMessage(proto: p1)
			_: Show(len: Len(blob))
			p2: pb.UnmarshalOldMessage(bytes: blob)
			return: Equal(p1, p2)
		}
		`,
		Arg:    nil,
		Result: true,
	},
}
