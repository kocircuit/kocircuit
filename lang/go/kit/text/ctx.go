package text

import (
	. "github.com/kocircuit/kocircuit/lang/go/kit/console"
)

const Abyss = "█" // "ᨆ", "ᨖ"
const DepthExhausted = "ᨑ"
const BreadthExhausted = "ᨓ"
const AlreadyShown = "ᨔ"

type PrintContext struct {
	Prefix  string
	Indent  string
	Width   int     // screen width
	Console Console // color highlighting
}

var (
	DefaultPrinter     = PrintContext{Indent: "\t", Width: 100, Console: ANSI}
	DefaultMonoPrinter = PrintContext{Indent: "\t", Width: 100, Console: Mono}
)
