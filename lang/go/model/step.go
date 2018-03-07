// Package model provides a language model for Go implementations of Ko programs.
package model

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

// GoStep captures the invocation of a substep logic.
// (Steps are generated in the go combiner.)
type GoStep struct {
	// Comment string       `ko:"name=comment" ctx:"expand"` // multiline comment
	Span    *Span        `ko:"name=span"`
	Label   string       `ko:"name=label"`   // step label
	Arrival []*GoArrival `ko:"name=receive"` // receive map from channels to fields
	Send    []string     `ko:"name=send"`    // labels of downstream channels waiting on result
	Returns GoType       `ko:"name=returns"`
	Result  bool         `ko:"name=result"` // this step's return value is returned by the enclosing func
	Logic   GoStepLogic  `ko:"name=logic"`  // step logic
	Cached  *AssignCache `ko:"name=cached"`
}

func (step *GoStep) Splay() Tree {
	return Quote{step.Comment()}
}

func (step *GoStep) Comment() string {
	return fmt.Sprintf("(GoStep: %s)", step.Span.SourceLine())
}

type GoArrival struct {
	FromExpr GoExpr         `ko:"name=expr"`
	FromChan *GoReceiveExpr `ko:"name=chan"`   // chan receptions also process upstream panics
	Shaper   Shaper         `ko:"name=shaper"` // reshaping of arriving data
	Slot     Slot           `ko:"name=slot"`
}

func (arrival *GoArrival) SlotExpr() GoExpr {
	return &GoVerbatimExpr{fmt.Sprintf("slot_%s", arrival.Slot.Label())}
}

// If arrival is from a channel, return the expression:
//	slot_result_SLOT := <-ch
//	if slot_result_SLOT.Recover() != nil {
//		panic(upstream_signal{})
//	}
func (arrival *GoArrival) FormReception() []GoExpr {
	if arrival.FromChan == nil {
		return nil
	}
	block := &GoBlockExpr{
		Line: []GoExpr{
			&GoColonAssignExpr{
				Left:  &GoVerbatimExpr{fmt.Sprintf("slot_result_%s", arrival.Slot.Label())},
				Right: arrival.FromChan,
			},
			&GoIfThenExpr{
				If:   &GoVerbatimExpr{fmt.Sprintf("slot_result_%s.Recover() != nil", arrival.Slot.Label())},
				Then: []GoExpr{&GoVerbatimExpr{"panic(upstream_signal{})"}},
			},
		},
	}
	return []GoExpr{block}
}

// If arrival is from a channel, return:
//	// shaper info
//	slot_SLOT := shape(slot_result_SLOT.Returned)
//	_ = slot_SLOT
// If arrival is from an expression, return:
//	// shaper info
//	slot_SLOT := shape(expr)
//	_ = slot_SLOT
func (arrival *GoArrival) FormShaping() []GoExpr {
	block := &GoBlockExpr{
	// Line: []GoExpr{
	// 	&GoCommentExpr{Sprint(arrival.Shaper)},
	// },
	}
	if arrival.FromChan != nil {
		block.Line = append(
			block.Line,
			&GoColonAssignExpr{
				Left: arrival.SlotExpr(),
				Right: &GoShapeExpr{
					Shaper: arrival.Shaper,
					Expr:   &GoVerbatimExpr{fmt.Sprintf("slot_result_%s.Returned", arrival.Slot.Label())},
				},
			},
		)
	} else {
		block.Line = append(
			block.Line,
			&GoColonAssignExpr{
				Left: arrival.SlotExpr(),
				Right: &GoShapeExpr{
					Shaper: arrival.Shaper,
					Expr:   arrival.FromExpr,
				},
			},
		)
	}
	block.Line = append(
		block.Line,
		&GoAssignExpr{
			Left:  UnderlineExpr,
			Right: arrival.SlotExpr(),
		},
	)
	return []GoExpr{block}
}

