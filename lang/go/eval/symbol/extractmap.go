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

package symbol

import (
	"reflect"

	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

func ExtractMap(span *model.Span, s Symbol, mapType reflect.Type) (Symbol, error) {
	if mapValue, err := extractMapValue(span, s, mapType); err != nil {
		return nil, err
	} else {
		return Deconstruct(span, mapValue), nil
	}
}

func extractMapValue(span *model.Span, s Symbol, mapType reflect.Type) (reflect.Value, error) {
	series := s.LiftToSeries(span)
	kt, vt := mapType.Key(), mapType.Elem()
	mapValue := reflect.MakeMapWithSize(mapType, len(series.Elem))
	for _, e := range series.Elem {
		row, ok := e.(*StructSymbol) // (key: K, value: V)
		if !ok {
			return reflect.Value{}, span.Errorf(nil, "expecting key/value struct")
		}
		kv, err := Integrate(span, row.Walk("key"), kt)
		if err != nil {
			return reflect.Value{}, span.Errorf(err, "expecting key of type %v", kt)
		}
		vv, err := Integrate(span, row.Walk("value"), vt)
		if err != nil {
			return reflect.Value{}, span.Errorf(err, "expecting value of type %v", vt)
		}
		mapValue.SetMapIndex(kv, vv)
	}
	return mapValue, nil
}
