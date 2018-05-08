
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

## LEARNING KO

Perhaps the best way to learn the language is by reading sequentially through
our [step-by-step lessons](https://github.com/kocircuit/kocircuit/tree/master/lessons).

## DESIGN AND THEORY

An initial formal specification of the language (its underlying computational model,
its syntax and its type system) can be found in the evolving [Ko Handbook](https://kocircuit.github.io/).
