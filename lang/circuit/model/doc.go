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

package model

import (
	"bytes"
	"strings"
)

func (f *Func) DocShort() string {
	var w bytes.Buffer
	w.WriteString(f.Name)
	w.WriteString("(")
	args := f.Args()
	args2 := make([]string, len(args))
	for i, a := range args {
		if f.Monadic == a {
			args2[i] = a + "?"
		} else {
			args2[i] = a
		}
	}
	w.WriteString(strings.Join(args2, ", "))
	w.WriteString(")")
	return w.String()
}

func (f *Func) DocLong() string {
	var w bytes.Buffer
	w.WriteString(f.DocShort())
	w.WriteString("\n\n")
	w.WriteString(f.Doc)
	return w.String()
}
