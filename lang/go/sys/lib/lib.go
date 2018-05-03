package lib

import (
	// macros
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros/circuit"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/weave"
	// standard library
	_ "github.com/kocircuit/kocircuit/lib/encoding/json"
	_ "github.com/kocircuit/kocircuit/lib/encoding/yaml"
	_ "github.com/kocircuit/kocircuit/lib/file"
	_ "github.com/kocircuit/kocircuit/lib/file/path"
	_ "github.com/kocircuit/kocircuit/lib/integer"
	_ "github.com/kocircuit/kocircuit/lib/net/http"
	_ "github.com/kocircuit/kocircuit/lib/shell"
	_ "github.com/kocircuit/kocircuit/lib/strings"
	_ "github.com/kocircuit/kocircuit/lib/sync"
	_ "github.com/kocircuit/kocircuit/lib/time"
	_ "github.com/kocircuit/kocircuit/lib/web/xml"
	// bootstrap
	_ "github.com/kocircuit/kocircuit/bootstrap/types"
)
