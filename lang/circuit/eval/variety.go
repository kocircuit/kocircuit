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
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// Variety is a gate shape.
type Variety struct {
	Macro Macro
	Arg   Arg
}

func (v Variety) Doc() string { return tree.Sprint(v) }

func (v Variety) String() string { return tree.Sprint(v) }

func (v Variety) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to variety")
}

func (v Variety) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into variety", path)
	}
	return v, nil, nil
}

func (v Variety) Augment(span *model.Span, arg Fields) (Shape, Effect, error) {
	aug := Fields{} // copy-and-append arg
	if v.Arg != nil {
		aug = append(aug, v.Arg.(Fields)...)
	}
	aug = append(aug, arg...)
	return Variety{Macro: v.Macro, Arg: aug}, nil, nil
}

func (v Variety) Invoke(span *model.Span) (Shape, Effect, error) {
	r, eff, err := v.Macro.Invoke(span, v.Arg)
	return r.(Shape), eff, err
}
