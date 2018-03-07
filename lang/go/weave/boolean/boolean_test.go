package boolean

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func TestWeave(t *testing.T) {
	booleanTests.Run(t)
}

var booleanTests = WeaveTests{
	{
		Enabled: true,
		Name:    "boolean",
		File: `
		Main() {
			return: Or(
				u: And(x: false, y: true)
				w: Not(false)
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: GoTrue,
	},
}
