package ir

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func DecodeRepo(gzipped []byte) (Repo, error) {
	r := bytes.NewBuffer(gzipped)
	gr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	if err := gr.Close(); err != nil {
		return nil, err
	}
	pbRepo := &pb.Repo{}
	if err := proto.Unmarshal(buf, pbRepo); err != nil {
		return nil, err
	}
	return DeserializeRepo(pbRepo)
}
