package yaml

import (
	"gopkg.in/yaml.v2"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoMarshal))
}

type GoMarshal struct {
	Value interface{} `ko:"name=value,monadic"`
}

func (g *GoMarshal) Play(ctx *runtime.Context) string {
	buf, err := yaml.Marshal(g.Value)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
