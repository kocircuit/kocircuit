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
	"fmt"

	"github.com/kocircuit/kocircuit/lang/go/kit/util"
)

type Chamber string

func (chamber Chamber) SheathID() *ID {
	return PtrID(StringID(string(chamber)))
}

func (chamber Chamber) SheathLabel() *string {
	return util.PtrString(string(chamber))
}

func (chamber Chamber) SheathString() *string {
	return util.PtrString(fmt.Sprintf("chamber=%s", chamber))
}

func RefineChamber(span *Span, ch string) *Span {
	return span.Refine(Chamber(ch))
}

func ChamberPath(span *Span) Path {
	if span == nil {
		return nil
	} else if chamber, ok := span.Sheath.(Chamber); ok {
		return append(ChamberPath(span.Parent), string(chamber))
	} else {
		return ChamberPath(span.Parent)
	}
}
