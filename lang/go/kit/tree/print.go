package tree

import (
	"bytes"
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/go/kit/text"
)

func Sprint(v interface{}) string {
	return RenderTreeString(Splay(v).Text())
}

func SprintType(t reflect.Type) string {
	return RenderTreeString(Explain(t).Text())
}

func RenderTreeString(text Textual) string {
	var w bytes.Buffer
	text.Render(DefaultPrinter, &w, DefaultPrinter.Prefix, DefaultPrinter.Width)
	return w.String()
}
