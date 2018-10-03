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

func (ss *StructSymbol) GetMonadic() (Symbol, bool) {
	for _, field := range ss.Field {
		if field.Monadic {
			return field.Value, true
		}
	}
	return nil, false
}

func (ss *StructSymbol) Link(span *model.Span, name string, monadic bool) (eval.Shape, eval.Effect, error) {
	return ss.LinkField(name, monadic), nil, nil
}

func (ss *StructSymbol) LinkField(name string, monadic bool) Symbol {
	if found := ss.FindName(name); found != nil {
		return found.Value
	} else if monadic {
		if found := ss.FindMonadic(); found != nil {
			return found.Value
		}
	}
	return EmptySymbol{}
}

func (ss *StructSymbol) Select(span *model.Span, path model.Path) (_ eval.Shape, _ eval.Effect, err error) {
	if len(path) == 0 {
		return ss, nil, nil
	} else {
		return ss.Walk(path[0]).Select(span, path[1:])
	}
}

func (ss *StructSymbol) Walk(step string) Symbol {
	if found := ss.FindName(step); found != nil {
		return found.Value
	} else {
		return EmptySymbol{}
	}
}
