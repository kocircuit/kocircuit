package gate

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// BindGate verifies the Go implementation before building and returning a harness object.
func BindGate(receiver reflect.Type) (Gate, error) {
	if receiver.Kind() != reflect.Ptr || receiver.Elem().Kind() != reflect.Struct {
		panic("receiver must be pointer to struct")
	}
	arg, err := BindStruct(receiver.Elem())
	if err != nil {
		return Gate{}, err
	}
	// verify signature of Play method
	play, ok := receiver.MethodByName("Play")
	if !ok {
		return Gate{}, fmt.Errorf("implementation %v missing a play method", receiver)
	}
	// expect one runtime context parameter
	if play.Type.NumIn() != 2 || play.Type.In(1) != reflect.TypeOf(&runtime.Context{}) {
		return Gate{}, fmt.Errorf("play method must have one context parameter")
	}
	if play.Type.NumOut() != 1 {
		return Gate{}, fmt.Errorf("play method must return a single value")
	}
	return Gate{Receiver: receiver, Struct: arg}, nil
}

// Gate represents a Go struct type with a Play method and Ko function bindings.
type Gate struct {
	Receiver reflect.Type `ko:"name=receiver"`
	Struct   Struct       `ko:"name=arg"`
}

func (f Gate) GoPkgPath() string { return f.Struct.PkgPath() }

func (f Gate) GoName() string { return f.Struct.Name() }

func (f Gate) Arg() []Field { return f.Struct.Field() }

func (f Gate) Returns() reflect.Type {
	play, ok := f.Receiver.MethodByName("Play")
	if !ok {
		panic("no play")
	}
	return play.Type.Out(0)
}

// BindStruct verifies that the given Go struct has a correctly-implemented Ko interface (i.e. bindings).
func BindStruct(t reflect.Type) (Struct, error) {
	if t.Kind() != reflect.Struct {
		return Struct{}, fmt.Errorf("implementation %s must be a go struct type", t.Name())
	}
	// verify all fields are public with ko tags
	koNames := map[string]bool{}
	monadic := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		koName, ok := Field{f}.getFieldKoName()
		if !ok {
			return Struct{}, fmt.Errorf("field %s missing ko name tag", f.Name)
		}
		if koNames[koName] {
			return Struct{}, fmt.Errorf("duplicate field %s", f.Name)
		}
		koNames[koName] = true
		if (Field{f}).IsMonadic() {
			if monadic {
				return Struct{}, fmt.Errorf("multiple monadic arguments")
			} else {
				monadic = true
			}
		}
	}
	return Struct{t}, nil
}

func MonadicStructField(s reflect.Type) (reflect.StructField, bool) {
	for i := 0; i < s.NumField(); i++ {
		f := Field{s.Field(i)}
		if f.IsMonadic() {
			return f.StructField, true
		}
	}
	return reflect.StructField{}, false
}

// Struct represents a Go struct type with Ko bindings.
type Struct struct{ reflect.Type }

func (s Struct) GoPkgPath() string { return s.Type.PkgPath() }

func (s Struct) GoName() string { return s.Type.Name() }

func (s Struct) Field() []Field {
	r := make([]Field, 0, s.Type.NumField())
	for i := 0; i < s.Type.NumField(); i++ {
		f := Field{s.Type.Field(i)}
		if f.IsVisibleInKo() {
			r = append(r, f)
		}
	}
	return r
}

func (s Struct) FieldByKoName(name string) (Field, bool) {
	for _, f := range s.Field() {
		if f.KoName() == name {
			return f, true
		}
	}
	return Field{}, false
}

func (s Struct) String() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "type %s struct {\n", s.GoName())
	for _, a := range s.Field() {
		fmt.Fprintf(&w, "\t%v\n", a)
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

// Field represents a Go struct field with Ko bindings.
type Field struct{ reflect.StructField }

func (f Field) IsGoExported() bool {
	n := f.StructField.Name
	return len(n) > 0 && strings.ToUpper(n[:1]) == n[:1]
}

func (f Field) IsVisibleInKo() bool {
	_, ok := f.getFieldKoName()
	return ok
}

func (f Field) String() string {
	return fmt.Sprintf("%s %v `%s`", f.GoName(), f.StructField.Type, f.StructField.Tag)
}

func (f Field) Name() KoGoName {
	return KoGoName{Ko: f.KoName(), Go: f.GoName()}
}

func (f Field) GoName() string {
	return f.StructField.Name
}

func (f Field) KoName() string {
	if name, ok := f.getFieldKoName(); ok {
		return name
	} else {
		return f.GoName()
	}
}

func (f Field) IsOptional() bool {
	switch f.StructField.Type.Kind() {
	case reflect.Ptr, reflect.Slice:
		return true
	case reflect.Interface, reflect.Map:
		return true
	}
	return false
}

const Monadic = "monadic"

func (f Field) IsMonadic() bool {
	return IsStructFieldMonadic(f.StructField)
}

func IsStructFieldMonadic(sf reflect.StructField) bool {
	tag := sf.Tag.Get("ko")
	for _, kv := range strings.Split(tag, ",") {
		x := strings.SplitN(kv, "=", 2)
		if len(x) > 0 && x[0] == Monadic {
			return true
		}
	}
	return false
}

func (f Field) getFieldKoName() (string, bool) {
	return StructFieldKoProtoGoName(f.StructField)
}

func StripFields(s reflect.Type) Fields {
	r := make(Fields, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		r[i] = &Field{s.Field(i)}
	}
	return r
}

type Fields []*Field

func (fields Fields) Monadic() (int, bool) {
	for i, f := range fields {
		if f.IsMonadic() {
			return i, true
		}
	}
	return -1, false
}

func (fields Fields) FieldByKoName(name string) (int, bool) {
	for i, f := range fields {
		if f.KoName() == name {
			return i, true
		}
	}
	return -1, false
}
