//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ir

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

func DecodeRepo(gzipped []byte) (model.Repo, error) {
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
