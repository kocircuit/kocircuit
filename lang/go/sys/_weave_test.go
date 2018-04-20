package sys

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/go/model"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func TestWeave(t *testing.T) {
	sysWeaveTests.Run(t)
}

var sysWeaveTests = WeaveTests{
	{
		Enabled: true,
		Name:    "fib",
		File: `
		import "boolean"
		import "integer"
		Fib(n) {
			return: Yield(
				if: boolean.Or(
					integer.Equal(n, 0)
					integer.Equal(n, 1)
				)
				then: Return[1]
				else: fibSum[n: n]
			)()
		}
		fibSum(n) {
			return: integer.Sum(
				Fib(n: integer.Sum(n, -1))
				Fib(n: integer.Sum(n, -2))
			)
		}
		Main() {
			return: Fix(
				Fib[n: Int64(0)]
			)
		}
		`,
		Arg:    NewGoStruct(),
		Result: nil,
	},
}
