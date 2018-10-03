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

func (vty *VarietySymbol) Invoke(span *model.Span) (eval.Shape, eval.Effect, error) {
	if r, eff, err := vty.Macro.Invoke(
		eval.RefineMacro(span, vty.Macro),
		MakeStructSymbol(vty.Arg),
	); err != nil {
		return nil, nil, err
	} else {
		return r.(Symbol), eff, nil
	}
}
