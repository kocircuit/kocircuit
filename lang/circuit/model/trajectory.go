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
