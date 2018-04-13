# VARIETIES

Varieties are the functional value type of Ko.
A _variety_ value holds a reference to a builtin or user-defined function,
as well as a set of argument assignments. 

Variety values can be constructed, augmented or invoked.

Augmenting a variety adds argument assignments to the variety, resulting in a new variety.
Augmentation is the Ko nomenclature for "functional closure".

Invoking a variety results in invoking the function referenced by the variety
with the argument assignments stored within the variety.

## CONSTRUCT

Varieties can be constructed by using a functional literal formula. For example:

	import "github.com/kocircuit/kocircuit/lib/strings"

	Joiner() {
		return: strings.Join // returns a variety of strings.Join
	}

Here `Joiner` returns a variety referencing `string.Join`, with no argument assignments.

## CONSTRUCT AND AUGMENT

Varieties can be constructed and augmented in the same formula. For example:

	SpacedJoiner() {
		return: strings.Join[delimiter: " "]
	}

Function `SpacedJoiner` returns a variety of `strings.Join` with argument `delimiter` set to `" "`.

## INVOKE

A variety value can be invoked using the standard invocation syntax. For example:

	JoinAliceAndBob() {
		joiner: SpacedJoiner()
		return: joiner[string: ("Alice", "Bob")] // augment and invoke joiner
	}

Here `joiner` is a variety of `strings.Join`, with argument `delimiter` set to `" "`.

The `return` step first augments `joiner` by assigning the value
`("Alice", "Bob")` to argument `string`. Then it invokes the resulting
variety.

Function `JoinAliceAndBob` is equivalent to:

	JoinAliceAndBobDirectly() {
		return: strings.Join(
			delimiter: " "
			strings: ("Alice", "Bob")
		)
	}

## USING VARIETIES TO BUILD INTERFACES

Using varieties and Ko's generic function semantics, one can
recreate the semantic offered by Go's interfaces, for instance.

	import "github.com/kocircuit/kocircuit/lib/strings"
	import "github.com/kocircuit/kocircuit/lib/integer" // for FormatInt64

	// Returns returns its default argument.
	Return(pass?) {
		return: pass
	}

	// Age returns the difference between currentYear and bornYear.
	Age(bornYear, currentYear) {
		return: Sum(currentYear, Negative(bornYear))
	}

	// AliceInfo returns a structure with two "methods", Name and Age.
	AliceInfo() {
		return: (
			Name: Return["Alice"]
			Age: Age[bornYear: 1901]
		)
	}

	// BobInfo returns a structure with two "methods", Name and Age.
	BobInfo() {
		return: (
			Name: Return["Bob"]
			Age: Age[bornYear: 1930]
		)
	}

	// PersonAge returns a message of the form "NAME is AGE years old."
	// NAME and AGE are extracted from the corresponding methods of personInfo.
	PersonAge(personInfo, currentYear) {
		return: strings.Join(
			string: personInfo.Name()
			string: "is"
			string: integer.FormatInt64(
				int64: personInfo.Age(currentYear: currentYear)
				base: 10
			)
			string: "years old."
			delimiter: " "
		)
	}

	// Run with:
	// ko play github.com/kocircuit/kocircuit/lessons/examples/AliceAge
	AliceAge() {
		return: PersonAge(personInfo: AliceInfo(), currentYear: 2018)
	}

	// Run with:
	// ko play github.com/kocircuit/kocircuit/lessons/examples/BobAge
	BobAge() {
		return: PersonAge(personInfo: BobInfo(), currentYear: 2018)
	}

This example is provided in:

	github.com/kocircuit/kocircuit/lessons/examples/variety.ko
