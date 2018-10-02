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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type testGateOk struct {
	Etc []int64     `ko:"name=etc,monadic"`
	X   *float64    `ko:"name=x"`
	Y   *testGateOk `ko:"name=y"`
}

type testGateSimple struct {
	S    string `ko:"str"`
	MonS string `ko:"mstr,monadic"`
}

func (testGateOk) Play(*runtime.Context) *complex128 {
	return nil
}

func (testGateSimple) Play(*runtime.Context) error {
	return nil
}

func TestBindGate(t *testing.T) {
	if _, err := BindGate(reflect.TypeOf(&testGateOk{})); err != nil {
		t.Errorf("bind (%v)", err)
	}
}

func TestGateFieldName(t *testing.T) {
	g, err := BindGate(reflect.TypeOf(&testGateOk{}))
	if err != nil {
		t.Fatalf("bind (%v)", err)
	}
	f, found := g.Struct.FieldByKoName("etc")
	assert.True(t, found)
	assert.Equal(t, f.KoName(), "etc")
	assert.Equal(t, f.GoName(), "Etc")
	assert.True(t, f.IsMonadic())

	f, found = g.Struct.FieldByKoName("x")
	assert.True(t, found)
	assert.Equal(t, f.KoName(), "x")
	assert.Equal(t, f.GoName(), "X")
	assert.False(t, f.IsMonadic())

	g, err = BindGate(reflect.TypeOf(&testGateSimple{}))
	if err != nil {
		t.Fatalf("bind (%v)", err)
	}

	f, found = g.Struct.FieldByKoName("str")
	assert.True(t, found)
	assert.Equal(t, f.KoName(), "str")
	assert.Equal(t, f.GoName(), "S")
	assert.False(t, f.IsMonadic())

	f, found = g.Struct.FieldByKoName("mstr")
	assert.True(t, found)
	assert.Equal(t, f.KoName(), "mstr")
	assert.Equal(t, f.GoName(), "MonS")
	assert.True(t, f.IsMonadic())

	_, found = g.Struct.FieldByKoName("notFound")
	assert.False(t, found)
}
