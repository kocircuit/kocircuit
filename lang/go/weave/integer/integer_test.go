package integer

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func TestWeave(t *testing.T) {
	integerTests.Run(t)
}

var integerTests = WeaveTests{
	{
		Enabled: true, // test needs package integer, move to sys
		Name:    "IntegerLen",
		File: `
		import "integer"
		Main(x) {
			array: (1, 2, 3)
			return: integer.Less(
				Len(x)
				Len(array)
			)
		}
		`,
		Arg: NewGoStruct(
			&GoField{Name: "Ko_X", Type: NewGoArray(2, GoString), Tag: KoTags("x", false)},
		),
		Result: nil,
	},
}
