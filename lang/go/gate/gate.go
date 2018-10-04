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
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

var (
	goErrorType = reflect.TypeOf((*error)(nil)).Elem()
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
	switch play.Type.NumOut() {
	case 1:
		// OK
	case 2:
		// Accept (result, error)
		if play.Type.Out(1) != goErrorType {
			return Gate{}, fmt.Errorf("second return value of a play method must be of type error")
		}
	default:
		return Gate{}, fmt.Errorf("play method must return a single value with an optional error")
	}
	return Gate{Receiver: receiver, Struct: arg}, nil
}

// Gate represents a Go struct type with a Play method and Ko function bindings.
type Gate struct {
	Receiver reflect.Type `ko:"name=receiver"`
	Struct   GateStruct   `ko:"name=arg"`
}

func (f Gate) GoPkgPath() string { return f.Struct.PkgPath() }

func (f Gate) GoName() string { return f.Struct.Name() }

func (f Gate) Arg() GateFields { return f.Struct.Field() }

func (f Gate) Returns() reflect.Type {
	play, ok := f.Receiver.MethodByName("Play")
	if !ok {
		panic("no play")
	}
	return play.Type.Out(0)
}

// BindStruct verifies that the given Go struct has a correctly-implemented Ko interface (i.e. bindings).
func BindStruct(t reflect.Type) (GateStruct, error) {
	if t.Kind() != reflect.Struct {
		return GateStruct{}, fmt.Errorf("implementation %s must be a go struct type", t.Name())
	}
	// verify all fields are public with ko tags
	koNames := map[string]bool{}
	monadic := false
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		koName, ok := GateField{f}.getFieldKoName()
		if !ok {
			return GateStruct{}, fmt.Errorf("field %s missing ko name tag", f.Name)
		}
		if koNames[koName] {
			return GateStruct{}, fmt.Errorf("duplicate field %s", f.Name)
		}
		koNames[koName] = true
		if (GateField{f}).IsMonadic() {
			if monadic {
				return GateStruct{}, fmt.Errorf("multiple monadic arguments")
			} else {
				monadic = true
			}
		}
	}
	return GateStruct{t}, nil
}

func MonadicStructField(s reflect.Type) (reflect.StructField, bool) {
	for i := 0; i < s.NumField(); i++ {
		f := GateField{s.Field(i)}
		if f.IsMonadic() {
			return f.StructField, true
		}
	}
	return reflect.StructField{}, false
}

// GateStruct represents a Go struct type with Ko bindings.
type GateStruct struct{ reflect.Type }

func (s GateStruct) GoPkgPath() string { return s.Type.PkgPath() }

func (s GateStruct) GoName() string { return s.Type.Name() }

func (s GateStruct) Field() GateFields {
	r := make(GateFields, 0, s.Type.NumField())
	for i := 0; i < s.Type.NumField(); i++ {
		f := GateField{s.Type.Field(i)}
		if f.IsVisibleInKo() {
			r = append(r, f)
		}
	}
	return r
}

func (s GateStruct) FieldByKoName(name string) (GateField, bool) {
	for _, f := range s.Field() {
		if f.KoName() == name {
			return f, true
		}
	}
	return GateField{}, false
}

func (s GateStruct) String() string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "type %s struct {\n", s.GoName())
	for _, a := range s.Field() {
		fmt.Fprintf(&w, "\t%v\n", a)
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

// GateField represents a Go struct field with Ko bindings.
type GateField struct{ reflect.StructField }

func (f GateField) IsGoExported() bool {
	n := f.StructField.Name
	return len(n) > 0 && strings.ToUpper(n[:1]) == n[:1]
}

func (f GateField) IsVisibleInKo() bool {
	_, ok := f.getFieldKoName()
	return ok
}

func (f GateField) String() string {
	return fmt.Sprintf("%s %v `%s`", f.GoName(), f.StructField.Type, f.StructField.Tag)
}

func (f GateField) Name() KoGoName {
	return KoGoName{Ko: f.KoName(), Go: f.GoName()}
}

func (f GateField) GoName() string {
	return f.StructField.Name
}

func (f GateField) KoName() string {
	if name, ok := f.getFieldKoName(); ok {
		return name
	} else {
		return f.GoName()
	}
}

func (f GateField) IsOptional() bool {
	return StructFieldIsOptional(f.StructField)
}

func StructFieldIsOptional(toField reflect.StructField) bool {
	switch toField.Type.Kind() {
	case reflect.Ptr, reflect.Slice: // to field is optional
		return true
	default:
		switch {
		case StructFieldIsProtoOptOrRep(toField):
			return true
		case StructFieldWithNoKoOrProtoName(toField):
			return true
		default:
			return false
		}
	}
}

const Monadic = "monadic"

func (f GateField) IsMonadic() bool {
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

func (f GateField) getFieldKoName() (string, bool) {
	return StructFieldKoProtoGoName(f.StructField)
}

func StripFields(s reflect.Type) GateFields {
	r := make(GateFields, s.NumField())
	for i := 0; i < s.NumField(); i++ {
		r[i] = GateField{s.Field(i)}
	}
	return r
}

type GateFields []GateField

func (fields GateFields) Monadic() (int, bool) {
	for i, f := range fields {
		if f.IsMonadic() {
			return i, true
		}
	}
	return -1, false
}

func (fields GateFields) FieldByKoName(name string) (int, bool) {
	for i, f := range fields {
		if f.KoName() == name {
			return i, true
		}
	}
	return -1, false
}
