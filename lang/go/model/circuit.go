// Package model provides a language model for Go implementations of Ko programs.
package model

import (
	"bytes"
	"fmt"
	"sort"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

// GoCircuit describes a Ko function (a circuit) implementation in Go.
// GoCircuit is produced by the GoCombiner.
type GoCircuit struct {
	Origin  *Func            `ko:"name=origin"`
	Comment string           `ko:"name=comment"` // go function comment
	Valve   *GoValve         `ko:"name=valve"`   // argument structure
	Step    []*GoStep        `ko:"name=step"`
	Effect  *GoCircuitEffect `ko:"name=effect"`
}

func (circuit *GoCircuit) DuctType() []GoType {
	if circuit.Effect == nil {
		return nil
	}
	return circuit.Effect.DuctType
}

func (circuit *GoCircuit) DuctFunc() []GoFuncExpr {
	if circuit.Effect == nil {
		return nil
	}
	return circuit.Effect.DuctFunc
}

func (circuit *GoCircuit) Name() string {
	return circuit.Valve.Address.Name
}

// TypeID captures the circuit's specific go arguments, go return values, and Ko namespace location.
func (circuit *GoCircuit) TypeID() string {
	return circuit.Valve.TypeID()
}

func (circuit *GoCircuit) Shapers() []Shaper {
	r := []Shaper{}
	for _, step := range circuit.Step {
		for _, a := range step.Arrival {
			r = append(r, a.Shaper)
		}
	}
	return r
}

func (circuit *GoCircuit) RenderImpl(fileCtx GoFileContext) string {
	var w bytes.Buffer
	comment := StringComment(circuit.Comment)
	fmt.Fprint(&w, comment)
	fmt.Fprint(&w, "type ")
	fmt.Fprintln(&w, circuit.Valve.RenderDef(fileCtx))
	// define the play method
	fmt.Fprintln(&w, ApplyTmpl(
		ParCircuitImplTmpl,
		M{
			"Comment":         comment,
			"TypeName":        circuit.Valve.Address.RenderExpr(fileCtx),
			"ReturnType":      circuit.Valve.Returns.RenderRef(fileCtx),
			"Steps":           circuit.RenderSteps(fileCtx),
			"ResultStepLabel": circuit.ResultStepLabel(),
		},
	))
	funcExpr := map[string]GoFuncExpr{}
	exprName := []string{}
	for _, dfx := range circuit.DuctFunc() {
		name := FuncExprRenderedName(dfx, fileCtx)
		if _, ok := funcExpr[name]; !ok {
			funcExpr[name] = dfx
			exprName = append(exprName, name)
		}
	}
	for _, shaper := range circuit.Shapers() {
		for _, sfx := range shaper.CircuitEffect().DuctFuncs() {
			name := FuncExprRenderedName(sfx, fileCtx)
			if _, ok := funcExpr[name]; !ok {
				funcExpr[name] = sfx
				exprName = append(exprName, name)
			}
		}
	}
	sort.Strings(exprName)
	for _, name := range exprName {
		w.WriteString(funcExpr[name].RenderExpr(fileCtx))
		w.WriteString("\n")
	}
	return w.String()
}

var ParCircuitImplTmpl = ParseTmpl(`{{- /* define step result harnesses */ -}}
{{- range .Steps}}
// go_circuit={{.Circuit}}, step={{.Label}}
{{comment .StepComment -}}
type step_{{.Circuit}}_{{.Label}}_result struct{
	Ctx      *runtime.Context
	Returned {{.ResultType}}
	Panic    interface{} {{/* kill_signal, upstream_signal */}}
	GoStack  []byte
}

func (sr *step_{{.Circuit}}_{{.Label}}_result) Context() *runtime.Context {return sr.Ctx }

func (sr *step_{{.Circuit}}_{{.Label}}_result) MustReturn() {{.ResultType}} {
	if sr.Panic != nil {
		panic(upstream_signal{})
	}
	return sr.Returned
}

func (sr *step_{{.Circuit}}_{{.Label}}_result) Recover() interface{} { return sr.Panic }

func (sr *step_{{.Circuit}}_{{.Label}}_result) Stack() []byte { return sr.GoStack }
{{end}}
{{/**/ -}}
{{.Comment -}}
func (arg *{{.TypeName}}) Play(ctx *runtime.Context) {{.ReturnType}} {
	{{- /* build step-to-step channels */ -}}
	{{- range .Steps}}{{range .Send}}
	{{- /* every step sends its result to each of its send channels exactly once, without blocking */}}
	{{.Chan}} := make(chan *step_{{.Circuit}}_{{.StepLabel}}_result, 1)
	{{- end}}{{end}}

	{{/* build common abort notification channel */ -}}
	{{/* every step can send to abort and done once without blocking */ -}}
	abort := make(chan runtime.Recoverer, {{len .Steps}})
	done := make(chan bool, {{len .Steps}})
	result := make(chan *step_{{.TypeName}}_{{.ResultStepLabel}}_result, 1)

	{{/* start steps */}}
	{{- range .Steps -}}
	{{indent .Impl}}
	{{end -}}

	{{/* handle first abort: send kill signal to all steps */ -}}
	var root_cause runtime.Recoverer
	go func() {
		var ok bool
		{{/* func's return logic kills this handler by closing the abort channel */ -}}
		if root_cause, ok = <-abort; !ok {
			return
		}
		{{- range .Steps}}
		close(kill_{{.Label}})
		{{- end}}
	}()
	{{/* wait until all steps complete */ -}}
	for i := 0; i < {{len .Steps}}; i++ {
		<-done
	}
	close(abort) {{/* no race: all step goroutines are done*/}}
	{{/* return */ -}}
	if root_cause != nil {
		panic(&runtime.Fault{
			Context: root_cause.Context(),
			Panic: root_cause.Recover(),
			GoStack: root_cause.Stack(),
		})
	}
	return (<-result).Returned
}
`)

type renderStep struct {
	Circuit     string
	Label       string
	ResultType  string
	Send        []*renderStepSend
	Impl        string
	StepComment string
}

type renderStepSend struct {
	Circuit   string
	Chan      string
	StepLabel string
}

func (circuit *GoCircuit) RenderSteps(fileCtx GoFileContext) []*renderStep {
	r := make([]*renderStep, len(circuit.Step))
	for i, step := range circuit.Step {
		send := make([]*renderStepSend, len(step.Send))
		for j, ch := range step.Send {
			send[j] = &renderStepSend{
				Circuit:   circuit.Valve.Address.RenderExpr(fileCtx),
				Chan:      ch,
				StepLabel: step.Label,
			}
		}
		r[i] = &renderStep{
			Circuit:     circuit.Valve.Address.RenderExpr(fileCtx),
			Label:       step.Label,
			ResultType:  step.Returns.RenderRef(fileCtx),
			Send:        send,
			Impl:        step.Render(circuit, fileCtx),
			StepComment: step.Comment(),
		}
	}
	return r
}

func (circuit *GoCircuit) ResultStepLabel() string {
	for _, s := range circuit.Step {
		if s.Result {
			return s.Label
		}
	}
	panic("result step not found")
}