func ArrivalToSlotExpr(arrival []*GoArrival) (arg []*GoSlotExpr) {
	arg = make([]*GoSlotExpr, len(arrival))
	for i, arrival := range arrival {
		arg[i] = &GoSlotExpr{
			Slot: arrival.Slot,
			Expr: arrival.SlotExpr(),
		}
	}
	return
}

func FormArrivals(arrival []*GoArrival) GoExpr {
	block := &GoBlockExpr{}
	block.Line = append(block.Line, &GoCommentExpr{"wait until all arguments arrive"})
	for _, a := range arrival {
		block.Line = append(block.Line, a.FormReception()...)
	}
	block.Line = append(block.Line, &GoCommentExpr{"apply implicit type conversions"})
	for _, a := range arrival {
		block.Line = append(block.Line, a.FormShaping()...)
	}
	return block
}

func (step *GoStep) Render(goCircuit *GoCircuit, fileCtx GoFileContext) string {
	return ApplyTmpl(
		ParStepImplTmpl,
		M{
			"Circuit": goCircuit.Valve.Address.RenderExpr(fileCtx),
			"Step":    step,
			"Line":    step.Span.SourceLine(),
			"StepSyntax": fmt.Sprintf("(%s) label=%s",
				step.Span.SourceLine(), step.Label,
			),
			"Arrivals":      FormArrivals(step.Arrival).RenderExpr(fileCtx),
			"StepLogic":     step.Logic.FormExpr(ArrivalToSlotExpr(step.Arrival)...).RenderExpr(fileCtx),
			"StepLogicType": fmt.Sprintf("%T", step.Logic),
		},
	)
}

var ParStepImplTmpl = ParseTmpl(`{{/**/ -}}
{{/* comment .Step.Comment */ -}}
{{comment .StepSyntax -}}
kill_{{.Step.Label}} := make(chan struct{})
go func() {
	step_ctx := &runtime.Context{
		Parent: ctx,
		Source: {{printf "%q" .Step.Comment}}, {{/* TODO: use standardized source loc data */}}
		Context: ctx.Context,
		Kill: kill_{{.Step.Label}},
	}
	defer func() {
		{{/* every step always sends a result to its send channels (this is the panic path) */ -}}
		{{/* catch panics from upstream or across (kill signal due to another step failing) */ -}}
		if recovered := recover(); recovered != nil {
			sr := &step_{{.Circuit}}_{{.Step.Label}}_result{Ctx: step_ctx}
			switch u := recovered.(type) {
			case kill_signal, upstream_signal:
				sr.Panic, sr.GoStack = u, __debug.Stack()
			default: {{/* this step produced the panic */}}
				sr.Panic, sr.GoStack = recovered, __debug.Stack()
				abort <- sr {{/* abort is fired only by panics arising from this step itself */}}
			}
			{{- range .Step.Send}}
			{{.}} <- sr; close({{.}})
			{{- end}}
			done <- true
			{{- if .Step.Result}}
			result <- sr
			{{- end}}
		}
	}()
	{{indent .Arrivals}}
	// ({{.Line}}) check for kill signal, before committing to execution
	select {
	case <-kill_{{.Step.Label}}:
		panic(kill_signal{})
	default:
	}
	// ({{.Line}}) invoke step logic {{.StepLogicType}}
	{{- if or .Step.Send .Step.Result}}
	step_logic_returned := {{indent .StepLogic}}
	{{- else}}
	_ = {{indent .StepLogic}}
	{{- end}}
	{{- /* every step always sends a result to its send channels (this is the success path) */ -}}
	{{- if or .Step.Send .Step.Result}}
	// ({{.Line}}) send result
	sr := &step_{{.Circuit}}_{{.Step.Label}}_result{Ctx: step_ctx, Returned: step_logic_returned}
	{{- end}}
	{{- range .Step.Send}}
	{{.}} <- sr; close({{.}})
	{{- end}}
	done <- true
	{{- if .Step.Result}}
	result <- sr
	{{- end}}
}()
`)
