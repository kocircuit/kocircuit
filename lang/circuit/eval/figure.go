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

type Figure interface{}

type Empty struct{}

func (e Empty) String() string { return tree.Sprint(e) }

func (e Empty) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to empty")
}

func (e Empty) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	return e, nil, nil
}

func (e Empty) Augment(span *model.Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting an empty")
}

func (e Empty) Invoke(span *model.Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking an empty")
}

type Integer struct{ Value_ int64 }

func (v Integer) String() string { return tree.Sprint(v) }

func (v Integer) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to integer")
}

func (v Integer) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into integer", path)
	}
	return v, nil, nil
}

func (v Integer) Augment(span *model.Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting an integer")
}

func (v Integer) Invoke(span *model.Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking an integer")
}

type Float struct{ Value_ float64 }

func (v Float) String() string { return tree.Sprint(v) }

func (v Float) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to float")
}

func (v Float) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into float", path)
	}
	return v, nil, nil
}

func (v Float) Augment(span *model.Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a float")
}

func (v Float) Invoke(span *model.Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a float")
}

type Bool struct{ Value_ bool }

func (v Bool) String() string { return tree.Sprint(v) }

func (v Bool) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to bool")
}

func (v Bool) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into bool", path)
	}
	return v, nil, nil
}

func (v Bool) Augment(span *model.Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a bool")
}

func (v Bool) Invoke(span *model.Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a bool")
}

type String struct{ Value_ string }

func (v String) String() string { return tree.Sprint(v) }

func (v String) Link(span *model.Span, name string, monadic bool) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "linking argument to string")
}

func (v String) Select(span *model.Span, path model.Path) (Shape, Effect, error) {
	if len(path) > 0 {
		return nil, nil, span.Errorf(nil, "selecting %v into string", path)
	}
	return v, nil, nil
}

func (v String) Augment(span *model.Span, _ Fields) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting a string")
}

func (v String) Invoke(span *model.Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking a string")
}
