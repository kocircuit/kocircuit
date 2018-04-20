package strings

import (
	"bytes"
	"strconv"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	RegisterEvalGate(new(GoEqualStrings))
	RegisterEvalGate(new(GoLenStrings))
	RegisterEvalGate(new(GoJoinStrings))
	RegisterEvalGate(new(GoStringTree))
	RegisterEvalGate(new(GoFlush))
	RegisterEvalGate(new(GoQuote))
}

type GoEqualStrings struct {
	String_ []string `ko:"name=string,monadic"`
}

func (g *GoEqualStrings) Play(ctx *runtime.Context) bool {
	for i := 1; i < len(g.String_); i++ {
		if g.String_[i] != g.String_[i-1] {
			return false
		}
	}
	return true
}

type GoStringTree struct {
	Delimiter *string              `ko:"name=delimiter"` // delimiter between renditions of the middle trees
	Prefix    []string             `ko:"name=prefix"`
	Middle    []GoStringTreeOpaque `ko:"name=middle"`
	Suffix    []string             `ko:"name=suffix"`
}

type GoStringTreeOpaque interface {
	GoStringTree() *GoStringTree
}

func (tree *GoStringTree) GoStringTree() *GoStringTree {
	return tree
}

func (tree *GoStringTree) Play(ctx *runtime.Context) GoStringTreeOpaque {
	return tree
}

type GoFlush struct {
	Tree []GoStringTreeOpaque `ko:"name=tree,monadic"`
}

func (flush *GoFlush) Play(ctx *runtime.Context) string {
	var w bytes.Buffer
	for _, tree := range flush.Tree {
		if _, err := flushTree(&w, tree.GoStringTree()); err != nil {
			ctx.Fatalf("flushing string tree (%v)", err)
		}
	}
	return w.String()
}

func flushTree(w *bytes.Buffer, tree *GoStringTree) (n int, err error) {
	delimiter := OptString(tree.Delimiter, "")
	if prefix := strings.Join(tree.Prefix, ""); prefix != "" {
		if k, err := w.WriteString(prefix); err != nil {
			return 0, err
		} else {
			n += k
		}
	}
	for i, middle := range tree.Middle {
		if i > 0 {
			if k, err := w.WriteString(delimiter); err != nil {
				return 0, err
			} else {
				n += k
			}
		}
		if k, err := flushTree(w, middle.GoStringTree()); err != nil {
			return 0, err
		} else {
			n += k
		}
	}
	if suffix := strings.Join(tree.Suffix, ""); suffix != "" {
		if k, err := w.WriteString(suffix); err != nil {
			return 0, err
		} else {
			n += k
		}
	}
	return
}

type GoJoinStrings struct {
	String_   []string `ko:"name=string,monadic"`
	Delimiter *string  `ko:"name=delimiter"`
}

func (g *GoJoinStrings) Play(ctx *runtime.Context) string {
	return strings.Join(g.String_, OptString(g.Delimiter, ""))
}

type GoLenStrings struct {
	String_ []string `ko:"name=string,monadic"`
}

func (g *GoLenStrings) Play(ctx *runtime.Context) (n int) {
	for _, s := range g.String_ {
		n += len(s)
	}
	return
}

type GoQuote struct {
	String_ string `ko:"name=string,monadic"`
}

func (g *GoQuote) Play(ctx *runtime.Context) string {
	return strconv.Quote(g.String_)
}
