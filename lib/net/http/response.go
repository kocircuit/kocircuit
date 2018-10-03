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
	"io/ioutil"
	"net/http"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("BodyAsString", new(goBodyAsString))
}

// Response is the response for an HTTP request.
type Response struct {
	StatusCode int `ko:"statusCode"`
	Response   interface{}
}

func (r *Response) getHTTPResponse() (*http.Response, error) {
	if r == nil || r.Response == nil {
		return nil, fmt.Errorf("response not set")
	}
	if httpResp, ok := r.Response.(*http.Response); ok {
		return httpResp, nil
	}
	return nil, fmt.Errorf("response has wrong type")
}

// newResponse creates a new Response from the given HTTP response.
func newResponse(resp *http.Response) *Response {
	return &Response{
		Response:   resp,
		StatusCode: resp.StatusCode,
	}
}

// goBodyAsString implements the BodyAsString function.
type goBodyAsString struct {
	Response *Response `ko:"response,monadic"`
}

func (g *goBodyAsString) Play(ctx *runtime.Context) (string, error) {
	httpResp, err := g.Response.getHTTPResponse()
	if err != nil {
		return "", err
	}
	body := httpResp.Body
	if body == nil {
		return "", nil
	}
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (g *goBodyAsString) Doc() string {
	return "BodyAsString(response?) returns the body of the given response as a string"
}

func (g *goBodyAsString) Help() string {
	return "BodyAsString(response?)"
}
