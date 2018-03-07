package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type GoExpr interface {
	RenderExpr(GoFileContext) string
}

func MergeLine(line ...[]GoExpr) (merged []GoExpr) {
	for _, l := range line {
		merged = append(merged, l...)
	}
	return
}

type GoNameExpr struct {
	Origin *Span  `ko:"name=origin"`
	Name   string `ko:"name=name"`
}

func (expr *GoNameExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s_%s", expr.Name, expr.Origin.SpanID().String())
}

type GoFmtExpr struct {
	Fmt string   `ko:"name=fmt"`
	Arg []GoExpr `ko:"name=arg"`
}

func (expr *GoFmtExpr) RenderExpr(fileCtx GoFileContext) string {
	arg := make([]interface{}, len(expr.Arg))
	for i, a := range expr.Arg {
		arg[i] = a.RenderExpr(fileCtx)
	}
	return fmt.Sprintf(expr.Fmt, arg...)
}

type GoReceiveExpr struct {
	FromChan string `ko:"name=chan"`
}

func (expr *GoReceiveExpr) RenderExpr(GoFileContext) string {
	return fmt.Sprintf("<-%s", expr.FromChan)
}

type GoZeroExpr struct {
	GoType `ko:"name=type"`
}

func (expr *GoZeroExpr) RenderExpr(fileCtx GoFileContext) string {
	return expr.GoType.RenderZero(fileCtx)
}

type GoSelectExpr struct {
	Into  GoExpr `ko:"name=into"`
	Field string `ko:"name=field"`
}

func (expr *GoSelectExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s.%s", expr.Into.RenderExpr(fileCtx), expr.Field)
}

type GoCommentExpr struct {
	Comment string `ko:"name=comment"`
}

func (expr *GoCommentExpr) RenderExpr(GoFileContext) string {
	return StringComment(expr.Comment)
}

func GoVerbatimf(format string, arg ...interface{}) *GoVerbatimExpr {
	return &GoVerbatimExpr{fmt.Sprintf(format, arg...)}
}

type GoVerbatimExpr struct {
	Verbatim string `ko:"name=verbatim"`
}

func (expr *GoVerbatimExpr) RenderExpr(GoFileContext) string {
	return expr.Verbatim
}

type GoIntegerExpr struct {
	Int int `ko:"name=int"`
}

func (expr *GoIntegerExpr) RenderExpr(GoFileContext) string {
	return strconv.Itoa(expr.Int)
}

type GoQuoteExpr struct {
	String string `ko:"name=string"`
}

func (expr *GoQuoteExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%q", expr.String)
}

type GoBracketExpr struct {
	Expr GoExpr `ko:"name=expr"`
}

func (expr *GoBracketExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("(%s)", expr.Expr.RenderExpr(fileCtx))
}

type GoMakeStructExpr struct {
	For   GoType         `ko:"name=for"`
	Field []*GoFieldExpr `ko:"name=field"`
}

func (expr *GoMakeStructExpr) RenderExpr(fileCtx GoFileContext) string {
	return RenderMakeStruct(fileCtx, expr.For, expr.Field)
}

type GoFieldExpr struct {
	Field *GoField `ko:"name=field"`
	Expr  GoExpr   `ko:"name=expr"`
}

func (expr *GoFieldExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s: %s", expr.Field.Name, expr.Expr.RenderExpr(fileCtx))
}

func RenderMakeStruct(fileCtx GoFileContext, harness GoType, fromFieldExpr []*GoFieldExpr) string {
	var amp string
	switch u := harness.(type) {
	case *GoPtr:
		amp, harness = "&", u.Elem
	case *GoNeverNilPtr:
		amp, harness = "&", u.Elem
	}
	return ApplyTmpl(
		projectTmpl,
		M{
			"Amp":   amp,
			"Ref":   harness.RenderRef(fileCtx),
			"Field": renderFieldExpr(fileCtx, fromFieldExpr),
		},
	)
}

func renderFieldExpr(fileCtx GoFileContext, fromFieldExpr []*GoFieldExpr) []M {
	m := make([]M, len(fromFieldExpr))
	for i, fromFieldExpr := range fromFieldExpr {
		m[i] = M{"Name": fromFieldExpr.Field.Name, "Expr": fromFieldExpr.Expr.RenderExpr(fileCtx)}
	}
	return m
}

