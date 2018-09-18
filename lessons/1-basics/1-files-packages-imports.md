# SOURCE FILES, PACKAGES AND IMPORTS

Ko source files reside within a repo, rooted somewhere in your local file system.
Similarly to Go, each directory within a repo represents a Ko package.

Ko is a very simple language: Every source file (and every package) provides
purely a list of function definitions.

To access the functions defined in a particular package, we use an import
statement. For instance,

```ko
import "github.com/kocircuit/kocircuit/lib/strings"
```

Functions from this package will be accessible within the importing source file
under the `strings` package alias. As in Go, the name of the package
directory itself is used as a package alias. For instance, package `strings`
contains a function named `Join` which can be referred to as
`strings.Join` within the context of the importing source file.

To prevent name collisions, Ko provides an extended import syntax which
can specify a desired package alias. Here's an example:

```ko
import "github.com/kocircuit/kocircuit/lib/strings" as stdlib_strings
```

Now the `Join` function would be accessible via `stdlib_strings.Join`.

## EXERCISE

Let's start writing a Ko program. Create a file called `helloworld.ko`
in the repo directory `github.com/kocircuit/kocircuit/lessons/examples`.
The full file path, relative to your repo root, should be
`github.com/kocircuit/kocircuit/lessons/examples/helloworld.ko`.

Import the strings package in preparation for the next lessons.

```ko
import "github.com/kocircuit/kocircuit/lib/strings"
```

The example sources are included in the repo for your convenience at
[github.com/kocircuit/kocircuit/lessons/examples/](github.com/kocircuit/kocircuit/lessons/examples/).
