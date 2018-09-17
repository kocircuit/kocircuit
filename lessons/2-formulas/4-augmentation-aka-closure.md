# AUGMENTATION

Augmentation is the Ko analog of functional closure.

An _augmentation_ formula represents a functional
value which is derived from a prior functional value
by assigning values to one or more of its arguments.

Augmentation formulas have identical syntax to invocation
formulas, but instead of using round brackets `()`, they
use square brackets `[]`.

There are two types of augmentations:
augmentation with named arguments and
augmentation with a default (aka monadic) argument.

## AUGMENTATION WITH NAMED ARGUMENTS

An augmentation formula with named arguments is expressed in the form:

```ko
<function_formula>[
  <arg_name_1>: <formula_1>
  ...
  <arg_name_n>: <formula_n>
]
```

Note the square brackets.

Let's use this in an example:

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

// MakeGreeter returns a functional value, which
// when invoked will return a greeting string, customized to the 
// given name argument.
MakeGreeter(name) {
  return: strings.Join[ // augment the function Join
    string: ("Hello", name)
    delimiter: " "
  ]
}

Greet(name) {
  greeter: MakeGreeter(name: name) // create a function that returns greetings to name
  return: greeter() // call the function from step greeter to return the actual greeting string
}

GreetTom() {
  return: Greet(name: "Tom") // this will return "Hello Tom"
}
```

This example can be found in
[github.com/kocircuit/kocircuit/lessons/examples/augment.ko](github.com/kocircuit/kocircuit/lessons/examples/augment.ko).

Run with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/GreetTom
```

## AUGMENTATION WITH A DEFAULT ARGUMENT

Functions that can be invoked with a default argument can be augmented
with a default argument. The syntax for such augmentation formulas is:

```ko
<function_formula>[<default_arg_formula>]
```

Again, note the use of square brackets , `[]`, instead of round ones, `()`, which are used for invocation.

As an example usage, we change the previous example to use default argument augmentation
of the `strings.Join` function:

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

MakeGreeter(name) {
  // augment Join by setting its default argument to the list ("Hello", name)
  return: strings.Join[("Hello", name)]
}

Greet(name) {
  greeter: MakeGreeter(name: name) // create a function that returns greetings to name
  return: greeter() // call the function from step greeter to return the actual greeting string
}

GreetTom() {
  return: Greet(name: "Tom") // this will return "Hello Tom"
}
```
