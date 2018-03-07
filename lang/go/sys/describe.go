package sys

import (
	"bytes"
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
	. "github.com/kocircuit/kocircuit/lang/go/weave"
)

func init() {
	RegisterGoGateAt("", "Print", new(GoPrint))
	RegisterGoGateAt("", "Println", new(GoPrintln))
	RegisterEvalGateAt("", "Print", new(GoPrint))
	RegisterEvalGateAt("", "Println", new(GoPrintln))
}

type GoPrint struct {
	String []string `ko:"name=string,monadic"`
}

func (p *GoPrint) Play(ctx *runtime.Context) string {
	switch len(p.String) {
	case 0:
		fmt.Print()
		return ""
	case 1:
		fmt.Print(p.String[0])
		return ""
	default:
		var w bytes.Buffer
		for _, s := range p.String {
			w.WriteString(s)
		}
		fmt.Print(w.String())
		return w.String()
	}
}

type GoPrintln struct {
	String []string `ko:"name=string,monadic"`
}

func (p *GoPrintln) Play(ctx *runtime.Context) string {
	q := &GoPrint{
		String: append(append([]string{}, p.String...), "\n"),
	}
	return q.Play(ctx)
}
