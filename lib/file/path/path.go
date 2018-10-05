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

package path

import (
	"path"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Join", new(goJoinPath))
	go_eval.RegisterNamedEvalGate("Base", new(goBasePath))
	go_eval.RegisterNamedEvalGate("Dir", new(goDirPath))
}

type goJoinPath struct {
	Paths []string `ko:"name=paths,monadic"`
}

func (g *goJoinPath) Play(ctx *runtime.Context) string {
	return path.Join(g.Paths...)
}

type goBasePath struct {
	Path string `ko:"name=path,monadic"`
}

func (g *goBasePath) Play(ctx *runtime.Context) string {
	return path.Base(g.Path)
}

type goDirPath struct {
	Path string `ko:"name=path,monadic"`
}

func (g *goDirPath) Play(ctx *runtime.Context) string {
	return path.Dir(g.Path)
}
