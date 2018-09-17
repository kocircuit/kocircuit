# TYPES

This chapter covers the _types_ of values used by the Ko language.

* Arithmetic types: string, boolean, integers and floats
* Sequence type: an ordered list of values of the same type
* Structure type: a set of named fields, each associated with a value of any type
* Variety type: a functional value that can be augmented (closure) or invoked

A distinguished "empty" type represents a missing value (e.g. an argument
that was not passed, or a selection into a non-existing structure field,
or an empty sequence).

* Empty type: a missing value

Ko provides an additional "opaque" type that can be used to hold
interface values from the underlying Go runtime. This type is used
in integration scenarios described in a separate chapter.

* Opaque type: a type for holding Go values as black-boxes from Ko's point of view
