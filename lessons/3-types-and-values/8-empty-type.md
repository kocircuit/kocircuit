# EMPTY VALUES

An _empty_ value represents "no value".

Empty values arise when a functional argument is not passed, resulting
in the argument carrying the empty value.

Ko provides two builtin functions for taking action depending on whether
a value is empty or not.

## WHEN: BRANCHING ON EMPTY

The builtin function `When` provides mechanism for branching based
on whether a value is empty or not.

`When` expects three arguments: `have`, `then` and `else`.

* When `have` is not empty, `When` returns the result of invoking (the variety) `then`,
  passing to its default argument the non-empty value of `have`.

* When `have` is empty, `When` returns the result of invoking (the variety) `else`.

For instance, in the following example the function `SmartGreeting` returns a different greeting
string depending on whether `lastName` was passed or not.

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

SmartGreeting(firstName, lastName) {
  first: String(firstName) // enforce that firstName is a string (and hence non-empty)
  return: When(
    have: lastName   // if lastName is not empty
    then: FormalGreeting[firstName: first]   // then call FormalGreeting passing lastName to its default argument
    else: InformalGreeting[firstName: first]   // otherwise call InformalGreeting
  )
}

FormalGreeting(firstName, lastName?) {
  return: strings.Join(
    string: ("Dear", firstName, lastName)
    delimiter: " "
  )
}

InformalGreeting(firstName) {
  return: strings.Join(
    string: ("Hi", firstName)
    delimiter: " "
  )
}
```

You can try `SmartGreeting` by running the functions:

```ko
SmartGreetAlice() {
  return: SmartGreeting(firstName: "Alice")   // returns "Hi Alice"
}

SmartGreetBob() {
  return: SmartGreeting(firstName: "Bob", lastName: "Thurston")   // returns "Dear Bob Thurston"
}
```

You can run those with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/SmartGreetAlice
ko play github.com/kocircuit/kocircuit/lessons/examples/SmartGreetBob
```

## ALL: ALL OR NOTHING

The builtin function `All` is designed to facilitate branching on
the condition that a few values are all non-empty.

`All` accepts any number of named arguments.
If all arguments are non-empty, `All` returns a structure containing all named arguments as fields.
If any one of the arguments is empty, `All` returns the empty value.

For instance, in the following example function `SmartGreeting2` adds onto the
logic of `SmartGreeting`.

* If both `middleName` and `lastName` are given, `SmartGreeting2` will concatenate
  them with a dash, `"-"`, and pass them to `SmartGreeting` as a last name.
* Otherwise, it will use the value of its argument `lastName` as the argument passed
  to `SmartGreeting` for last name.

```ko
SmartGreeting2(firstName, middleName, lastName) {
  middleLast: When(
    have: All(middle: middleName, last: lastName)
    then: dashedMiddleLast
    else: Return[lastName]
  )
  return: SmartGreeting(firstName: firstName, lastName: middleLast)
}

dashedMiddleLast(value?) {
  return: strings.Join(
    string: value.middle
    string: value.last
    delimiter: "-"
  )
}
```

You can try `SmartGreeting2` by running the functions:

```ko
SmartGreetAda() {
  return: SmartGreeting2(firstName: "Ada", middleName: "Lee")   // returns "Hi Ada"
}

SmartGreetEarl() {
  return: SmartGreeting2(firstName: "Earl", middleName: "Lee", lastName: "Chu")   // returns "Dear Earl Lee-Chu"
}
```

You can run those with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/SmartGreetAda
ko play github.com/kocircuit/kocircuit/lessons/examples/SmartGreetEarl
```