var projectTmpl = ParseTmpl(`{{/**/ -}}
{{.Amp}}{{.Ref}}{
{{- range .Field}}
	{{.Name}}: {{indent .Expr}},
{{- end}}
}`)

type GoSwitchExpr struct {
	Over    GoExpr              `ko:"name=over"`
	Case    []*GoSwitchCaseExpr `ko:"name=case"`
	Default GoExpr              `ko:"name=default"`
}

type GoSwitchCaseExpr struct {
	Predicate GoExpr `ko:"name=predicate"`
	Expr      GoExpr `ko:"name=expr"`
}

func (expr *GoSwitchExpr) RenderExpr(fileCtx GoFileContext) string {
	case_ := make([]M, len(expr.Case))
	for i, c := range expr.Case {
		case_[i] = M{
			"Predicate": c.Predicate.RenderExpr(fileCtx),
			"Expr":      c.Expr.RenderExpr(fileCtx),
		}
	}
	return ApplyTmpl(switchTmpl, M{
		"Over":    expr.Over.RenderExpr(fileCtx),
		"Case":    case_,
		"Default": expr.Default.RenderExpr(fileCtx),
	})
}

var switchTmpl = ParseTmpl(`{{/**/ -}}
{{- if .Over}}
switch {{.Over}} {
{{- else}}
switch {
{{- end}}
{{- range .Case}}
case {{.Predicate}}:
	{{indent .Expr}}
{{- end}}
{{- if .Default}}
default:
	{{indent .Default}}
{{- end}}
}`)

func GoBlock(line ...GoExpr) GoExpr {
	filter := []GoExpr{}
	for _, l := range line {
		switch u := l.(type) {
		case nil:
		case *GoBlockExpr:
			filter = append(filter, u.Line...)
		case GoExpr:
			filter = append(filter, l)
		}
	}
	return &GoBlockExpr{Line: filter}
}

type GoBlockExpr struct {
	Line []GoExpr `ko:"name=line"`
}

func (expr *GoBlockExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	for _, line := range expr.Line {
		fmt.Fprintln(&w, strings.TrimSpace(line.RenderExpr(fileCtx)))
	}
	return strings.TrimRight(w.String(), "\n")
}

type GoReturnExpr struct {
	Expr GoExpr `ko:"name=expr"`
}

func (expr *GoReturnExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("return %s", expr.Expr.RenderExpr(fileCtx))
}

type GoDerefExpr struct {
	Expr GoExpr `ko:"name=expr"`
}

func (expr *GoDerefExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("*%s", expr.Expr.RenderExpr(fileCtx))
}

type GoNegativeExpr struct {
	Expr GoExpr `ko:"name=expr"`
}

func (expr *GoNegativeExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("-%s", expr.Expr.RenderExpr(fileCtx))
}

type GoMonadicExpr struct {
	Left GoExpr `ko:"name=left"`
	Op   string `ko:"name=op"`
}

func (expr *GoMonadicExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s%s", expr.Left.RenderExpr(fileCtx), expr.Op)
}

