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

import "github.com/kocircuit/kocircuit/lang/circuit/model"

type Boundary interface {
	Figure(*model.Span, Figure) (Shape, Effect, error)
	Enter(*model.Span, Arg) (Shape, Effect, error)
	Leave(*model.Span, Shape) (Return, Effect, error)
}

type IdentityBoundary struct{}

func (IdentityBoundary) Figure(_ *model.Span, figure Figure) (Shape, Effect, error) {
	switch u := figure.(type) {
	case Macro:
		return Variety{Macro: u}, nil, nil
	case Shape:
		return u, nil, nil
	}
	panic("unknown figure")
}

func (IdentityBoundary) Enter(_ *model.Span, arg Arg) (Shape, Effect, error) {
	if arg == nil {
		return nil, nil, nil
	}
	return arg.(Shape), nil, nil
}

func (IdentityBoundary) Leave(_ *model.Span, shape Shape) (Return, Effect, error) {
	return shape, nil, nil
}
