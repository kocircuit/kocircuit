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

package http

import (
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("SetHeader", new(goSetHeader))
}

const (
	headerKeyContentType = "Content-Type"
	contentTypeJSON      = "text/json"
)

// HeaderEntry holds a key/value pair passed to the response header of an HTTP request
type HeaderEntry struct {
	Key   string `ko:"key"`
	Value string `ko:"value"`
}

// goSetHeader implements the SetHeader(ctx, entry) function
type goSetHeader struct {
	Context *HandlerContext `ko:"ctx"`
	Entry   []HeaderEntry   `ko:"entry"`
}

func (g *goSetHeader) Play(ctx *runtime.Context) *HandlerContext {
	for _, e := range g.Entry {
		g.Context.Writer.Header().Set(e.Key, e.Value)
	}
	return g.Context
}

func (g *goSetHeader) Doc() string {
	return "SetHeader(ctx, entry) sends the given header entry (key, value) into the given handler context and returns the handler context"
}

func (g *goSetHeader) Help() string {
	return "SetHeader(ctx, entry)"
}
