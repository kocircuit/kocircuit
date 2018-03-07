package path

import (
	"path"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterEvalGate(new(GoJoinPath))
	RegisterEvalGate(new(GoBasePath))
	RegisterEvalGate(new(GoDirPath))
	RegisterGoGate(new(GoJoinPath))
	RegisterGoGate(new(GoBasePath))
	RegisterGoGate(new(GoDirPath))
}

type GoJoinPath struct {
	Paths []string `ko:"name=paths,monadic"`
}

func (g *GoJoinPath) Play(ctx *runtime.Context) string {
	return path.Join(g.Paths...)
}

type GoBasePath struct {
	Path string `ko:"name=path,monadic"`
}

func (g *GoBasePath) Play(ctx *runtime.Context) string {
	return path.Base(g.Path)
}

type GoDirPath struct {
	Path string `ko:"name=path,monadic"`
}

func (g *GoDirPath) Play(ctx *runtime.Context) string {
	return path.Dir(g.Path)
}
