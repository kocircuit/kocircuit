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
	"strings"
)

func (span *Span) Trajectory() string {
	return span.trajectory().String()
}

func (span *Span) trajectory() Trajectory {
	if span == nil {
		return nil
	} else if span.Sheath != nil && span.Sheath.SheathString() != nil {
		return append(span.Parent.trajectory(), TrajectoryPoint{span})
	} else {
		return span.Parent.trajectory()
	}
}

type Trajectory []TrajectoryPoint

func (traj Trajectory) String() string {
	ss := make([]string, len(traj))
	for i, point := range traj {
		ss[i] = fmt.Sprintf("[%d] %v", i, point)
	}
	return strings.Join(ss, "\n")
}

type TrajectoryPoint struct {
	Span *Span `ko:"name=span"`
}

func (p TrajectoryPoint) String() string {
	return fmt.Sprintf("%s, %s", p.Span.SourceLine(), *p.Span.Sheath.SheathString())
}
