// Package truss provides a framework for code generation.
package truss

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/toolchain"
)

type TrussCtx struct {
	Go *GoToolchain `ko:"name=go"` // compiler + repo
}
