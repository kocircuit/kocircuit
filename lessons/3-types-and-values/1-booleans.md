# BOOLEAN

Ko booleans can assume the values `true` or `false`.

A builtin function `Bool` is provided to assert that a value is boolean.

For example, the following function returns its argument unchanged,
while also asserting that it is boolean. If it is not, a panic is produced
(resulting in an error message, unless it is recovered from):

```ko
PassBool(x) {
  return: Bool(x) // Bool returns x unchanged and panics if it is not boolean
}
```

## BOOLEAN ARITHMETIC

Ko provides a few builtin functions for common arithmetic manipulations over booleans.

### NOT: BOOLEAN NEGATION

The builtin function `Not` expects a single unnamed boolean argument.
It returns its boolean neagation.

Example usage:

```ko
NotJohn(personName) {
  return: Not(Equal(personName, "John"))
}
```

### AND: BOOLEAN CONJUNCTION

The builtin function `And` expects a single unnamed argument, which is a sequence of booleans.
It returns the conjunction of their values. If the sequence is empty, `And` returns `true`.

The following examples returns `true` if the integral argument `y` is strictly between
the integral arguments `x` and `z`:

```ko
IsBetween(x, y, z) {
  return: And(
    Less(x, y)
    Less(y, z)
  )
}
```

### OR: BOOLEAN DISJUNCTION

The builtin function `Or` expects a single unnamed argument, which is a sequence of booleans.
It returns the disjunction of their values. If the sequence is empty, `Or` returns `false`.

Note that Ko (in contrast to many imperative languages) will always calculate all
arguments of the `Or` function, even if one is already `true`.

The following example returns `true` if any two of its three arguments, `x`, `y` and `z`, are equal:

```ko
HasEqualPair(x, y, z) {
  return: Or(
    Equal(x, y)
    Equal(y, z)
    Equal(x, z)
  )
}
```

### XOR: BOOLEAN EXCLUSIVE-OR

The builtin function `Xor` expects a single unnamed argument, which is a sequence of booleans.
It returns the exclusive-or of their values: `true` if an odd number of booleans are `true`,
and `false` otherwise. (If the sequence is empty, `Xor` returns `false`.)

The following example function returns `true`,
if either both of its arguments are `true` or both are `false`.

```ko
BothOrNone(x, y) {
  return: Xor(true, x, y)
}
```

## YIELD: BRANCHING ON A BOOLEAN VALUE

The builtin function `Yield` provides a mechanism for branching on a boolean value.

`Yield` expects three arguments: `if`, `then` and `else`.
If the boolean argument `if` is `true`, then `Yield` returns the value of `then`.
If the boolean argument `if` is `false`, then `Yield` returns the value of `else`.

For instance, the function `GreetAsRequested` below will taylor
a different greeting message based on the value of the boolean
argument `beFormal`.

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

GreetAsRequested(beFormal, firstName, lastName) {
  return: Yield(
    if: beFormal
    then: strings.Join(
      string: ("Dear", firstName, lastName)
      delimiter: " "
    )
    else: strings.Join(
      string: ("Hi", firstName)
      delimiter: " "
    )
  )
}
```

Try this example by running:

```ko
GreetAliceFormally() {
  return: GreetAsRequested(beFormal: true, firstName: "Alice")
}
```

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/GreetAliceFormally
```

### YIELD VARIETIES TO BUILD RECURSIVE FUNCTIONS

Suppose you want to compute the n-th Fibonacci number `F(n)`,
which is defined recursively: `F(0) = 1`, `F(1) = 2` and
`F(n) = F(n-1) + F(n-2)` for `n > 1`.

The following function `Fib` demonstrates an implementation 
computing the n-th Fibonacci number.

It demonstrates how to implement self-referential recursion
using `Yield` in conjunction with variety (functional) values.

The `Yield` statement in `Fib` first determines which 
case we are in: the "base" cases `n == 0` or `n == 1` or
the "recursive" case `n > 1`.

* In the "base" case, `Yield` returns a variety (functional value)
  that will return `1` when invoked.

* In the "recursive" case, `Yield` returns a variety that will
  compute and return the sum of the previous two Fibonacci
  numbers when invoked.

After `Yield` returns the chosen variety it is invoked,
which is accomplished with the invocation formula `()`
appended after the `Yield` formula.

```ko
Fib(n?) {
  return: Yield(
    if: Or(Equal(n, 0), Equal(n, 1))   // if n == 0 or n == 1,
    then: fibBase   // then return a variety that returns 1
    else: fibRecurse[n]   // otherwise return a variety that calls Fib recursively
  )()   // invoke whichever variety was returned by Yield
}

fibBase() {
  return: 1
}

fibRecurse(m?) {
  return: Sum(
    Fib(Sum(m, -1))
    Fib(Sum(m, -2))
  )
}
```

Note that if we had implemented `Fib` to invoke `fibBase` and `fibRecurse`,
instead of passing them as varieties, both of them would be invoked before
the execution of `Yield` and this would result in a non-halting exection.
The erroneous implementation is shown below:

```ko
FibNonHalting(n?) {
  return: Yield(
    if: Or(Equal(n, 0), Equal(n, 1))
    then: fibBase()
    else: fibRecurse(n)   // fibRecurse would be invoked before Yield, resulting in infinite recursion
  )
}
```

Try the above example by running:

```ko
Fib13() {
  return: Fib(13)
}
```

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/Fib13
```
