package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func GoNameFor(koName string) string {
	if koName == NoLabel {
		return "Monadic"
	}
	return fmt.Sprintf("Ko_%s", koName)
}
