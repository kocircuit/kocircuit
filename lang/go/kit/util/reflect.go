package util

import (
	"fmt"
	"reflect"
)

func InterfaceTypeAddress(v interface{}) string {
	return TypeAddress(reflect.TypeOf(v))
}

func TypeAddress(t reflect.Type) string {
	return fmt.Sprintf("go:%q.%s", t.PkgPath(), t.Name())
}
