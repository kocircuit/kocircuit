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

package sys

import (
	"fmt"
	"path"

	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

// Play is a service that plays a function in a Ko repository.
type Play struct {
	Pkg     string               `ko:"name=pkg"`  // e.g. github.com/kocircuit/kocircuit/codelab
	Func    string               `ko:"name=func"` // e.g. HelloWorld
	Repo    model.Repo           `ko:"name=repo"` // compiled ko repo
	Faculty eval.Faculty         `ko:"name=faculty"`
	Arg     *symbol.StructSymbol `ko:"name=arg"` // arg can be nil
}

// Play the function
func (w *Play) Play(ctx *runtime.Context) *PlayResult {
	pfe := &PlayFuncEval{
		Func: w.Repo[w.Pkg][w.Func],
		Eval: go_eval.NewEvaluator(w.Faculty, w.Repo),
		Arg:  w.Arg,
	}
	if pfe.Func == nil {
		return &PlayResult{
			Play:  w,
			Error: fmt.Errorf("cannot find main circuit %s", path.Join(w.Pkg, w.Func)),
		}
	}
	return pfe.Play(ctx)
}

type PlayFuncEval struct {
	Func *model.Func          `ko:"name=func"`
	Eval *go_eval.Evaluate    `ko:"name=eval"`
	Arg  *symbol.StructSymbol `ko:"name=arg"` // arg can be nil
}

func (w *PlayFuncEval) Play(ctx *runtime.Context) *PlayResult {
	r := &PlayResult{PlayFuncEval: w}
	span := model.NewSpan()
	var arg symbol.Symbol
	if w.Arg == nil {
		arg = symbol.MakeStructSymbol(nil)
	} else {
		arg = w.Arg
	}
	r.Returned, _, _, r.Error = w.Eval.Eval(span, w.Func, arg)
	return r
}

type PlayResult struct {
	Play         *Play         `ko:"name=play"`
	PlayFuncEval *PlayFuncEval `ko:"name=playFuncEval"`
	Returned     symbol.Symbol `ko:"name=returned"`
	Error        error         `ko:"name=error"`
}
