package json

import (
	"encoding/json"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoMarshal))
	RegisterEvalGate(new(GoMarshalIndent))
}

type GoMarshal struct {
	Value interface{} `ko:"name=value,monadic"`
}

func (g *GoMarshal) Play(ctx *runtime.Context) string {
	buf, err := json.Marshal(g.Value)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

type GoMarshalIndent struct {
	Value  interface{} `ko:"name=value,monadic"`
	Prefix *string     `ko:"name=prefix"`
	Indent *string     `ko:"name=indent"`
}

func (g *GoMarshalIndent) Play(ctx *runtime.Context) string {
	buf, err := json.MarshalIndent(
		g.Value,
		OptString(g.Prefix, ""),
		OptString(g.Indent, "\t"),
	)
	if err != nil {
		panic(err)
	}
	return string(buf)
}
