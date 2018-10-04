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
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("BodyAsString", new(goBodyAsString))
	go_eval.RegisterNamedEvalGate("BodyAsJSON", new(goBodyAsJSON))
}

// Response is the response for an HTTP request.
type Response struct {
	StatusCode int                  `ko:"statusCode"`
	Response   *symbol.OpaqueSymbol // Wrapping a *http.Response
}

func (r *Response) getHTTPResponse() (*http.Response, error) {
	if r == nil || r.Response == nil {
		return nil, fmt.Errorf("response not set")
	}
	if httpResp, ok := r.Response.Value.Interface().(*http.Response); ok {
		return httpResp, nil
	}
	return nil, fmt.Errorf("response has wrong type")
}

// newResponse creates a new Response from the given HTTP response.
func newResponse(resp *http.Response) *Response {
	return &Response{
		Response:   &symbol.OpaqueSymbol{Value: reflect.ValueOf(resp)},
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

// goBodyAsJSON implements the BodyAsJSON function.
type goBodyAsJSON struct {
	Response *Response `ko:"response,monadic"`
}

func (g *goBodyAsJSON) Play(ctx *runtime.Context) (symbol.Symbol, error) {
	httpResp, err := g.Response.getHTTPResponse()
	if err != nil {
		return nil, err
	}
	body := httpResp.Body
	if body == nil {
		return nil, err
	}
	defer body.Close()
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	var value interface{}
	if err := json.Unmarshal(content, &value); err != nil {
		return nil, err
	}
	return symbol.Deconstruct(model.NewSpan(), mapsToStructs(reflect.ValueOf(value))), nil
}

func (g *goBodyAsJSON) Doc() string {
	return "BodyAsJSON(response?) returns the body of the given response as a JSON value"
}

func (g *goBodyAsJSON) Help() string {
	return "BodyAsJSON(response?)"
}

// mapsToStructs recursively converts map[string]interface{} to structs.
func mapsToStructs(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Map:
		if v.Type().Key().Kind() == reflect.String && v.Type().Elem().Kind() == reflect.Interface {
			// Conversion needed
			keys := v.MapKeys()
			// Create struct fields
			fields := make([]reflect.StructField, 0, len(keys))
			fieldValues := make([]reflect.Value, 0, len(keys))
			for i, k := range keys {
				fv := v.MapIndex(k)
				fv = mapsToStructs(fv)
				fields = append(fields, reflect.StructField{
					Name: fmt.Sprintf("Field%d", i),
					Type: fv.Type(),
					Tag:  reflect.StructTag(fmt.Sprintf(`ko:"name=%s" json:"%s"`, k.String(), k.String())),
				})
				fieldValues = append(fieldValues, fv)
			}
			structType := reflect.StructOf(fields)
			vsRef := reflect.New(structType)
			vs := vsRef.Elem()
			// Fill struct fields
			for i := range keys {
				vs.Field(i).Set(fieldValues[i])
			}
			return vsRef
		}
		// No conversion needed, just recurse into all fields
		keys := v.MapKeys()
		for _, k := range keys {
			fv := v.MapIndex(k)
			fv = mapsToStructs(fv)
			v.SetMapIndex(k, fv)
		}
		return v
	case reflect.Interface:
		if v.IsNil() {
			return v
		}
		return mapsToStructs(v.Elem())
	case reflect.Ptr:
		if v.IsNil() {
			return v
		}
		return mapsToStructs(v.Elem()).Addr()
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			sv := v.Index(i)
			sv = mapsToStructs(sv)
			v.Index(i).Set(sv)
		}
		return v
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			fv = mapsToStructs(fv)
			v.Field(i).Set(fv)
		}
		return v
	default:
		return v
	}
}
