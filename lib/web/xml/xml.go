package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoRenderElement))
	RegisterEvalGate(new(GoTagElement))
	RegisterEvalGate(new(GoDocTypeElement))
	RegisterEvalGate(new(GoTextElement))
	RegisterEvalGate(new(GoGroupElement))
}

// Elements: GoTagElement, GoDocTypeElement, GoTextElement, GoGroupElement
type Element interface {
	WriteXML(io.Writer)
}

type GoTagElement struct {
	Tag      string    `ko:"name=tag"`
	Attrs    []string  `ko:"name=attrs"`
	Children []Element `ko:"name=children"`
}

func (g *GoTagElement) Play(ctx *runtime.Context) Element {
	return g
}

func (g *GoTagElement) WriteXML(w io.Writer) {
	if len(g.Children) == 0 {
		fmt.Fprintf(w, "<%s %s />", g.Tag, strings.Join(g.Attrs, " "))
	} else {
		fmt.Fprintf(w, "<%s %s>", g.Tag, strings.Join(g.Attrs, " "))
		for i, child := range g.Children {
			child.WriteXML(w)
			if i+1 < len(g.Children) {
				w.Write([]byte(" "))
			}
		}
		fmt.Fprintf(w, "</%s>", g.Tag)
	}
}

type GoDocTypeElement struct {
	String_ string `ko:"name=string,monadic"`
}

func (g *GoDocTypeElement) Play(ctx *runtime.Context) Element {
	return g
}

func (g *GoDocTypeElement) WriteXML(w io.Writer) {
	fmt.Fprintf(w, "<!DOCTYPE %s>", g.String_)
}

type GoTextElement struct {
	String_ string `ko:"name=string,monadic"`
}

func (g *GoTextElement) Play(ctx *runtime.Context) Element {
	return g
}

func (g *GoTextElement) WriteXML(w io.Writer) {
	xml.EscapeText(w, []byte(g.String_))
}

type GoGroupElement struct {
	Elem []Element `ko:"name=elem,monadic"`
}

func (g *GoGroupElement) Play(ctx *runtime.Context) Element {
	return g
}

func (g *GoGroupElement) WriteXML(w io.Writer) {
	for i, e := range g.Elem {
		e.WriteXML(w)
		if i+1 < len(g.Elem) {
			w.Write([]byte(" "))
		}
	}
}

type GoRenderElement struct {
	Elem []Element `ko:"name=elem,monadic"`
}

func (g *GoRenderElement) Play(ctx *runtime.Context) string {
	var w bytes.Buffer
	for _, elem := range g.Elem {
		elem.WriteXML(&w)
	}
	return w.String()
}
