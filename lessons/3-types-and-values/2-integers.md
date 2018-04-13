# INTEGERS

Supported integer types are:

* Signed integers: `Int8`, `Int16`, `Int32`, `Int64`
* Unsigned integers: `Uint8`, `Uint16`, `Uint32`, `Uint64`

## INTEGER CONVERSIONS

Integer values can be converted to other integral types using the builtin
conversion functions: `Int8`, `Int16`, `Int32`, `Int64`, `Uint8`, `Uint16`, `Uint32`, `Uint64`.

All integral conversion functions are called with a single unnamed integral argument.

For instance:

	ConvertToInt32(x) {
		return: Int32(x) // if the value of x is not integral, a type panic will be issued
	}

## INTEGER ARITHMETIC

Ko provides a few generic arithmetic functions for manipulating integral types conveniently.

The functions are generic in the sense that they work with all integral types:

### NEGATION

* Function `Negative` expects a single integral argument. It returns the
negation of its argument as the same integral type.

An example usage:

	Negate(x) {
		return: Negative(x)
	}

### COMPARISON

* Function `Less` expects a sequence of two integral arguments of the same integral type.
It returns `true` if the first is strictly smaller than the second argument.

An example usage:

	IsXSmallerThan3(x) {
		return: Less(x, 3)
	}

Alternatively, one could also pass a sequence value as an argument to `Less`:

	IsXSmallerThan3(x) {
		pair: (x, 3)
		return: Less(pair)
	}

### SUMMATION

* Function `Sum` expects a sequence of one or more integral arguments of the same type.
It returns the sum of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	SumXY(x, y) {
		return: Sum(x, y)
	}

The following more advanced example shows a function for summing three
arguments. This function handles also the cases when the user does not
supply some or all of the arguments:

	SumXYZ(x, y, z) {
		return: Sum(0, x, y, z)
	}

The first element `0` ensures that the sequence of integers, `(0, x, y, z)`, passed to
`Sum` is not empty. (Note that the sequence `(x, y, z)` could be empty, if
all of the arguments `x`, `y`, `z` were not passed, or they had empty values.)

### MULTIPLICATION

* Function `Prod` expects a sequence of one or more integral arguments of the same type.
It returns the product of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	ProdXY(x, y) {
		return: Prod(x, y)
	}

### DIVISION

* Function `Ratio` expects a sequence of one or more integral arguments of the same type.
It returns the ratio of the elements in the sequence using the same integral type
as the elements themselves.

An example usage:

	RatioXYZ(x, y, z) {
		return: Ratio(x, y, z) // returns x divided by y divided by z
	}

### MODULI

* Function `Moduli` expects a sequence of two integral arguments of the same integral type.
`Moduli` returns the value of the first argument modulus the value of the second one.

An example usage:

	Mod3(x) {
		return: Moduli(x, 3) // returns the value of "x modulus 3"
	}
