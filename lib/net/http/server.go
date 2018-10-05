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
	"encoding/json"
	"net/http"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Serve", new(goServe))
	go_eval.RegisterNamedEvalGate("SendJSON", new(goSendJSON))
}

// goServe implements the Serve(...) function
type goServe struct {
	// The desired http server address to listen on
	Address string `ko:"name=address"`
	// Handler called for every request
	Handler symbol.Symbol `ko:"name=handler"`
}

// HandlerContext is passed to a Serve handler.
type HandlerContext struct {
	// Requested Path
	Path string `ko:"path"`
	// Writer to send the response on.
	Writer http.ResponseWriter `ko:"writer"`
}

func (g *goServe) Play(ctx *runtime.Context) bool {
	ctx.Printf("serving %s", g.Address)
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hctx := &HandlerContext{
			Path:   req.URL.Path,
			Writer: w,
		}
		span := model.NewSpan()
		fields := eval.Fields{
			eval.Field{
				Name:  "ctx",
				Shape: symbol.DeconstructInterface(span, hctx),
			},
		}
		v, _, err := g.Handler.Augment(span, fields)
		if err != nil {
			ctx.Panicf("Augmenting handler failed: %v", err)
		}
		if _, _, err := v.Invoke(span); err != nil {
			ctx.Panicf("Handler invocation failed: %v", err)
		}
	})
	if err := http.ListenAndServe(g.Address, handler); err != nil {
		ctx.Fatalf("serving %s failed: %s", g.Address, err)
		return false // never reached
	}
	return true
}

func (g *goServe) Doc() string {
	return "Serve(address, handler) runs an HTTP server that invokes a handler upon incoming requests"
}

func (g *goServe) Help() string {
	return "Serve(address, handler)"
}

// goSendJSON implements the SendJSON(ctx, statusCode, value?) function
type goSendJSON struct {
	Context    *HandlerContext `ko:"ctx"`
	Value      interface{}     `ko:"value,monadic"`
	StatusCode *int            `ko:"statusCode"`
}

func (g *goSendJSON) Play(ctx *runtime.Context) bool {
	w := g.Context.Writer
	if w.Header().Get(headerKeyContentType) == "" {
		w.Header().Set(headerKeyContentType, contentTypeJSON)
	}
	w.WriteHeader(util.OptInt(g.StatusCode, 200))
	encoded, err := json.Marshal(g.Value)
	if err != nil {
		panic(err)
	}
	if _, err := w.Write(encoded); err != nil {
		panic(err)
	}
	return true
}

func (g *goSendJSON) Doc() string {
	return "SendJSON(ctx, statusCode, value?) sends the given value as JSON encoded text back to the client"
}

func (g *goSendJSON) Help() string {
	return "SendJSON(ctx, statusCode, value?)"
}
