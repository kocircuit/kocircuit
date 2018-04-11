# INTEGERS

Supported integer types are:

* Signed integers: `Int8`, `Int16`, `Int32`, `Int64`
* Unsigned integers: `Uint8`, `Uint16`, `Uint32`, `Uint64`

## INTEGER CONVERSIONS

Integer values can be converted to other integral types using the builtin
conversion functions: `Int8`, `Int16`, `Int32`, `Int64`, `Uint8`, `Uint16`, `Uint32`, `Uint64`.

All integral conversion functions are called with a formula of the form:

	<conversion_function_name>(<formula>)

For instance:

	ConvertToInt32(x) {
		return: Int32(x) // if the value of x is not integral, a compiler error will be issued
	}

## INTEGER ARITHMETIC

Ko includes a builtin package called "integer" which provides a few generic
arithmetic functions for manipulating integral types conveniently.

The functions are generic in the sense that they work with all integral types:

### NEGATION

* Function `Negative` expects a single integral argument. It returns the
negation of its argument as the same integral type.

An example usage:

	import "integer"

	Negate(x) {
		return: integer.Negative(x)
	}

### COMPARISON

* Function `Less` expects a sequence of two integral arguments of the same integral type.
It returns `true` if the first is strictly smaller than the second argument.

An example usage:

	import "integer"

	IsXSmallerThan3(x) {
		return: integer.Less(x, 3)
	}

Alternatively, one could also pass a sequence value as an argument to `Less`:

	IsXSmallerThan3(x) {
		pair: (x, 3)
		return: integer.Less(pair)
	}

### SUMMATION

* Function `Sum` expects a sequence of one or more integral arguments of the same type.
It returns the sum of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	import "integer"
	
	SumXY(x, y) {
		return: integer.Sum(x, y)
	}

The following more advanced example shows a function for summing three
arguments. This function handles also the cases when the user does not
supply some or all of the arguments:

	SumXYZ(x, y, z) {
		return: integer.Sum(0, x, y, z)
	}

The first element `0` ensures that the sequence of integers, `(0, x, y, z)`, passed to
`Sum` is not empty. (Note that the sequence `(x, y, z)` could be empty, if
all of the arguments `x`, `y`, `z` were not passed, or they had empty values.)

### MULTIPLICATION

* Function `Prod` expects a sequence of one or more integral arguments of the same type.
It returns the product of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	import "integer"
	
	ProdXY(x, y) {
		return: integer.Prod(x, y)
	}

### DIVISION

* Function `Ratio` expects a sequence of one or more integral arguments of the same type.
It returns the ratio of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	import "integer"
	
	RatioXYZ(x, y, z) {
		return: integer.Ratio(x, y, z) // returns x divided by y divided by z
	}

### MODULI

* Function `Moduli` expects a sequence of two integral arguments of the same integral type.
`Moduli` returns the value of the first argument modulus the value of the second one.

An example usage:

	import "integer"

	Mod3(x) {
		return: integer.Moduli(x, 3) // returns the value of "x modulus 3"
	}
