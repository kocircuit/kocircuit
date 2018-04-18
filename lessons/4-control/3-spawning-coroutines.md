# SPAWNING CO-ROUTINES

Ko provides a mechanism for spawning a function execution
in a new co-routine and eventually retrieving
the result (or panic value) from the spawned execution.

The builtin `Spin` function executes an argument function in a new co-routine.

`Spin` expects a single unnamed variety (functional value) argument.
It creates a new co-routine and executes its argument function there.

Once execution commences `Spin` returns a "handle" structure.
The handle structure contains a single field, called `Wait`, which holds a functional value.

Calling `Wait` will block until the underlying spawned execution completes.
When it does, `Wait` returns the value returned by the spawned function.
If the spawned execution panics, `Wait` will reproduce that same panic (and its value).

The following example computes two Fibonacci numbers, placing each computation
in a separate co-routine. 

	ComputeTwoFibonacciNumbers(n1, n2) {
		handle1: Spin(IterativeFib[n1])
		handle2: Spin(IterativeFib[n2])
		return: (
			fib1: handle1.Wait()
			fib2: handle2.Wait()
		)
	}

The function `IterativeFib` for computing a Fibonacci number was covered in a previous
article, and its source can be found in:

	github.com/kocircuit/kocircuit/lessons/examples/loop.ko

You can try `ComputeTwoFibonacciNumbers` by running, for instance:

	ComputeFib17AndFib42() {
		return: ComputeTwoFibonacciNumbers(n1: 17, n2: 42)
	}

Run this with:

	ko play github.com/kocircuit/kocircuit/lessons/examples/ComputeFib17AndFib42
