package model

import (
	"hash/fnv"
	"reflect"
	"strconv"
	"strings"
)

func PtrID(id ID) *ID {
	v := id
	return &v
}

type ID struct {
	Data uint64 `ko:"name=data"`
}

func IDFromProtoData(d uint64) ID {
	return ID{Data: d}
}

func (id ID) ProtoData() uint64 {
	return id.Data
}

func (id ID) String() string {
	return id.Format(7)
}

func (id ID) Format(width int) string {
	trim := strconv.FormatUint(id.Data, 36)
	if len(trim) >= width {
		return trim[:width]
	}
	return strings.Join([]string{strings.Repeat("0", width-len(trim)), trim}, "")
}

func InterfaceID(v interface{}) ID {
	return ValueID(reflect.ValueOf(v))
}

func ValueID(v reflect.Value) ID {
	return Blend(
		StringID(v.Type().String()),
		StringID(v.String()),
	)
}

func StringID(s string) ID {
	h := fnv.New64a() // alloc *uint64
	h.Write([]byte(s))
	return ID{Data: h.Sum64()}
}

func BytesID(b []byte) ID {
	h := fnv.New64a() // alloc *uint64
	h.Write(b)
	return ID{Data: h.Sum64()}
}

func Blend(id ...ID) ID {
	h := fnv.New64a() // alloc *uint64
	for _, id := range id {
		h.Write(Uint64ToBytes(id.Data)) // fnv is as fast as direct arithmetic
	}
	return ID{Data: h.Sum64()}
}

func Uint64ToBytes(u uint64) []byte {
	b := make([]byte, 8)
	b[0] = byte((u >> 8 * 0) & 0xff)
	b[1] = byte((u >> 8 * 1) & 0xff)
	b[2] = byte((u >> 8 * 2) & 0xff)
	b[3] = byte((u >> 8 * 3) & 0xff)
	b[4] = byte((u >> 8 * 4) & 0xff)
	b[5] = byte((u >> 8 * 5) & 0xff)
	b[6] = byte((u >> 8 * 6) & 0xff)
	b[7] = byte((u >> 8 * 7) & 0xff)
	return b
}
