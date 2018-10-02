//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package strings

import (
	"bytes"
	"strconv"
	"strings"

	go_eval "github.com/kocircuit/kocircuit/lang/go/eval"
	"github.com/kocircuit/kocircuit/lang/go/kit/util"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

func init() {
	go_eval.RegisterNamedEvalGate("Equal", new(goEqualStrings))
	go_eval.RegisterNamedEvalGate("Len", new(goLenStrings))
	go_eval.RegisterNamedEvalGate("Join", new(goJoinStrings))
	go_eval.RegisterNamedEvalGate("Tree", new(goStringTree))
	go_eval.RegisterNamedEvalGate("Flush", new(goFlush))
	go_eval.RegisterNamedEvalGate("Quote", new(goQuote))
	go_eval.RegisterNamedEvalGate("HasPrefix", new(goHasPrefix))
	go_eval.RegisterNamedEvalGate("HasSuffix", new(goHasSuffix))
}

// goEqualStrings implements the Equal(string?) function.
type goEqualStrings struct {
	String []string `ko:"name=string,monadic"`
}

func (g *goEqualStrings) Play(ctx *runtime.Context) bool {
	for i := 1; i < len(g.String); i++ {
		if g.String[i] != g.String[i-1] {
			return false
		}
	}
	return true
}

func (g *goEqualStrings) Help() string {
	return "Equal(string?)"
}

func (g *goEqualStrings) Doc() string {
	return "Equal(string?) returns true if all given strings are equal"
}

// goStringTree implements the Tree(delimited,prefix,middle,suffix) function.
type goStringTree struct {
	Delimiter *string              `ko:"name=delimiter"` // delimiter between renditions of the middle trees
	Prefix    []string             `ko:"name=prefix"`
	Middle    []goStringTreeOpaque `ko:"name=middle"`
	Suffix    []string             `ko:"name=suffix"`
}

type goStringTreeOpaque interface {
	GoStringTree() *goStringTree
}

func (tree *goStringTree) GoStringTree() *goStringTree {
	return tree
}

func (tree *goStringTree) Play(ctx *runtime.Context) goStringTreeOpaque {
	return tree
}

// goFlush implements the Flush(tree) function.
type goFlush struct {
	Tree []goStringTreeOpaque `ko:"name=tree,monadic"`
}

func (flush *goFlush) Play(ctx *runtime.Context) string {
	var w bytes.Buffer
	for _, tree := range flush.Tree {
		if _, err := flushTree(&w, tree.GoStringTree()); err != nil {
			ctx.Fatalf("flushing string tree (%v)", err)
		}
	}
	return w.String()
}

func flushTree(w *bytes.Buffer, tree *goStringTree) (n int, err error) {
	delimiter := util.OptString(tree.Delimiter, "")
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

// goJoinStrings implements the Join(string?, delimited) function
type goJoinStrings struct {
	String    []string `ko:"name=string,monadic"`
	Delimiter *string  `ko:"name=delimiter"`
}

func (g *goJoinStrings) Play(ctx *runtime.Context) string {
	return strings.Join(g.String, util.OptString(g.Delimiter, ""))
}

func (g *goJoinStrings) Help() string {
	return "Join(string?, delimiter)"
}

func (g *goJoinStrings) Doc() string {
	return "Join(string?, delimiter) concatenates all given strings with an optional delimited between all strings"
}

// goLenStrings implements the Len(string?) function
type goLenStrings struct {
	String []string `ko:"name=string,monadic"`
}

func (g *goLenStrings) Play(ctx *runtime.Context) (n int) {
	for _, s := range g.String {
		n += len(s)
	}
	return
}

func (g *goLenStrings) Help() string {
	return "Len(string?)"
}

func (g *goLenStrings) Doc() string {
	return "Len(string?) returns the length of all given strings combined"
}

// goQuote implements to Quote(string?) function
type goQuote struct {
	String string `ko:"name=string,monadic"`
}

func (g *goQuote) Play(ctx *runtime.Context) string {
	return strconv.Quote(g.String)
}

func (g *goQuote) Help() string {
	return "Quote(string?)"
}

func (g *goQuote) Doc() string {
	return "Quote(string?) returns a double-quotes Ko string representing the given string"
}

// goHasPrefix implements to HasPrefix(string?, prefix) function
type goHasPrefix struct {
	String string `ko:"name=string,monadic"`
	Prefix string `ko:"name=prefix"`
}

func (g *goHasPrefix) Play(ctx *runtime.Context) bool {
	return strings.HasPrefix(g.String, g.Prefix)
}

func (g *goHasPrefix) Help() string {
	return "HasPrefix(string?, prefix)"
}

func (g *goHasPrefix) Doc() string {
	return "HasPrefix(string?, prefix) returns true if the given string starts with the given prefix"
}

// goHasSuffix implements to HasSuffix(string?, suffix) function
type goHasSuffix struct {
	String string `ko:"name=string,monadic"`
	Suffix string `ko:"name=suffix"`
}

func (g *goHasSuffix) Play(ctx *runtime.Context) bool {
	return strings.HasSuffix(g.String, g.Suffix)
}

func (g *goHasSuffix) Help() string {
	return "HasSuffix(string?, suffix)"
}

func (g *goHasSuffix) Doc() string {
	return "HasSuffix(string?, suffix) returns true if the given string ends with the given suffix"
}
