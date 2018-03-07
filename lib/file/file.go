package file

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterEvalGate(new(GoWriteLocalFile))
	RegisterGoGate(new(GoWriteLocalFile))
}

type GoWriteLocalFile struct {
	Path string `ko:"name=path"`
	Body string `ko:"name=body"`
}

func (g *GoWriteLocalFile) Play(ctx *runtime.Context) bool {
	if err := os.MkdirAll(path.Dir(g.Path), 0755); err != nil {
		ctx.Fatalf("mkdir %q (%v)", path.Dir(g.Path), err)
		return false
	}
	if err := ioutil.WriteFile(g.Path, []byte(g.Body), 0666); err != nil {
		ctx.Fatalf("write file %q (%v)", g.Path, err)
		return false
	}
	return true
}
