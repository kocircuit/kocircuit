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

func StructFieldIsProtoOptOrRep(sf reflect.StructField) bool {
	tag := sf.Tag.Get("protobuf")
	for _, kv := range strings.Split(tag, ",") {
		if kv == "opt" || kv == "rep" {
			return true
		}
	}
	return false
}

func StructFieldIsProtoOpt(sf reflect.StructField) bool {
	tag := sf.Tag.Get("protobuf")
	for _, kv := range strings.Split(tag, ",") {
		if kv == "opt" {
			return true
		}
	}
	return false
}

func structFieldProtoName(sf reflect.StructField) (string, bool) {
	tag := sf.Tag.Get("protobuf")
	for _, kv := range strings.Split(tag, ",") {
		x := strings.SplitN(kv, "=", 2)
		switch len(x) {
		case 2:
			if x[0] == "name" {
				return x[1], true
			}
		}
	}
	return "", false
}
