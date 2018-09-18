# STRUCTURES

A _structure_ type holds zero or more named fields, each associated with a value.

Structure values can be constructed using formulas of the form:

```ko
(
  <field_name_1>: <formula_1>
  ...
  <field_name_n>: <formula_n>
)
```

Here `<field_name_x>` is a name identifier for the field and
`<formula_x>` is a formula for the value of the corresponding
field.

Field values need not be of the same type.

For example,

```ko
BusinessCard() {
  return: (
    name: "Paul Erdos"
    yearOfBirth: 1913
    occupations: ("teacher", "researcher", "traveller")
  )
}
```

The field `occupations` is a sequence of strings. Ko supports an alternative
syntax, called _repeated assignment_, which allows the above function to
be written equivalently as:

```ko
BusinessCard() {
  return: (
    name: "Paul Erdos"
    yearOfBirth: 1913
    occupations: "teacher"
    occupations: "researcher"
    occupations: "traveller"
  )
}
```
