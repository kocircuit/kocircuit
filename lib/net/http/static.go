package http

import (
	"net/http"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterEvalGate(new(GoServeLocalDir))
	RegisterGoGate(new(GoServeLocalDir))
}

type GoServeLocalDir struct {
	Address string `ko:"name=address"`
	Dir     string `ko:"name=dir"`
	URL     string `ko:"name=url"` // root URL
}

func (g *GoServeLocalDir) Play(ctx *runtime.Context) bool {
	http.Handle(g.URL, http.FileServer(http.Dir(g.Dir)))
	ctx.Printf("serving %q on %s%s", g.Dir, g.Address, g.URL)
	if err := http.ListenAndServe(g.Address, nil); err != nil {
		ctx.Fatalf("serving %q on %s%s (%v)", g.Dir, g.Address, g.URL, err)
		return false // never reached
	} else {
		return true
	}
}
