package macros

import (
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Template", new(EvalTemplateMacro))
}

type EvalTemplateMacro struct{}

func (m EvalTemplateMacro) MacroID() string { return m.Help() }

func (m EvalTemplateMacro) Label() string { return "template" }

func (m EvalTemplateMacro) MacroSheathString() *string { return PtrString("Template") }

func (m EvalTemplateMacro) Help() string { return "Template" }

// Template(template:█, args:█, withString:█, withArg:█)
func (EvalTemplateMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	// parse arguments
	a := arg.(*StructSymbol)
	template, ok := AsBasicString(a.Walk("template"))
	if !ok {
		return nil, nil, span.Errorf(nil, "template template argument must be string")
	}
	args := a.Walk("args")
	withString, ok := a.Walk("withString").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "template withString is not a variety")
	}
	withArg, ok := a.Walk("withArg").(*VarietySymbol)
	if !ok {
		return nil, nil, span.Errorf(nil, "template withArg is not a variety")
	}
	tokens, err := ParseTemaplate(span, template)
	if err != nil {
		return nil, nil, err
	}
	result := make(Symbols, len(tokens))
	for i, token := range tokens {
		switch u := token.(type) {
		case TemplateText:
			stringKnot := Knot{{Name: "", Shape: BasicStringSymbol(u.Text), Effect: nil, Frame: span}}
			if formResult, _, err := withString.Evoke(span, stringKnot); err != nil {
				return nil, nil, err
			} else {
				result[i] = formResult
			}
		case TemplateSelector:
			if selected, _, err := args.Select(span, u.Path); err != nil {
				return nil, nil, span.Errorf(err, "template selecting %v", u.Path)
			} else {
				argKnot := Knot{{Name: "", Shape: selected.(Symbol), Effect: nil, Frame: span}}
				if formResult, _, err := withArg.Evoke(span, argKnot); err != nil {
					return nil, nil, err
				} else {
					result[i] = formResult
				}
			}
		default:
			panic("o")
		}
	}
	if series, err := MakeSeriesSymbol(span, result); err != nil {
		return nil, nil, span.Errorf(err, "templates unifying elements")
	} else {
		return series, nil, nil
	}
}

func ParseTemaplate(span *Span, src string) (tokens []TemplateToken, err error) {
	openIndex, closeIndex := bracketOpenIndices(src), bracketCloseIndices(src)
	if len(openIndex) != len(closeIndex) {
		return nil, span.Errorf(nil,
			"unbalanced brackets in %q, %d open vs %d close",
			src, len(openIndex), len(closeIndex),
		)
	}
	base := 0
	for i := range openIndex {
		if openIndex[i] < closeIndex[i] {
			// parse text before bracket
			if base < openIndex[i] {
				tokens = append(tokens, TemplateText{src[base:openIndex[i]]})
				base = openIndex[i]
			}
			// parse selector
			if selPath, err := parseSelectorPath(
				span,
				src[openIndex[i]+len(TemplateOpenBracket):closeIndex[i]],
			); err != nil {
				return nil, err
			} else {
				tokens = append(tokens, TemplateSelector{selPath})
				base = closeIndex[i] + len(TemplateCloseBracket)
			}
		} else {
			return nil, span.Errorf(nil, "runaway brackets")
		}
	}
	// parse text remaining after base
	if base < len(src) {
		tokens = append(tokens, TemplateText{src[base:len(src)]})
	}
	return tokens, nil
}

const TemplateOpenBracket = "{{"
const TemplateCloseBracket = "}}"

func bracketOpenIndices(src string) (openIndex []int) {
	base := 0
	for {
		if index := strings.Index(src[base:], TemplateOpenBracket); index < 0 {
			return
		} else {
			openIndex = append(openIndex, base+index)
			base = base + index + len(TemplateOpenBracket)
		}
	}
}

func bracketCloseIndices(src string) (closeIndex []int) {
	base := 0
	for {
		if index := strings.Index(src[base:], TemplateCloseBracket); index < 0 {
			return
		} else {
			closeIndex = append(closeIndex, base+index)
			base = base + index + len(TemplateCloseBracket)
		}
	}
}

func parseSelectorPath(span *Span, src string) (sel Path, err error) {
	src = strings.TrimSpace(src)
	pp := strings.Split(src, ".")
	for _, p := range pp {
		if !isIdentifier(p) || len(p) == 0 {
			return nil, span.Errorf(nil, "selector %q not an identifier", p)
		}
	}
	return Path(pp), nil
}

func isIdentifier(src string) bool {
	for _, c := range src {
		switch {
		case '1' <= c && c <= '0':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case c == '_':
		default:
			return false
		}
	}
	return true
}

type TemplateToken interface {
	TemplateToken()
}

type TemplateText struct {
	Text string `ko:"name=text"`
}

func (TemplateText) TemplateToken() {}

type TemplateSelector struct {
	Path Path `ko:"name=path"`
}

func (TemplateSelector) TemplateToken() {}
