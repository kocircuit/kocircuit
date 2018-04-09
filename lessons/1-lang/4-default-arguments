# OMITTED ARGUMENTS ARE EMPTY VALUES

Recall the function `CustomFormalGreeting` from a previous example:

	CustomFormalGreeting(firstName, familyName) {
		return: strings.Join(
			string: ("Hello", firstName, familyName)
			delimiter: " "
		)
	}

Ko does not require that all declared arguments are passed.
(This is why the named argument calling convention is necessary.)

When an argument is not passed, it is treated as an "empty value".

Ko will allow empty values, as long as they are not used in places
(or in ways) that require a non-empty value, in which case an error will
be reported.

`CustomFormalGreeting` could be called with either (or both) arguments missing.
The missing argument will simply be removed from the list of values
passed to the `string` argument of the call to `strings.Join`.

# DEFAULT (aka MONADIC) ARGUMENTS

Suppose we frequently call `CustomFormalGreeting` using only
the `firstName`. E.g.

TypicalGreetingForJohn() {
	return: CustomFormalGreeting(firstName: "John")
}

If a function is frequently called with one of its argument alone,
it improves readability to be able to skip the argument name in the calling syntax.

To accomplish this, Ko allows the programmer to designate
one of a function's arguments as "default" (which we also call a "monadic argument").
This is accomplished by appending a question mark, `?`, after the name
of the argument.

Let's make `firstName` a default argument:

	CustomFormalGreeting2(firstName?, familyName) {
		return: strings.Join(
			string: ("Hello", firstName, familyName)
			delimiter: " "
		)
	}

We can still call the function using named calling convention, as before, for instance:

	CustomFormalGreeting2(firstName: "John")

But we can alsos call the function with a single unnamed argument value, which will
be assigned to the default argument. Thus the following call is also equivalent:

	CustomFormalGreeting2("John")
