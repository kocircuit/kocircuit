# STRINGS

Ko provides a native string type.

A builtin function `String` is provided to assert that a value is a string.

For example, the following function returns its argument unchanged,
while also asserting that it is a string. If it is not, a panic is produced
(resulting in an error message, unless it is recovered from):

PassString(x) {
	return: String(x) // String returns x unchanged and panics if it is not a string
}

## STRING MANIPULATIONS

Functions for manipulating strings are not built into the Ko compiler,
rather they are provided by the (evolving) Ko standard library package
`"github.com/kocircuit/kocircuit/lib/strings"`.
