# PANIC AND RECOVER

Ko provides a mechanism for throwing runtime panics (aka exceptions),
which is analogous to its counterparts in languages like Go, Java, Python, etc.

## PANIC: SIGNALING RUNTIME ERRORS

The builtin function `Panic` throws a panic up the stack.
If the panic is not handled (using `Recover`, covered below),
it results in program termination with a message to the log output.
The panic message includes the stack of the location that caused
the panic, as well as a user-supplied panic value.

`Panic` accepts any number of named arguments (with any values).
`Panic` creates a structure whose fields are the named arguments 
passed to it, and attaches it as a _panic value_ to the panic itself.
(Panic values can be retrieved using `Recover`, as explained in the next section.)

`Panic` never returns into the function calling it.

For example, the function `AgeDifference` below computes and returns the
difference between the ages of a father and his child, passed as arguments.

`AgeDifference` also checks whether the age of the child is smaller than
that of the father. If it is not, it will panic with the structure value
`(ageDifferenceError: "a father cannot be older than his child")`.

```ko
AgeDifference(childAge, fatherAge) {
  check: Yield(
    if: Not(Less(childAge, fatherAge))
    then: Panic[ageDifferenceError: "a father cannot be older than his child"]
    else: []
  ) () // invoke whichever variety was yielded
  return: Sum(fatherAge, Negative(childAge))
}
```

You can try calling `AgeDifference` with invalid arguments, for instance,
by running the function:

```ko
InvalidAgeDifference() {
  return: AgeDifference(childAge: 21, fatherAge: 19)
}
```

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/InvalidAgeDifference
```

## RECOVER: HANDLING PANICS

The builtin function `Recover` provides a mechanism for handling panics,
caused by runtime conditions.

`Recover` expects two arguments, `invoke` and `panic`, both of which must be varieties (functional values).

`Recover` starts by invoking the functional value `invoke` (without passing any arguments):

* If the invocation succeeds in returning a value without panicking, `Recover` will return that value.

* If the invocation panics, `Recover` will capture the panic and invoke
  the functional value of the `panic` argument, while also passing the panic value
  as a default (aka monadic) argument to `panic`. Whatever the call to `panic` returns,
  will be returned by `Recover`.

For example, function `RecoverInvalidAgeDifference` below upgrades our
previous example `InvalidAgeDifference` to handle the panic and return
a message to the user.

```ko
RecoverInvalidAgeDifference() {
  return: Recover(
    invoke: AgeDifference[childAge: 21, fatherAge: 19]
    panic: recoverFromAgeDifferencePanic
  )
}

recoverFromAgeDifferencePanic(panicValue?) {
  return: (age_difference_failed_with_message: panicValue.ageDifferenceError)
}
```

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/RecoverInvalidAgeDifference
```
