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
	go_json "encoding/json"
	"fmt"
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Unmarshal", new(goUnmarshal))
}

type goUnmarshal struct {
	Value string `ko:"name=value,monadic"`
}

func (g *goUnmarshal) Play(ctx *runtime.Context) (symbol.Symbol, error) {
	return Unmarshal([]byte(g.Value))
}

func (g *goUnmarshal) Help() string {
	return "Unmarshal(value?)"
}

func (g *goUnmarshal) Doc() string {
	return "Unmarshal(value?) decodes the given JSON value"
}

// Unmarshal the given encoded JSON to a Go value.
func Unmarshal(encoded []byte) (symbol.Symbol, error) {
	var result interface{}
	if err := go_json.Unmarshal(encoded, &result); err != nil {
		return nil, err
	}
	return symbol.Deconstruct(model.NewSpan(), mapsToStructs(reflect.ValueOf(result))), nil
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
