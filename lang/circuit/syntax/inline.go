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

package syntax

type Inline struct {
	Design []Design `ko:"name=design"` // inline function definitions
	Series []Term   `ko:"name=series"` // inline step definitions, arising from series composition
}

func (inline Inline) Union(u Inline) Inline {
	return Inline{
		Design: append(append([]Design{}, inline.Design...), u.Design...),
		Series: append(append([]Term{}, inline.Series...), u.Series...),
	}
}
