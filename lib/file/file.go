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

package file

import (
	"io/ioutil"
	"os"
	"path"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("WriteLocalFile", new(goWriteLocalFile))
}

// goWriteLocalFile implements the WriteLocalFile(path, body) function
type goWriteLocalFile struct {
	Path string `ko:"name=path"`
	Body string `ko:"name=body"`
}

func (g *goWriteLocalFile) Play(ctx *runtime.Context) bool {
	if err := os.MkdirAll(path.Dir(g.Path), 0755); err != nil {
		ctx.Fatalf("mkdir %q (%v)", path.Dir(g.Path), err)
		return false
	}
	if err := ioutil.WriteFile(g.Path, []byte(g.Body), 0666); err != nil {
		ctx.Fatalf("write file %q (%v)", g.Path, err)
		return false
	}
	return true
}
