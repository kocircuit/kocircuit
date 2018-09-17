# FUNCTION DEFINITION

Ko source files have simple syntax. They are comprised of package imports
and function definitions.

Every function in Ko returns exactly one value.

This choice ensures fluidity in composing functions with each other,
as we will see over time.

(Error conditions can be communicated either as fields within returned values
or via a panic/recover mechanism analogous to Go, described later.)

Let's start with a simple function, call it `Greeting`, which
takes no arguments and returns a greeting message string.
Place this function in `helloworld.ko`:

```ko
Greeting() {
  return: "Hello, there."
}
```

This function demonstrates the required elements of a function definition.
A function name (any valid identifier),
followed by a list of arguments between round brackets (none in this case),
followed by a function body between curly brackets.

Every function body has a mandatory return statement, which takes the form

```ko
return: <formula>
```

In this example, `<formula>` is the string literal "Hello, there.", which
represents the string value "Hello, there." (Formulas are discussed later.)

You can execute the `Greeting` function using:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/Greeting
```

You should see a printout:

```text
"Hello, there."
```

The Ko interpreter always prints the returned value of the interpreted function to standard output.
You will notice that the printout is quoted:
The interpreter prints values using a canonical syntax for printing values.
