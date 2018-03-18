package eval

import (
	"context"
	"fmt"
	"reflect"
	// goruntime "runtime"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type EvalCallMacro struct {
	Gate gate.Gate `ko:"name=gate"`
}

func (m *EvalCallMacro) Splay() Tree { return Quote{m.Help()} }

func (m *EvalCallMacro) Label() string { return "call" }

func (m *EvalCallMacro) MacroID() string { return m.Help() }

func (m *EvalCallMacro) MacroSheathString() *string { return nil }

func (m *EvalCallMacro) Help() string {
	return fmt.Sprintf("%q.%s", m.Gate.GoPkgPath(), m.Gate.GoName())
}

func (call *EvalCallMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	ss := arg.(*StructSymbol)
	var receiver reflect.Value
	if receiver, err = Integrate(span, ss, call.Gate.Receiver); err != nil {
		return nil, nil, err
	} else {
		defer func() {
			if r := recover(); r != nil {
				// buf := make([]byte, 2e5) // 200K
				// stack := string(buf[:goruntime.Stack(buf, true)])
				returns, effect, err = nil, nil, span.Errorf(nil, "calling gate panic: %v", r)
			}
		}()
		ctx := NewEvalRuntimeContext(span)
		result := receiver.MethodByName("Play").Call([]reflect.Value{reflect.ValueOf(ctx)})
		// lift result to declared return value
		m, _ := call.Gate.Receiver.MethodByName("Play")
		result[0] = result[0].Convert(m.Type.Out(0))
		//
		if returned, err := Deconstruct(span, result[0]); err != nil {
			return nil, nil, span.Errorf(err, "deconstructing gate return value")
		} else {
			return returned, nil, nil
		}
	}
}

func NewEvalRuntimeContext(span *Span) *runtime.Context {
	return &runtime.Context{
		Parent:  nil,
		Source:  span.CommentLine(),
		Context: context.Background(),
		Kill:    nil,
	}
}
