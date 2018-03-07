package runtime

import (
	"time"
)

var start = time.Now()

func ExecutionID() string {
	return start.Format("2006-01-02-15-04-05")
}
