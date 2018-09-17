# FUNCTIONAL ARGUMENTS

Let's elaborate the greeting example to return a customized greeting.

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

CustomGreeting(name) {
  return: strings.Join(
    string: ("Hello", name)
    delimiter: " "
  )
}
```

This function declares an argument `name`. It then invokes the 
function `Join` provided by the `strings` package to join
the string "Hello" and the value of the argument `name` into
a single string, delimiting them with a space, " ".

Let's write a higher level function `GreetJohn` that calls `CustomGreeting`
and specializes its argument to be the string "John":

```ko
GreetJohn() {
  return: CustomGreeting(name: "John")
}
```

This example shows how to call functions and set their arguments.
As shown, Ko uses named arguments when calling functions.

You can try running the function `GreetJohn` with:

```ko
ko play github.com/kocircuit/kocircuit/lessons/examples/GreetJohn
```

To get the output:

```text
"Hello John"
```

Now that we've established the named argument calling convention,
you can guess from the implementation of `CustomGreeting` that
the function `strings.Join` expects two arguments, one named 
`string`, the other named `delimiter`.

Naturally, functions can have more than one argument. For instance,
let's extend the example to distinguish a first and a family name:

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

CustomFormalGreeting(firstName, familyName) {
  return: strings.Join(
    string: ("Hello", firstName, familyName)
    delimiter: " "
  )
}

GreetJohnFormally() {
  return: CustomFormalGreeting(firstName: "John", familyName: "Jovi")
}
```

Ko treats commas and new lines identically. For clarity, we could re-format
the last function as:

```ko
GreetJohnFormally() {
  return: CustomFormalGreeting(
    firstName: "John"
    familyName: "Jovi"
  )
}
```
