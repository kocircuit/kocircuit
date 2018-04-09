# FORMULAS

As we already saw, a Ko function body consists of a sequence of _steps_,
where every step is labeled by a unique identifier and is associated with
a _formula_ which describes how to compute its value.

Every function must have a step labeled `return`, whose value is
used as the return value of the function itself.

For instance, we looked at the function `DoubleGreeting`:

	DoubleGreeting(name1, name2) {
		firstGreeting: CustomFormalGreeting(firstName: name1)
		secondGreeting: CustomFormalGreeting(firstName: name2)
		return: strings.Join(
			string: (firstGreeting, "and", secondGreeting)
			delimiter: " "
		)
	}

This function has three steps, named `firstGreeting`, `secondGreeting` and `return`,
respectively.

In general, every step definition begins on a new line (within the function body)
and conforms to the syntax rule:

	<label>: <formula>

The expression `<label>` holds the unique idenitifier of the step (and its computed value),
whereas the expression `<formula>` describes how the step value is computed.

A _formula_ can be one of the following:

1. Literal expression, for describing constants (booleans, strings, etc.)
2. Reference to a functional argument, for referring to the value of an argument
3. Reference to another step (in the same function), for refering to a prior step's value
4. Invocation of a function (or a functional value)
5. Augmentation of a function (or a functional value), which is Ko's nomenclature for closure
6. Selection of a field into a value of a structure type

The following articles cover each of these types of formulas.
