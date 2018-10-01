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
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Macro interface {
	Figure
	MacroID() string // MacroID uniquely identifies the meaning of this macro
	// The top-level circuit step in the frame argument is the circuit invocation step
	// corresponding to this operator invocation.
	Invoke(*model.Span, Arg) (Return, Effect, error)
	Label() string              // string identifier
	Help() string               // human-readable line of text, used in "ko list"
	MacroSheathString() *string // human-readable line ot text show in debug frames (if not nil)
	Doc() string
}

func RefineMacro(span *model.Span, macro Macro) *model.Span {
	return span.Refine(MacroSheath{macro})
}

type MacroSheath struct {
	Macro Macro `ko:"name=macro"`
}

func (sh MacroSheath) SheathID() *model.ID {
	return model.PtrID(model.StringID(sh.Macro.MacroID()))
}

func (sh MacroSheath) SheathLabel() *string {
	return util.PtrString(sh.Macro.Label())
}

func (sh MacroSheath) SheathString() *string {
	return sh.Macro.MacroSheathString()
}
