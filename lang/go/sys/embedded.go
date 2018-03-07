package sys

import (
	// Register weave macros
	_ "github.com/kocircuit/kocircuit/lang/go/weave/boolean"
	_ "github.com/kocircuit/kocircuit/lang/go/weave/integer"
	// Register eval macros
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros/circuit"
	// Packages with side-effect of registering go gate implementations
	_ "github.com/kocircuit/kocircuit/lib/file"
	_ "github.com/kocircuit/kocircuit/lib/file/path"
	_ "github.com/kocircuit/kocircuit/lib/integer"
	_ "github.com/kocircuit/kocircuit/lib/net/http"
	_ "github.com/kocircuit/kocircuit/lib/shell"
	_ "github.com/kocircuit/kocircuit/lib/strings"
	_ "github.com/kocircuit/kocircuit/lib/time"
	_ "github.com/kocircuit/kocircuit/lib/web/xml"
	// _ "github.com/kocircuit/kocircuit/x/tree"
)
