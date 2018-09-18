# DOCUMENTATION

The Ko compilation environment consists of a set of functions available to the user.
Some functions are built into the compiler and some are user-defined (in the repo).

All functions reside in a package, identified by a path.

The Ko compiler provides two commands, `list` and `doc`, to view the available
builtin functions and to display documentation for both builtin and user-defined functions
and packages.

## LIST OF BUILTIN FUNCTIONS

The Ko `list` command lists all functions that are built into the compiler
(this excludes user-defined functions):

```bash
ko list
```

## PACKAGE DOCUMENTATION

To view the contents and documentation of a package, use a command of the form:

```bash
ko doc <pkg_path>...
```

For instance, to display the documentation for
package `"github.com/kocircuit/kocircuit/lib/strings"`, run:

```bash
ko doc github.com/kocircuit/kocircuit/lib/strings...
```

## FUNCTION DOCUMENTATION

To view the documentation for any function, use a command of the form:

```bash
ko doc <pkg_path>/<func_name>
```

For instance, to display the documentation for function `Join`
in package `"github.com/kocircuit/kocircuit/lib/strings"`, run:

```bash
ko doc github.com/kocircuit/kocircuit/lib/strings/Join
```

A few builtin functions belong to the empty string, `""`, package,
such as `Range`, `Yield`, etc. To view the documentation for those,
use a command of the form:

```bash
ko doc <func_name>
```

For instance, to display the documentation for `Yield`, run:

```bash
ko doc Yield
```