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

package eval

import (
	"context"
	"fmt"
	"reflect"
	// goruntime "runtime"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/gate"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// EvalCallMacro is a macro that calls a Gate.
type EvalCallMacro struct {
	Gate gate.Gate `ko:"name=gate"`
}

func (m *EvalCallMacro) Splay() tree.Tree { return tree.Quote{String_: m.Help()} }

func (m *EvalCallMacro) Label() string { return "call" }

func (m *EvalCallMacro) MacroID() string { return m.Help() }

func (m *EvalCallMacro) MacroSheathString() *string { return nil }

// Help returns the a short help message of the macro
func (m *EvalCallMacro) Help() string {
	if helpMethod, found := m.Gate.Receiver.MethodByName("Help"); found && helpMethod.Type.NumIn() == 1 {
		instance := reflect.New(m.Gate.Receiver).Elem()
		result := helpMethod.Func.Call([]reflect.Value{instance})
		if len(result) == 1 {
			return result[0].String()
		}
	}
	return fmt.Sprintf("%q.%s", m.Gate.GoPkgPath(), m.Gate.GoName())
}

// Doc returns the usage documentation of the macro
func (m *EvalCallMacro) Doc() string {
	if docMethod, found := m.Gate.Receiver.MethodByName("Doc"); found && docMethod.Type.NumIn() == 1 {
		instance := reflect.New(m.Gate.Receiver).Elem()
		result := docMethod.Func.Call([]reflect.Value{instance})
		if len(result) == 1 {
			return result[0].String()
		}
	}
	return fmt.Sprintf("Run: go doc %s.%s", m.Gate.GoPkgPath(), m.Gate.GoName())
}

func (call *EvalCallMacro) Invoke(span *model.Span, arg eval.Arg) (returns eval.Return, effect eval.Effect, err error) {
	ss := arg.(*symbol.StructSymbol)
	var receiver reflect.Value
	if receiver, err = symbol.Integrate(span, ss, call.Gate.Receiver); err != nil {
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
		return symbol.Deconstruct(span, result[0]), nil, nil
	}
}

func NewEvalRuntimeContext(span *model.Span) *runtime.Context {
	return &runtime.Context{
		Parent:  nil,
		Source:  span.CommentLine(),
		Context: context.Background(),
		Kill:    nil,
	}
}
