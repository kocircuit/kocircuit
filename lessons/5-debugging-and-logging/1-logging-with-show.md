# LOGGING WITH SHOW

Ko provides two builtin functions, `Show` and `ShowType`, to enable
convenient integration of logging into your program flow.

`Show` pretty-prints the values passed to it, whereas
`ShowType` pretty-prints the types of the values passed to it.

## SHOWING A VALUE AND PASSING IT THROUGH

When `Show` is invoked with an unnamed argument, it pretty-prints
the argument value to the log output (standard error) and
returns the value.

For instance, the function `VerboseSum` below will return the sum of its
(necessarily integer) arguments `x` and `y`. Before summation
is invoked however, `VerboseSum` will print the values of `x` and `y` to
standard error.

```ko
VerboseSum(x, y) {
  return: Sum(Show(x), Show(y))
}
```

You can try `VerboseSum` by running the function:

```ko
VerboseSumOneAndTwo() {
  return: VerboseSum(x: 1, y: 2)
}
```

This will produce the log output:

```text
1
2
3
```

Note that the first two lines (`1`, `2`) are the result of `Show`.
The last line (`3`) is the result of ko printing the return value of the function.

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/VerboseSumOneAndTwo
```

## SHOWING CUSTOM MESSAGES

When `Show` is invoked with named arguments, it will accept any number of named arguments.
It will return a structure containing all named arguments as fields, and it will also
print that structure to standard error.

This way of calling `Show` is particularly convenient for generating
readable log messages. This is demonstrated in the following example:

```ko
VerboseSum2(x, y) {
  ignore: Show(
    message_by: "VerboseSum2"
    arg_x: x
    arg_y: y
  )
  return: Sum(x, y)
}
```

You can try `VerboseSum2` by running the function:

```ko
VerboseSum2OneAndTwo() {
  return: VerboseSum2(x: 1, y: 2)
}
```

This will produce the log output:

```text
(message_by: "VerboseSum2", arg_x: 1, arg_y: 2)
```

You can run this with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/VerboseSum2OneAndTwo
```

## SHOWING TYPES

The builtin function `ShowType` pretty-prints the deep type structure of its values.
Otherwise it behaves identically to its analog `Show` for printing values.

The following example demonstrates pretty-printing various types:

```ko
ShowTypeExample() {
  return: ShowType(
    arg0: -3.14e-11   // 64-bit floating-point
    arg1: -7   // 64-bit integer
    arg2: "abc"   // string
    arg3: true   // boolean
    arg4: ("foo", "bar")   // sequence of strings
    arg5: (   // structure
      subarg1: "def"   // string
      subarg2: ShowTypeExample   // variety (aka functional type, aka lambda)
    )
  )
}
```

Running `ShowTypeExample` will produce the following log output:

```ko
(
  arg0: Float64
  arg1: Int64
  arg2: String
  arg3: Bool
  arg4: (String)
  arg5: (subarg1: String, subarg2: Variety)
)
```

You can run this example with:

```bash
ko play github.com/kocircuit/kocircuit/lessons/examples/ShowTypeExample
```
