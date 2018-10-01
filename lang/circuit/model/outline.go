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

import "github.com/kocircuit/kocircuit/lang/go/kit/util"

type Outline string

func (outline Outline) SheathID() *ID {
	return PtrID(StringID(string(outline)))
}

func (outline Outline) SheathLabel() *string {
	return util.PtrString(string(outline))
}

func (outline Outline) SheathString() *string {
	return util.PtrString(string(outline))
}

func RefineOutline(span *Span, outline string) *Span {
	return span.Refine(Outline(outline))
}
