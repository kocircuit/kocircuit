# LITERALS

A _literal_ is a _formula_ which describes an arithmetic or functional constant.

Ko supports four types of constants: boolean, string, integral, floating-point and functional.

## Boolean literals

Boolean constants are captured by the identifiers `true` and `false`.

## String literals

String constants are expressed using either double-quoted or back-quoted literal notation,
following the Go definition of those:

_Double-quoted_ string literals can span a single source line and they interpret
standard ASCII escape sequences (like, `\n`, `\t`, `\r`, etc.) For example:

```ko
ExampleDoubleQuoted() {
  return: "\tHello\nworld."
}
```

_Back-quoted_ string literals can span multiple lines and their contents is interpreted
verbatim (ASCII escape sequences are not interpreted). For example:

```ko
ExampleBackQuoted() {
  return: `
This is a longer text,
which spans multiple lines.`
}
```

## Integer literals

Integer constants are expressed using Go-style integers, e.g. 

```ko
ExampleInteger() {
  return: -314
}
```

Ko interprets all integer literals as a 64-bit signed integer type (Int64).
(The can be explicitly converted to other integral types.)

## Floating-point literals

Floating-point constants are expressed using Go-style floating-point numbers, e.g.

```ko
ExampleFloat() {
  return: -3.14e+11
}
```

Ko interprets all floating-point literals as 64-bit floats (Float64).
(The can be explicitly converted to 32-bit floats, Float32.)

## Functional literals

Ko supports functional values. A functional literal expression
refers to a user or builtin function and has a functional value type (which
in Ko is called a _variety_ type).

There are three types of functional literals:

1. A functional literal which refers to another function defined
   in the same package is the name identifier of this function. E.g.

```ko
F() { return: "Hello" }

ReturnF() {
  return: F // F is a functional literal refering to F (above)
}
```

2. A functional literal which refers to a function defined in another
   package is written in the form

```ko
<pkg_alias>.<func_name>
```

For instance,

```ko
import "github.com/kocircuit/kocircuit/lib/strings"

ReturnJoin() {
  return: strings.Join // strings.Join refers to a function defined in package strings
}
```

3. A functional literal which refers to a builtin function is the name identifier of the function.
  For context, Ko supplies a number of builtin functions (like, `Print`, `Yield`, `Range`, etc.)
  which can be referred to from any source file (without importing any packages).

  For instance,

```ko
ReturnYield() {
  return: Yield // returns a functional value pointing to the builtin Yield function
}
```
