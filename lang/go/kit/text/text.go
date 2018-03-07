package text

import (
	"fmt"
	"io"

	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type TextRenderer interface {
	RenderText() Textual
}

type TextualCtx interface { //XXX:add to Textual.Render
	Rewrite(Textual) Textual
}

type Textual interface { // implemented by TextRubber, TextSlab and TextGlue
	Render(ctx PrintContext, w io.Writer, linePrefix string, width int)
	Width() int
}

// TextSlab is a verbatim/indivisible string.
type TextSlab struct {
	String string `ko:"name=string,monadic"`
}

func (t TextSlab) Play(*runtime.Context) TextSlab { return t }

func (t TextSlab) Width() int {
	return len(t.String)
}

func (t TextSlab) Render(_ PrintContext, w io.Writer, _ string, _ int) {
	fmt.Fprint(w, t.String)
}

// TextTile is a verbatim string printed with a highlight.
type TextTile struct {
	String string `ko:"name=string,monadic"`
}

func (t TextTile) Play(*runtime.Context) TextTile { return t }

func (t TextTile) Width() int {
	return len(t.String)
}

func (t TextTile) Render(ctx PrintContext, w io.Writer, _ string, _ int) {
	fmt.Fprint(w, ctx.Console.Field(t.String))
}

// TextRubber is a sequence of textual objects which can render horizontally or vertically,
// depending on available width.
type TextRubber struct {
	Header *string   `ko:"name=header"`
	Open   *string   `ko:"name=open"`  // e.g. "{"
	Close  *string   `ko:"name=close"` // e.g. "}"
	Field  []Textual `ko:"name=text,monadic"`
}

func (g TextRubber) header() string { return OptString(g.Header, "") }
func (g TextRubber) open() string   { return OptString(g.Open, "") }
func (g TextRubber) close() string  { return OptString(g.Close, "") }

func (g TextRubber) Play(*runtime.Context) TextRubber { return g }

func (g TextRubber) Width() int {
	s := len(g.header()) + len(g.open()) + len(g.close()) + 2*len(g.Field)
	for _, f := range g.Field {
		s += f.Width()
	}
	return s
}

func (g TextRubber) Render(ctx PrintContext, w io.Writer, linePrefix string, width int) {
	unfold := g.Width() > width
	fmt.Fprint(w, g.header(), g.open())
	for i, f := range g.Field {
		if unfold {
			fmt.Fprint(w, "\n", linePrefix, ctx.Indent)
		} else if i > 0 {
			fmt.Fprint(w, ", ")
		}
		f.Render(ctx, w, linePrefix+ctx.Indent, width-len(ctx.Indent))
	}
	if unfold && len(g.Field) > 0 {
		fmt.Fprint(w, "\n", linePrefix)
	}
	if g.close() != "" {
		fmt.Fprint(w, g.close())
	}
}

// TextGlue is a list of textual objects which cannot be broken up at their boundary.
type TextGlue struct {
	Text []Textual `ko:"name=text,monadic"`
}

func (g TextGlue) Play(*runtime.Context) TextGlue { return g }

func (g TextGlue) Width() int {
	s := 0
	for _, r := range g.Text {
		s += r.Width()
	}
	return s
}

func (g TextGlue) Render(ctx PrintContext, w io.Writer, linePrefix string, width int) {
	for _, r := range g.Text {
		r.Render(ctx, w, linePrefix, width)
	}
}
