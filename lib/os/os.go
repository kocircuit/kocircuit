package os

import (
	"os"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	// weave
	RegisterGoGate(new(GoMkdir))
	RegisterGoGate(new(GoEnv))
	RegisterGoGate(new(GoTempDir))
	// eval
	RegisterEvalGate(new(GoMkdir))
	RegisterEvalGate(new(GoEnv))
	RegisterEvalGate(new(GoTempDir))
}

type GoMkdir struct {
	Path string `ko:"name=path,monadic"`
}

func (mkdir *GoMkdir) Play(ctx *runtime.Context) error {
	return os.MkdirAll(mkdir.Path, 0755)
}

type GoEnv struct {
	Name string `ko:"name=name,monadic"`
}

func (env *GoEnv) Play(ctx *runtime.Context) string {
	return os.Getenv(env.Name)
}

type GoTempDir struct{}

func (GoTempDir) Play(ctx *runtime.Context) string {
	return os.TempDir()
}
