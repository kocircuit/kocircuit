package gate

import (
	"reflect"
	"strings"
)

func StructFieldProtoName(sf reflect.StructField) (string, bool) {
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
