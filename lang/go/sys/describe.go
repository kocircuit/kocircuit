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
	"bytes"
	"fmt"

	"github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	eval.RegisterEvalGateAt("", "Print", new(GoPrint))
	eval.RegisterEvalGateAt("", "Println", new(GoPrintln))
}

type GoPrint struct {
	String []string `ko:"name=string,monadic"`
}

func (p *GoPrint) Play(ctx *runtime.Context) string {
	switch len(p.String) {
	case 0:
		fmt.Print()
		return ""
	case 1:
		fmt.Print(p.String[0])
		return ""
	default:
		var w bytes.Buffer
		for _, s := range p.String {
			w.WriteString(s)
		}
		fmt.Print(w.String())
		return w.String()
	}
}

type GoPrintln struct {
	String []string `ko:"name=string,monadic"`
}

func (p *GoPrintln) Play(ctx *runtime.Context) string {
	q := &GoPrint{
		String: append(append([]string{}, p.String...), "\n"),
	}
	return q.Play(ctx)
}
