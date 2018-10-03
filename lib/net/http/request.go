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
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Request", new(goRequest))
}

// goRequest implements the Request function.
type goRequest struct {
	// Client to create the request on. If not set, the default client is used.
	Client *Client `ko:"client"`
	// Method of the request
	Method *string `ko:"method"`
	// URL to request
	URL string `ko:"url"`
	// Body to send with the request
	Body interface{} `ko:"body"`
}

func (g *goRequest) Play(ctx *runtime.Context) (*Response, error) {
	client := g.Client
	if client == nil {
		client = getDefaultClient()
	}
	method := util.OptString(g.Method, "GET")
	body, err := g.createBodyReader()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, g.URL, body)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return newResponse(resp), nil
}

func (g *goRequest) Doc() string {
	return "Request(client, method, url, body) performs an HTTP request on the given (or default) client"
}

func (g *goRequest) Help() string {
	return "Request(client, method, url, body)"
}

// createBodyReader creates a Reader to fetch the body, depending on
// the type of the given body.
func (g *goRequest) createBodyReader() (io.Reader, error) {
	if g.Body == nil {
		return nil, nil
	}
	if x, ok := g.Body.(string); ok {
		return strings.NewReader(x), nil
	}
	return nil, fmt.Errorf("Unknown Body type %s", reflect.TypeOf(g.Body).String())
}
