
[![Build Status](https://travis-ci.org/kocircuit/kocircuit.svg?branch=master)](https://travis-ci.org/kocircuit/kocircuit)

# Ko

Ko is a concurrent, immutable, functional language.

Ko is both generic (function arguments and return values do not declare types)
and type-safe (a new static-type inference algorithm ensures deep static type
safety everywhere).

Ko is built on top of the Go runtime, in order to benefit from Go's rich ecosystem of 
integrations with industrial technologies.

Existing Go libraries and clients can be "exposed" in Ko with little relative effort.
Protocol definitions, like Protocol Buffers or OpenAPI, can also be exposed in Ko
using simple code-generation.

## Learning Ko

Perhaps the best way to learn the language is by reading sequentially through
our [step-by-step lessons](https://github.com/kocircuit/kocircuit/tree/master/lessons).

## Design, specifications and theory

An initial formal specification of the language (its underlying computational model,
its syntax and its type system) can be found in the evolving [Ko Handbook](https://kocircuit.github.io/).

## Why use Ko?

There are four main aspects of Ko which make it an interesting proposal:

* __Language.__ Ko is generic and at the same time entirely type-safe.
Genericity means that functions do not have to declare argument and return 
value types, which makes them highly reusable. At the same time, when entire
Ko programs are compiled against external protocols, services or types,
they are fully-verified for type compliance.

* __Types.__ Ko uses an type-system which is the common denominator of
industry protocol standards, like Protocol Buffers, Thrift, OpenAPI, and so on.
These type systems are captured by structures, sequences (repeated types), primitive types,
optional types and map types. This type system is already familiar to most programmers.

* __Architecture__. The forthcoming Ko compiler compiles Ko programs to
a [high-level intermediate representation](https://github.com/kocircuit/kocircuit/blob/master/bootstrap/asm/proto/asm.proto) (IR) which can be used to code-generate an actual implementation in
any language (e.g. Go, Java, C++, etc.) with relatively little effort.
The IR produced by the Ko compiler is a collection of functions in SSA form,
with deep [type annotations](https://github.com/kocircuit/kocircuit/blob/master/bootstrap/types/proto/types.proto) everywhere.

* __Integrations__. The Ko interpreter, being built on top of Go, can
gain access to any technology available in Go by binding dynamically to it
with little effort. This includes libraries and clients written in Go,
as well as standards like Protocol Buffers and OpenAPI which
have bindings for Go.

The forthcoming Ko compiler, being a code-generation technology,
can benefit from integrations with any target language, as described in
the architecture bullet (above).
