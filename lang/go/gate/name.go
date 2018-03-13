package gate

import (
	"reflect"
	"strings"
)

type KoGoName struct {
	Ko string `ko:"name=ko"`
	Go string `ko:"name=go"`
}

func StructFieldKoProtoGoName(sf reflect.StructField) (string, bool) {
	if koName, ok := StructFieldKoName(sf); ok {
		return koName, true
	}
	if protoName, ok := StructFieldProtoName(sf); ok {
		return protoName, true
	}
	return StructFieldGoName(sf)
}

func StructFieldKoName(sf reflect.StructField) (string, bool) {
	tag := sf.Tag.Get("ko")
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

func StructFieldGoName(sf reflect.StructField) (string, bool) {
	if sf.PkgPath != "" { // unexported field
		return "", false
	} else {
		return sf.Name, true
	}
}
