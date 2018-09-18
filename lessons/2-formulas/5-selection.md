# SELECTION

A _selection_ formula represents a value acquired by
selecting the value of a named field in a given structure.

The syntax for selection takes the form:

```ko
<structure_formula>.<field_name>
```

Here `<structure_formula>` is a formula that evaluates to a structure value or the empty value,
and `<field_name>` is the name of the field whose value we are selecting.

Creating structures is covered in greater detail in the following chapter.
We are going to look at a quick example here:

```ko
StevesCredentials() { // Credentials returns a structure with three fields
  return: (
    firstName: "Steven"
    lastName: "Dworkin"
    age: 37
  )
}

FirstName(credentials) {
  return: credentials.firstName // selection formula: select field "firstName" from argument "credentials"
}

StevesFirstName() { // returns "Steven"
  return: FirstName(credentials: StevesCredentials())
}
```

## SELECTING MISSING FIELDS

If a formula tries to select a field that is not present in a structure,
the selection evaluates to the empty value. In other words, it 
does not produce a compiler type-error.

This design is intentional to facilitate writing fluid generic programs.

At the same time this design is safe: If a selection results in an
empty value, then downstream in the program execution the compiler
will catch an error if the empty value is used in a context which
does not allow for one.

## SELECTING INTO EMPTY VALUES

Selection into an empty value always results into an empty value.
