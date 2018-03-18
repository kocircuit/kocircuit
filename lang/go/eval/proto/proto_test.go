package proto

import (
	"testing"

	_ "github.com/kocircuit/kocircuit/lang/go/eval/proto/testdata"
)

func TestProto(t *testing.T) {
	RegisterEvalProtoFile("test.proto")
}
