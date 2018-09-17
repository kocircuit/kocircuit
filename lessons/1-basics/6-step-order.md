# STEP EXECUTION ORDER

As we saw, function steps are a way of computing
immutable intermediate values which might be necessary
for computing the final return value of a function.

The Ko interpreter guarantees two conditions:

(a) Every step in a function definition is executed before the function returns.
This includes steps whose values are not used to compute the final return value
of the function.

(b) If the formula for computing step B refers to the value computed by step A,
then step A will be executed before step B.

We already saw an example before:

```ko
DoubleGreeting(name1, name2) {
  firstGreeting: CustomFormalGreeting(firstName: name1)
  secondGreeting: CustomFormalGreeting(firstName: name2)
  return: strings.Join(
    string: (firstGreeting, "and", secondGreeting)
    delimiter: " "
  )
}
```

In this example, the values computed by the steps `firstGreeting` and `secondGreeting`
are needed by the formula for the `return` step.

This ensures that the formulas `CustomFormalGreeting(firstName: name1)`
and `CustomFormalGreeting(firstName: name2)` are executed before
the formula for the `return` step begins exection.

There is no guarantee on the relative order of execution of the steps `firstGreeting`
and `secondGreeting`, as they have no dependence on each other. This is
a key difference from other programming languages!

## FORCING EXECUTION ORDER EXPLICITLY

In some applications, one needs to enforce an execution order between
steps that don't depend on each other's values.

To accomplish this (and for other reasons), the Ko language allows
passing arguments to function invocations, which are not declared by
the called function.

For instance, if we want to make sure that step `firstGreeting` is
computed before `secondGreeting`, we could make the formula
for `secondGreeting` depend on it "artificially":

```ko
DoubleGreeting(name1, name2) {
  firstGreeting: CustomFormalGreeting(firstName: name1)
  secondGreeting: CustomFormalGreeting(
    firstName: name2
    _after: firstGreeting
  )
  return: strings.Join(
    string: (firstGreeting, "and", secondGreeting)
    delimiter: " "
  )
}
```

The formula for `secondGreeting`, which used to be
`CustomFormalGreeting(firstName: name2)`, now takes an additional
argument called `_after` which uses the value `firstGreeting`.

This change will tell the Ko interpreter that the value of `firstGreeting`
must be computed before the formula for `secondGreeting` can be executed.
Yet it will not affect the value computed by `secondGreeting` in any way
because the function `CustomFormalGreeting` does not recognize an argument
called `_after` and it will simply ignore it.

The argument label `_after` is arbitrary. One can use any argument name
not used by the function being called.