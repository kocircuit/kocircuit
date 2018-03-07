package hash

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"strconv"
	"strings"
)

const ScopeHashLen = 7

func Mix(s ...string) string {
	return MixN(ScopeHashLen, s...)
}

func MixInterface(f ...interface{}) string {
	return MixNInterface(ScopeHashLen, f...)
}

func MixScalar(v interface{}) string {
	return MixValue(reflect.ValueOf(v))
}

func MixValue(v reflect.Value) string {
	return MixNValue(ScopeHashLen, v)
}

func MixNValue(n int, v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		// reflect.Value.String handles Strings differently than the other kinds
		return MixN(n, "string", v.String())
	case reflect.Bool:
		return MixN(n, "bool", strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return MixN(n, "integer", strconv.FormatInt(v.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return MixN(n, "unsigned-integer", strconv.FormatUint(v.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		return MixN(n, "floating", strconv.FormatFloat(v.Float(), 'e', -1, 64))
	case reflect.UnsafePointer:
		return MixN(n, "unsafe", strconv.FormatUint(uint64(v.UnsafeAddr()), 16))
	case reflect.Complex64, reflect.Complex128:
		return MixN(n, "complex", fmt.Sprintf("%v", v.Interface()))
	case reflect.Invalid:
		return MixN(n, "invalid")
	case reflect.Chan:
		panic("channel hashing not supported")
	case reflect.Func:
		panic("func hashing not supported")
	}
	panic("o")
}

func MixNInterface(n int, s ...interface{}) string {
	ss := make([]string, len(s))
	for i := range s {
		id := reflect.ValueOf(s[i]).InterfaceData()
		ss[i] = fmt.Sprintf("%v.%v", id[0], id[1])
	}
	return MixN(n, ss...)
}

func MixN(n int, s ...string) string {
	h := fnv.New64a() // alloc *uint64
	for _, s := range s {
		h.Write([]byte(s)) // fnv is as fast as direct arithmetic
	}
	u := h.Sum64()
	trim := strconv.FormatUint(u, 36)
	if len(trim) >= n {
		return trim[:n]
	}
	return strings.Join([]string{strings.Repeat("0", n-len(trim)), trim}, "")
}
