# LOOPING

The builtin function `Loop` is a mechanism for
running a user function in a loop, providing ways
to carry values from one invocation to the next
and an optional stopping condition.

`Loop` expects three arguments named: `start`, `step` and `stop`.

* `start` is the initial _carry_ value,
* `step` is a variety (functional value) which accepts a default (aka monadic) argument,
* `stop` is an optional variety which accepts a default argument and returns a boolean.

`Loop` invokes `step` in a loop. On the first invocation, `Loop` passes the value
of `start` to the default argument of `step`. On following invocations, `Loop` passes
the value returned by the previous invocation of `step`. We call the value returned 
by `step` a _carry_ value.

If `stop` is provided, `Loop` will call `stop` after each invocation of `step` passing
it the value returned by the immediately preceding `step` invocation.

* If `stop` returns `true`, looping stops and `Loop` returns the last carry value
(which triggered `stop` to return `true`.)

* If `stop` returns `false`, looping continues.

If `stop` is not provided, looping will never end.

## EXAMPLE: COMPUTING FIBONACCI NUMBERS ITERATIVELY

In the example below, function `IterativeFib` computes Fibonacci numbers
using an iterative algorithm. 

The state of the algorithm between iterations is captured by a structure
with three fields:

* Field `n` is an integer, holding the index of the Fibonacci number to be
computed on the next iteration,
* Field `prev` is an integer, holding the value of the `(n-1)`-st Fibonacci number, and
* Field `prevPrev` is an integer, holding the value of the `(n-2)`-nd Fibonacci number.

The starting state is `(n: 2, prev: 1, prevPrev: 1)`, reflecting the fact 
the `0`-th and `1`-st Fibonacci numbers both equal `1` and the index of
the next Fibonacci number to be computed is `2`.

Given a state structure, say `state`, each iteration of the algorithm
proceeds as follows:

* Compute the `state.n`-th Fibonacci number as the sum of `state.prev` and `state.prevPrev`,
* Return a new state object with:
	* Field `n` set to `state.n` plus `1`,
	* Field `prev` set to the `state.n`-th Fibonacci number (computed above),
	* Field `prevPrev` set to the `state.prev`

The implementation follows:

	// IterativeFib computes the n-th Fibonacci number, using an iterative algorithm.
	// n must be a number bigger than 1.
	IterativeFib(n?) {
		return: Loop(
			// n = Fibonacci number to compute on the next call to step
			// prev = the value of the (n-1)-st Fibonacci number
			// prevPrev = the value of the (n-2)-nd Fibonacci number
			start: (n: 2, prev: 1, prevPrev: 1)   // prepare the initial iterative state
			step: iterativeFibStep
			stop: iterativeFibStop[n: n]
		).prev
	}

	iterativeFibStep(state?) {
		fibN: Sum(state.prev, state.prevPrev)   // compute the state.n-th Fibonacci number
		return: (   // return the updated state
			n: Sum(state.n, 1)   // n <- n+1
			prev: fibN   // prev <- state.n-th Fibonacci number
			prevPrev: state.prev   // prevPrev <- (state.n-1)-st Fibonacci number
		)
	}

	iterativeFibStop(n, state?) {
		return: Less(n, state.n)
	}

You can try this example by running:

	IterativeFib40() {
		return: IterativeFib(40)
	}

You can run this with:

	ko play github.com/kocircuit/kocircuit/lessons/examples/IterativeFib40
