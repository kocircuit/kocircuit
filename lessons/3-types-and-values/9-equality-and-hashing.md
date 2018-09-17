# EQUALITY AND HASHING

Ko provides generic functions for equality between values (of any type),
as well as for computing a hash of any value.

## EQUALITY

The builtin `Equal` function expects a single unnamed argument of a sequence type.
It returns `true` if the values of all elements in the sequence are equal or if the
sequence is empty. Otheriwse, it returns `false`.

For example:

```ko
DoJohnAndPaulHaveTheSameOccupation() { // returns `true`
  return: Equal( // check equality of the "occupation" structures of John and Paul
    John().occupation
    Paul().occupation
  )
}

John() {
  return: (
    name: "John"
    occupation: (employer: "IBM", role: "Engineer")
}

Paul() {
  return: (
    name: "Paul"
    occupation: (employer: "IBM", role: "Engineer")
}
```

## HASHING

The builtin `Hash` function expects a single unnamed argument of any type.
It returns a short string representing a unique (up to collisions) hash
for the supplied argument value. The hashing algorithm used is based on FNV.

Two hashes are equal, if and only if the hashed values are equal.

For example:

```ko
DetectHashCollision(x, y) {
  return: And(
    Equal(x, y)
    Not(Equal(Hash(x), Hash(y)))
  )
}
```
