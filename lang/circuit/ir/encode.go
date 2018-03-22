package ir

import (
	"bytes"
	"compress/gzip"

	"github.com/golang/protobuf/proto"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func EncodeRepo(repo Repo) ([]byte, error) {
	pbRepo := SerializeRepo(repo)
	buf, err := proto.Marshal(pbRepo)
	if err != nil {
		return nil, err
	}
	var w bytes.Buffer
	gw := gzip.NewWriter(&w)
	if _, err := gw.Write(buf); err != nil {
		return nil, err
	}
	if err = gw.Close(); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
