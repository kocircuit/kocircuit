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

	"github.com/kocircuit/kocircuit/lang/circuit/compile"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	eval.RegisterEvalGateAt("ko", "Eval", &Eval{})
}

// Eval is a service that evaluates a Ko expression.
type Eval struct {
	KoExpr string `ko:"name=koExpr"`
}

// EvalResult holds the result of the Eval service.
type EvalResult struct {
	Eval     *Eval         `ko:"name=eval"`
	Error    error         `ko:"name=error"`
	Returned symbol.Symbol `ko:"name=returned"`
}

// Play the Eval service.
func (e *Eval) Play(ctx *runtime.Context) *EvalResult {
	r := &EvalResult{Eval: e}
	framedExpr := fmt.Sprintf("Cell() {\nreturn: %s\n}", e.KoExpr)
	repo, err := compile.CompileString("jail", "expr.ko", framedExpr)
	if err != nil {
		r.Error = err
		return r
	}
	ev := eval.NewEvaluator(eval.EvalFaculty(), repo)
	span := model.NewSpan()
	r.Returned, _, _, r.Error = ev.Eval(span, repo["jail"]["Cell"], symbol.MakeStructSymbol(nil))
	return r
}
