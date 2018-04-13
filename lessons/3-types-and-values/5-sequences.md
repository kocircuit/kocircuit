# SEQUENCES

A _sequence_ type holds an ordered sequence of values of the same type.

Sequences can be constructed using a formula with the following syntax:

	(
		<formula_1>
		...
		<formula_n>
	)

Here `<formula_x>` represents a formula expressing the value of the corresponding 
sequence element.

The compiler will ensure that all element formulas can be captured by the same common type.

For instance, the following function returns a sequence of two strings containing
`"Hello"` and the value of argument `name`:

	AppendNameToGreeting(name) {
		return: (
			"Hello"
			name
		)
	}

Since in Ko new lines and commas are syntactically equivalent, the previous 
function can also be written as:

	AppendNameToGreeting(name) {
		return: ("Hello", name)
	}

Sequence construction ignores empty values. Hence if the caller of 
`AppendNameToGreeting` does not pass a value for the `name` argument,
the function will return the singletong sequence `("Hello")`.

If we want to ensure that name is passed and it is a string, we could write instead:

	AppendRequiredNameToGreeting(name) {
		return: ("Hello", String(name)) // String ensures that name is a string
	}

## PASSING SEQUENCES AS ARGUMENTS

Sequence constructions (as the ones above) can be used as function
arguments on invocation, as in this example:

	import "github.com/kocircuit/kocircuit/lib/strings"

	MergeTwoStrings(x, y) {
		return: strings.Join(
			string: (x, y) // sequence (x, y) passed to argument string of strings.Join
			delimiter: ""
		)
	}

Ko supports an alternative syntax, called _repeated assignment_,
for constructing sequences that are passed
to a function argument. This syntax is inspired by an analog in the textual
syntax for Protocol Buffer values.

The function `MergeTwoStrings` could be equivalently written as:

	MergeTwoStrings(x, y) {
		return: strings.Join(
			string: x // passing the argumet "string" multiple times constructs a sequence
			string: y
			delimiter: ""
		)
	}

## SEQUENCE LENGTH

One can obtain the length of a sequence value in the form of a 64-bit integer, `Int64`,
using the builtin `Len` function. `Len` expects a single unnamed sequence argument.
For instance,

	Len("a", "b", "c") // returns 3

Ko treats the empty value as equivalent to an empty sequence, thus `Len()` would evaluate to `0`.
The following function will return the number of non-empty arguments passed to it:

	CountNonEmptyArgs(x, y, z) {
		return: Len(x, y, z)
	}

Ko also treats non-empty, non-sequence values (e.g. primitive values, structures, varieties, etc.)
as equivalent to a sequence with a singleton element holding that value.
Thus `Len("a")` and `Len(5)` both evaluate to `1`.

## MERGING SEQUENCES

Two or more sequences holding values of compatible (aka unifiable) types can be merged into
a single sequence holding the concatenation of the original sequences. This is accomplished
with the builtin function `Merge`.

	Merge(("a", "b"), ("x", "y", "z")) // returns ("a", "b", "x", "y", "z")
	Merge("a", ("x", "y", "z")) // returns ("a", "x", "y", "z")

Note that singleton values, like `"a"` above, are lifted to a sequence, in this case `("a")`,
before merging.

## TAKING AN ELEMENT FROM THE FRONT

The builtin function `Take` will split an input sequence into its first element and the remainder.

`Take` expects a single unnamed sequence argument.
It returns a structure with two fields `first` and `remainder`.
Field `first` holds the value of the first element in the input sequence;
if the sequence is empty, it holds the empty value.
Field `remainder` holds the remainder of the input sequence, after the first
element.

For instance,

	Take("a", "b", "c") // returns (first: "a", remainder: ("b", "c"))
	Take("a") // returns (first: "a", remainder: ())
	Take() // returns (first: (), remainder: ())

## RANGING OVER A SEQUENCE

Ko provides the builtin `Range` function for ranging over the elements of a sequence.

`Range` iterates sequentially over the elements of an input sequence.

A user-supplied iterator function is invoked for each sequence element.
The iterator function is expected to return a structure with two fields: `emit` and `carry`.

Values stored in `emit` are collected across all invocations of the user
function and merged into a new sequence and returned by `Range`.

The value stored in `carry` is passed to the next invocation of the iterator function.
The `carry` returned by the last invocation of the iterator is returned by `Range`.

`Range` expects three input arguments:
* `over` holds the sequence value to be ranged over,
* `with` holds the user-supplied iterator function,
* `start` holds the carry value to be passed to the first invocation of the iterator.

`Range` passes two arguments to the user-supplied iterator:
* `elem` holds the element currently being iterated over
* `carry` holds the carry value from the previous invocation of the iterator,
or in the case of the first iteration, it holds the value of `start`

`Range` returns a structure with two fields: `image` and `residue`.
* `image` holds a merged sequence of all values emitted by the per-element
invocations to the iterator,
* `residue` holds the value of the carry returned by the last iterator invocation;
if the input sequence is empty, `residue` holds the value of `start`.

### EXAMPLE: WORD LENGTH AND AGGREGATE LENGTH

The following example function traverses a sequence of strings.

It returns a new sequence of integers representing the length of each string in the input sequence,
as well as the total length of all strings.

	import "github.com/kocircuit/kocircuit/lib/strings"

	// Run with:
	// ko play github.com/kocircuit/kocircuit/lessons/examples/DemoCountWordLengths
	// You should get the output:
	// (stringLengthsSeqeunce: (1, 1, 3, 3), totalLength: 8)
	DemoCountWordLengths() {
		return: CountWordLengths(stringSeq: ("a", "b", "foo", "bar"))
	}

	CountWordLengths(stringSeq) {
		ranged: Range(
			over: stringSeq   // sequence of strings to range over
			with: stringLenIterator   // iterator function
			start: (totalLength: 0)   // initial carry
		)
		return: (
			stringLengthsSeqeunce: ranged.image
			totalLength: ranged.residue.totalLength
		)
	}

	stringLenIterator(carry, elem) {
		elemLen: strings.Len(elem)
		return: (
			emit: elemLen
			carry: (
				totalLength: Sum(carry.totalLength, elemLen)
			)
		)
	}

This example can be found in:

	github.com/kocircuit/kocircuit/lessons/examples/count.ko

More advanced usage of `Range` can be found in the standard library, e.g.
function `Map` in package `"github.com/kocircuit/kocircuit/lib/series"`.
