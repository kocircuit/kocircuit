# DEBUGGING WITH PEEK

Ko provides two builtin functions, `Peek` and `PeekType`, to enable 
convenient integration of debugging into your program flow.
These functions are the "`printf`" of Ko, for debugging purposes.

`Peek` and `PeekType` are essentially identical to their counterparts `Show` and `ShowType`,
described in the previous article on logging.
The only difference is that `Peek` and `PeekType`
also print the program stack, describing the location from which they were invoked
as well as source code location information.

For instance, executing `PeekTypeExample` below will invoke `PeekType`
from within the subroutine call to `peekTypeSubroutine`. You will see this
reflect in the print out.

	PeekTypeExample() {
		return: peekTypeSubroutine(arg: 3.14)
	}

	peekTypeSubroutine(arg) {
		return: PeekType(arg)
	}

You can run this example from the repo with:

	ko play github.com/kocircuit/kocircuit/lessons/examples/PeekTypeExample

It will produce an output like this:

	[0] span:37a4j2u, github.com/kocircuit/kocircuit/lessons/examples/debug.ko:42:10, PeekTypeExample:return
	[1] span:276hkds, github.com/kocircuit/kocircuit/lessons/examples/debug.ko:46:10, peekTypeSubroutine:return
	[2] span:33y9q36, github.com/kocircuit/kocircuit/lessons/examples/debug.ko:46:10, PeekType
	(span:33y9q36, kocircuit/lang/go/eval/symbol/invoke.go:9, github.com/kocircuit/kocircuit/lessons/examples/debug.ko:46:10) Float64

Notice that you can simply "erase" the invocation of `PeekType` from the implementation of `peekTypeSubroutine`
without affecting the result. This makes it easy to temporarily
"insert" invocations to `Peek` or `PeekType` into pre-existing code for debugging purposes.
