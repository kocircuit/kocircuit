package time

import (
	"fmt"
	"time"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoSleep))
	RegisterEvalGate(new(GoSecond))
	RegisterEvalGate(new(GoMinute))
	RegisterEvalGate(new(GoHour))
	RegisterEvalGate(new(GoFormatDurationSeconds))
}

type GoSleep struct {
	Duration time.Duration `ko:"name=duration,monadic"`
}

func (g GoSleep) Play(ctx *runtime.Context) time.Duration {
	time.Sleep(g.Duration)
	return g.Duration
}

type GoSecond struct {
	Scale *int64 `ko:"name=scale,monadic"`
}

func (g GoSecond) Play(ctx *runtime.Context) time.Duration {
	return time.Duration(OptInt64(g.Scale, 1) * int64(time.Second))
}

type GoMinute struct {
	Scale *int64 `ko:"name=scale,monadic"`
}

func (g GoMinute) Play(ctx *runtime.Context) time.Duration {
	return time.Duration(OptInt64(g.Scale, 1) * int64(time.Minute))
}

type GoHour struct {
	Scale *int64 `ko:"name=scale,monadic"`
}

func (g GoHour) Play(ctx *runtime.Context) time.Duration {
	return time.Duration(OptInt64(g.Scale, 1) * int64(time.Hour))
}

type GoFormatDurationSeconds struct {
	Duration int64 `ko:"name=duration,monadic"` // nanoseconds
}

func (g GoFormatDurationSeconds) Play(ctx *runtime.Context) string {
	return fmt.Sprintf("%1.3fs\n", float64(g.Duration)/1e9)
}
