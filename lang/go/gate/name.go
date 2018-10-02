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

package gate

import (
	"reflect"
	"strings"
)

type KoGoName struct {
	Ko string `ko:"name=ko"`
	Go string `ko:"name=go"`
}

func StructFieldWithNoKoOrProtoName(sf reflect.StructField) bool {
	if _, hasProto := structFieldProtoName(sf); hasProto {
		return false
	} else if _, hasKo := structFieldProtoName(sf); hasKo {
		return false
	} else {
		return true
	}
}

func StructFieldKoProtoGoName(sf reflect.StructField) (string, bool) {
	if koName, ok := structFieldKoName(sf); ok {
		return koName, true
	}
	if protoName, ok := structFieldProtoName(sf); ok {
		return protoName, true
	}
	return structFieldGoName(sf)
}

// structFieldKoName looks up the Ko name of the given field.
// Returns: KoName, NameFound
func structFieldKoName(sf reflect.StructField) (string, bool) {
	splittedTag := strings.Split(sf.Tag.Get("ko"), ",")
	for _, kv := range splittedTag {
		x := strings.SplitN(kv, "=", 2)
		switch len(x) {
		case 2:
			if x[0] == "name" {
				return x[1], true
			}
		}
	}
	// No explicit `name=Foo` found, use first entry
	name := splittedTag[0]
	return name, name != ""
}

func structFieldGoName(sf reflect.StructField) (string, bool) {
	if sf.PkgPath != "" { // unexported field
		return "", false
	} else {
		return sf.Name, true
	}
}
