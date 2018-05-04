package source

import (
	pb "github.com/kocircuit/kocircuit/bootstrap/source/proto"
	eval_proto "github.com/kocircuit/kocircuit/lang/go/eval/proto"
)

func init() {
	eval_proto.RegisterEvalProtoFileBytes(pb.FileDescriptorBytes())
}
