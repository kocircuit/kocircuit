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

package symbol

import (
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

type EvokeArg struct {
	Name  string `ko:"name=name"`
	Value Symbol `ko:"name=value"`
}

func Evoke(vty *VarietySymbol, args ...EvokeArg) Symbol {
	ff := make(eval.Fields, len(args))
	for i, arg := range args {
		ff[i] = eval.Field{Name: arg.Name, Shape: arg.Value}
	}
	if returns, _, err := vty.Evoke(model.NewSpan(), ff); err != nil {
		panic(err)
	} else {
		return returns.(Symbol)
	}
}

func (vty *VarietySymbol) Evoke(span *model.Span, fields eval.Fields) (Symbol, eval.Effect, error) {
	if augmented, _, err := vty.Augment(span, fields); err != nil {
		return nil, nil, err
	} else if returns, _, err := augmented.Invoke(span); err != nil {
		return nil, nil, err
	} else {
		return returns.(Symbol), nil, nil
	}
}
