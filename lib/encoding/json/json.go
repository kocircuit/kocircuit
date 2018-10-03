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

package json

import (
	"encoding/json"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Marshal", new(goMarshal))
	go_eval.RegisterNamedEvalGate("MarshalIndent", new(goMarshalIndent))
}

type goMarshal struct {
	Value interface{} `ko:"name=value,monadic"`
}

func (g *goMarshal) Play(ctx *runtime.Context) string {
	buf, err := json.Marshal(g.Value)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (g *goMarshal) Help() string {
	return "Marshal(value?)"
}

func (g *goMarshal) Doc() string {
	return "Marshal(value?) encodes the given value as JSON"
}

type goMarshalIndent struct {
	Value  interface{} `ko:"name=value,monadic"`
	Prefix *string     `ko:"name=prefix"`
	Indent *string     `ko:"name=indent"`
}

func (g *goMarshalIndent) Play(ctx *runtime.Context) string {
	buf, err := json.MarshalIndent(
		g.Value,
		util.OptString(g.Prefix, ""),
		util.OptString(g.Indent, "\t"),
	)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (g *goMarshalIndent) Help() string {
	return "MarshalIndent(value?, prefix, indent)"
}

func (g *goMarshalIndent) Doc() string {
	return "MarshalIndent(value?, prefix, indent) encodes the given value as indented JSON"
}