type GoDyadicExpr struct {
	Left  GoExpr `ko:"name=left"`
	Op    string `ko:"name=op"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoDyadicExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s %s %s", expr.Left.RenderExpr(fileCtx), expr.Op, expr.Right.RenderExpr(fileCtx))
}

type GoColonAssignExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoColonAssignExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s := %s", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoAssignExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoAssignExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s = %s", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoColonPairAssignExpr struct {
	Left  [2]GoExpr `ko:"name=left"`
	Right [2]GoExpr `ko:"name=right"`
}

func (expr *GoColonPairAssignExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf(
		"%s, %s := %s, %s",
		expr.Left[0].RenderExpr(fileCtx),
		expr.Left[1].RenderExpr(fileCtx),
		expr.Right[0].RenderExpr(fileCtx),
		expr.Right[1].RenderExpr(fileCtx),
	)
}

type GoPairAssignExpr struct {
	Left  [2]GoExpr `ko:"name=left"`
	Right [2]GoExpr `ko:"name=right"`
}

func (expr *GoPairAssignExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf(
		"%s, %s = %s, %s",
		expr.Left[0].RenderExpr(fileCtx),
		expr.Left[1].RenderExpr(fileCtx),
		expr.Right[0].RenderExpr(fileCtx),
		expr.Right[1].RenderExpr(fileCtx),
	)
}

type GoNotExpr struct {
	Expr GoExpr `ko:"name=expr"`
}

func (expr *GoNotExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("!%s", expr.Expr.RenderExpr(fileCtx))
}

type GoAndExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoAndExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("(%s && %s)", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoOrExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoOrExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("(%s || %s)", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoEqualityExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoEqualityExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s == %s", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoInequalityExpr struct {
	Left  GoExpr `ko:"name=left"`
	Right GoExpr `ko:"name=right"`
}

func (expr *GoInequalityExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s != %s", expr.Left.RenderExpr(fileCtx), expr.Right.RenderExpr(fileCtx))
}

type GoTypeRefExpr struct {
	Type GoType `ko:"name=type"`
}

func (expr *GoTypeRefExpr) RenderExpr(fileCtx GoFileContext) string {
	return expr.Type.RenderRef(fileCtx)
}

type GoPanicExpr struct{}

func (expr *GoPanicExpr) RenderExpr(fileCtx GoFileContext) string {
	return "panic(\"o\")"
}

type GoFuncExpr interface {
	GoExpr
	NameExpr() GoExpr
}

func FuncExprRenderedName(expr GoFuncExpr, fileCtx GoFileContext) string {
	return expr.NameExpr().RenderExpr(fileCtx)
}

type GoShaperFuncExpr struct {
	Comment    string   `ko:"name=comment"`
	FuncName   GoExpr   `ko:"name=funcName"`
	ArgName    GoExpr   `ko:"name=argName"`
	ArgType    GoType   `ko:"name=argType"`
	ReturnType GoType   `ko:"name=returnType"`
	Line       []GoExpr `ko:"name=line"`
}

func (expr *GoShaperFuncExpr) NameExpr() GoExpr {
	return expr.FuncName
}

func (expr *GoShaperFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	if expr.Comment != "" {
		fmt.Fprint(&w, StringComment(expr.Comment))
	}
	fmt.Fprintf(
		&w,
		"func %s(%s %s) %s {",
		expr.FuncName.RenderExpr(fileCtx),
		expr.ArgName.RenderExpr(fileCtx),
		expr.ArgType.RenderRef(fileCtx),
		expr.ReturnType.RenderRef(fileCtx),
	)
	bodyExpr := &GoBlockExpr{Line: expr.Line}
	if body := bodyExpr.RenderExpr(fileCtx); body != "" {
		fmt.Fprint(&w, "\n\t")
		fmt.Fprintln(&w, Indent(body))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

type GoDuctFuncExpr struct {
	Comment    string   `ko:"name=comment"`
	FuncName   GoExpr   `ko:"name=funcName"`
	ArgName    GoExpr   `ko:"name=argName"`
	ArgType    GoType   `ko:"name=argType"`
	ReturnName GoExpr   `ko:"name=returnName"`
	ReturnType GoType   `ko:"name=returnType"`
	Line       []GoExpr `ko:"name=line"`
}

func (expr *GoDuctFuncExpr) NameExpr() GoExpr {
	return expr.FuncName
}

func (expr *GoDuctFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	if expr.Comment != "" {
		fmt.Fprint(&w, StringComment(expr.Comment))
	}
	returnExpr := expr.ReturnType.RenderRef(fileCtx)
	if expr.ReturnName != nil {
		returnExpr = fmt.Sprintf("(%s %s)", expr.ReturnName.RenderExpr(fileCtx), returnExpr)
	}
	fmt.Fprintf(
		&w,
		"func %s(step_ctx *runtime.Context, %s %s) %s {",
		expr.FuncName.RenderExpr(fileCtx),
		expr.ArgName.RenderExpr(fileCtx),
		expr.ArgType.RenderRef(fileCtx),
		returnExpr,
	)
	bodyExpr := &GoBlockExpr{Line: expr.Line}
	if body := bodyExpr.RenderExpr(fileCtx); body != "" {
		fmt.Fprint(&w, "\n\t")
		fmt.Fprintln(&w, Indent(body))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

type GoShapeExpr struct {
	Shaper Shaper `ko:"name=shaper"`
	Expr   GoExpr `ko:"name=expr"`
}

func (expr *GoShapeExpr) RenderExpr(fileCtx GoFileContext) string {
	return expr.Shaper.RenderExprShaping(fileCtx, expr.Expr)
}

var (
	EmptyExpr     = &GoVerbatimExpr{""}
	NilExpr       = &GoVerbatimExpr{"nil"}
	ZeroExpr      = &GoVerbatimExpr{"0"}
	OneExpr       = &GoVerbatimExpr{"1"}
	UnderlineExpr = &GoVerbatimExpr{"_"}
	BreakExpr     = &GoVerbatimExpr{"break"}
	TrueExpr      = &GoVerbatimExpr{"true"}
	FalseExpr     = &GoVerbatimExpr{"false"}
	MakeExpr      = &GoVerbatimExpr{"make"}
	LenExpr       = &GoVerbatimExpr{"len"}
	BrokenExpr    = &GoVerbatimExpr{`{{!@#%}}`}
)

// SliceOrArrayType{elem_expr1, elem_expr2}
type GoMakeSequenceExpr struct {
	Type GoType   `ko:"name=type"` // sequence type, e.g. []int or [3]*struct{}
	Elem []GoExpr `ko:"name=elem"`
}

func (expr *GoMakeSequenceExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, expr.Type.RenderRef(fileCtx))
	fmt.Fprint(&w, "{")
	for i, elem := range expr.Elem {
		fmt.Fprint(&w, elem.RenderExpr(fileCtx))
		if i+1 < len(expr.Elem) {
			fmt.Fprint(&w, ", ")
		}
	}
	fmt.Fprint(&w, "}")
	return w.String()
}

type GoIfThenExpr struct {
	If   GoExpr   `ko:"name=if"`
	Then []GoExpr `ko:"name=then"` // lines
}

func (expr *GoIfThenExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "if %s {", expr.If.RenderExpr(fileCtx))
	if len(expr.Then) > 0 {
		thenBlock := &GoBlockExpr{Line: expr.Then}
		fmt.Fprintf(&w, "\n\t%s\n", Indent(thenBlock.RenderExpr(fileCtx)))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

type GoIfThenElseExpr struct {
	If   GoExpr   `ko:"name=if"`
	Then []GoExpr `ko:"name=then"`
	Else []GoExpr `ko:"name=else"`
}

func (expr *GoIfThenElseExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "if %s {", expr.If.RenderExpr(fileCtx))
	if len(expr.Then) > 0 {
		thenBlock := &GoBlockExpr{Line: expr.Then}
		fmt.Fprintf(&w, "\n\t%s\n", Indent(thenBlock.RenderExpr(fileCtx)))
	}
	fmt.Fprintf(&w, "} else {")
	if len(expr.Else) > 0 {
		elseBlock := &GoBlockExpr{Line: expr.Else}
		fmt.Fprintf(&w, "\n\t%s\n", Indent(elseBlock.RenderExpr(fileCtx)))
	}
	fmt.Fprintf(&w, "}")
	return w.String()
}

type GoCallExpr struct {
	Func GoExpr   `ko:"name=func"`
	Arg  []GoExpr `ko:"name=arg"`
}

func (expr *GoCallExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	fmt.Fprintf(&w, "%s(", expr.Func.RenderExpr(fileCtx))
	for i, a := range expr.Arg {
		fmt.Fprint(&w, a.RenderExpr(fileCtx))
		if i+1 < len(expr.Arg) {
			fmt.Fprint(&w, ", ")
		}
	}
	fmt.Fprintf(&w, ")")
	return w.String()
}

type GoCallFuncExpr struct {
	Func GoFuncExpr `ko:"name=func"`
	Arg  []GoExpr   `ko:"name=arg"`
}

func (expr *GoCallFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	rewrite := &GoCallExpr{
		Func: expr.Func.NameExpr(),
		Arg:  expr.Arg,
	}
	return rewrite.RenderExpr(fileCtx)
}

type GoCallShaperFuncExpr struct {
	Func *GoShaperFuncExpr `ko:"name=func"`
	Arg  GoExpr            `ko:"name=arg"`
}

func (expr *GoCallShaperFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	rewrite := &GoCallFuncExpr{
		Func: expr.Func,
		Arg:  []GoExpr{expr.Arg},
	}
	return rewrite.RenderExpr(fileCtx)
}

type GoCallDuctFuncExpr struct {
	Func *GoDuctFuncExpr `ko:"name=func"`
	Arg  GoExpr          `ko:"name=arg"`
}

func (expr *GoCallDuctFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	rewrite := &GoCallFuncExpr{
		Func: expr.Func,
		Arg:  []GoExpr{&GoVerbatimExpr{"step_ctx"}, expr.Arg},
	}
	return rewrite.RenderExpr(fileCtx)
}

// go func() {
// 	█
// }()
type GoGoFuncExpr struct {
	Line []GoExpr `ko:"name=line"`
}

func (expr *GoGoFuncExpr) RenderExpr(fileCtx GoFileContext) string {
	block := &GoBlockExpr{Line: expr.Line}
	return fmt.Sprintf("go func() {\n\t%s\n}()", Indent(block.RenderExpr(fileCtx)))
}

// for RANGE {
// 	█
// }
type GoForExpr struct {
	Range GoExpr   `ko:"name=range"`
	Line  []GoExpr `ko:"name=line"`
}

func (expr *GoForExpr) RenderExpr(fileCtx GoFileContext) string {
	block := &GoBlockExpr{Line: expr.Line}
	if expr.Range != nil {
		return fmt.Sprintf("for %s {\n\t%s\n}",
			expr.Range.RenderExpr(fileCtx),
			Indent(block.RenderExpr(fileCtx)),
		)
	} else {
		return fmt.Sprintf("for {\n\t%s\n}", Indent(block.RenderExpr(fileCtx)))
	}
}

type GoIncrementExpr struct {
	Zero      GoExpr `ko:"name=zero"`
	Invariant GoExpr `ko:"name=invariant"`
	Increment GoExpr `ko:"name=increment"`
}

func (expr *GoIncrementExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf(
		"%s; %s; %s",
		expr.Zero.RenderExpr(fileCtx),
		expr.Invariant.RenderExpr(fileCtx),
		expr.Increment.RenderExpr(fileCtx),
	)
}

type GoVarDeclExpr struct {
	Name GoExpr `ko:"name=name"`
	Type GoType `ko:"name=type"`
}

func (expr *GoVarDeclExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf(
		"var %s %s",
		expr.Name.RenderExpr(fileCtx),
		expr.Type.RenderRef(fileCtx),
	)
}

type GoListExpr struct {
	Elem []GoExpr `ko:"name=expr"`
}

func (expr *GoListExpr) RenderExpr(fileCtx GoFileContext) string {
	var w bytes.Buffer
	for i, e := range expr.Elem {
		w.WriteString(e.RenderExpr(fileCtx))
		if i+1 < len(expr.Elem) {
			w.WriteString(", ")
		}
	}
	return w.String()
}

type GoRangeExpr struct {
	Range GoExpr `ko:"name=range"`
}

func (expr *GoRangeExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("range %s", expr.Range.RenderExpr(fileCtx))
}

// append(base, e0, e1, e2)
type GoAppendExpr struct {
	Base GoExpr   `ko:"name=base"`
	Elem []GoExpr `ko:"name=elem"`
}

func (expr *GoAppendExpr) RenderExpr(fileCtx GoFileContext) string {
	list := &GoListExpr{Elem: expr.Elem}
	return fmt.Sprintf("append(%s)", list.RenderExpr(fileCtx))
}

type GoIndexExpr struct {
	Container GoExpr `ko:"name=container"`
	Index     GoExpr `ko:"name=index"`
}

func (expr *GoIndexExpr) RenderExpr(fileCtx GoFileContext) string {
	return fmt.Sprintf("%s[%s]", expr.Container.RenderExpr(fileCtx), expr.Index.RenderExpr(fileCtx))
}

type GoInitExpr struct {
	Line []GoExpr `ko:"name=line"`
}

func (expr *GoInitExpr) RenderExpr(fileCtx GoFileContext) string {
	block := &GoBlockExpr{Line: expr.Line}
	return fmt.Sprintf("func init() {\n\t%s\n}", Indent(block.RenderExpr(fileCtx)))
}
