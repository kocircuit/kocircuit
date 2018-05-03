package types

import (
	pb "github.com/kocircuit/kocircuit/bootstrap/types/proto"
	eval_proto "github.com/kocircuit/kocircuit/lang/go/eval/proto"
)

func init() {
	eval_proto.RegisterEvalProtoFileBytes(pb.FileDescriptorBytes())
}
