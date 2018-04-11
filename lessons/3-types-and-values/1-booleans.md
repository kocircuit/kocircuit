# BOOLEAN

Ko booleans can assume the values `true` or `false`.

A builtin function `Bool` is provided to assert that a value is boolean.

For example, the following function returns its argument unchanged,
while also asserting that it is boolean. If it is not, a panic is produced
(resulting in an error message, unless it is recovered from):

PassBool(x) {
	return: Bool(x) // Bool returns x unchanged and panics if it is not boolean
}

## BOOLEAN ARITHMETIC

Ko provides a few builtin functions for common arithmetic manipulations over booleans.

### NEGATION

The builtin function `Not` expects a single unnamed boolean argument.
It returns its boolean neagation.

Example usage:

	NotJohn(personName) {
		return: Not(Equal(personName, "John"))
	}

### CONJUNCTION

The builtin function `And` expects a single unnamed argument, which is a sequence of booleans.
It returns the conjunction of their values. If the sequence is empty, `And` returns `true`.

The following examples returns `true` if the integral argument `y` is strictly between
the integral arguments `x` and `z`:

	import "integer"

	IsBetween(x, y, z) {
		return: And(
			integer.Less(x, y)
			integer.Less(y, z)
		)
	}

### DISJUNCTION

The builtin function `Or` expects a single unnamed argument, which is a sequence of booleans.
It returns the disjunction of their values. If the sequence is empty, `Or` returns `false`.

The following example returns `true` if any two of its three arguments, `x`, `y` and `z`, are equal:

	import "integer"

	AnyIsNonZero(x, y, z) {
		return: Or(
			Equal(x, y)
			Equal(y, z)
			Equal(x, z)
		)
	}

### EXCLISIVE DISJUNCTION

The builtin function `Xor` expects a single unnamed argument, which is a sequence of booleans.
It returns the exclusive-or of their values: `true` if an odd number of booleans are `true`,
and `false` otherwise. (If the sequence is empty, `Xor` returns `false`.)

The following example function returns `true`,
if either both of its arguments are `true` or both are `false`.

	BothOrNone(x, y) {
		return: Xor(true, x, y)
	}
