package ir

import (
	"bytes"
	"compress/gzip"

	"github.com/golang/protobuf/proto"

	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func EncodeRepo(repo Repo) ([]byte, error) {
	if _, buf, err := SerializeEncodeRepo(repo); err != nil {
		return nil, err
	} else {
		return buf, nil
	}
}

func SerializeEncodeRepo(repo Repo) (*pb.Repo, []byte, error) {
	pbRepo := SerializeRepo(repo)
	buf, err := proto.Marshal(pbRepo)
	if err != nil {
		return nil, nil, err
	}
	var w bytes.Buffer
	gw := gzip.NewWriter(&w)
	if _, err := gw.Write(buf); err != nil {
		return nil, nil, err
	}
	if err = gw.Close(); err != nil {
		return nil, nil, err
	}
	return pbRepo, w.Bytes(), nil
}
