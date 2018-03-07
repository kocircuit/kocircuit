package gate

import (
	"reflect"
	"testing"

	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type testGateOk struct {
	Etc []int64     `ko:"name=etc,monadic"`
	X   *float64    `ko:"name=x"`
	Y   *testGateOk `ko:"name=y"`
}

func (testGateOk) Play(*runtime.Context) *complex128 {
	return nil
}

func TestBindGate(t *testing.T) {
	if _, err := BindGate(reflect.TypeOf(&testGateOk{})); err != nil {
		t.Errorf("bind (%v)", err)
	}
}
