package macros

import (
	"log"
	"sync"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalMacro("Spin", new(EvalSpinMacro))
}

type EvalSpinMacro struct{}

func (m EvalSpinMacro) MacroID() string { return m.Help() }

func (m EvalSpinMacro) Label() string { return "spin" }

func (m EvalSpinMacro) MacroSheathString() *string { return PtrString("Spin") }

func (m EvalSpinMacro) Help() string { return "Spin" }

func (EvalSpinMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if vty, ok := arg.(*StructSymbol).SelectMonadic().(*VarietySymbol); !ok {
		return nil, nil, span.Errorf(nil, "spin expects a variety, got %v", arg)
	} else {
		done := make(chan *waitResult, 1)
		go func() {
			defer func() {
				if r := recover(); r != nil {
					evalPanic := r.(*EvalPanic)
					log.Println(
						evalPanic.Origin.Errorf(nil, "panic inside spin: %v", evalPanic.Panic),
					)
					done <- &waitResult{Panic: evalPanic, Returned: EmptySymbol{}}
					close(done)
				}
			}()
			if returned, _, invErr := vty.Invoke(span); invErr != nil {
				log.Printf("spin error (%v)", invErr)
				done <- &waitResult{Error: invErr, Returned: EmptySymbol{}}
			} else {
				done <- &waitResult{Returned: returned.(Symbol)}
			}
			close(done)
		}()
		return MakeStructSymbol(
			FieldSymbols{
				{
					Name: "Wait",
					Value: MakeVarietySymbol(
						&evalWaitMacro{waiter: &waiter{done: done}},
						nil,
					),
				},
			},
		), nil, nil
	}
}

type waitResult struct {
	Panic    *EvalPanic `ko:"name=panic"`
	Error    error      `ko:"name=error"`
	Returned Symbol     `ko:"name=returned"`
}

type waiter struct {
	sync.Mutex
	done   chan *waitResult
	result *waitResult
}

func (w *waiter) Wait() *waitResult {
	w.Lock()
	defer w.Unlock()
	if wr, ok := <-w.done; ok {
		w.result = wr
	}
	return w.result
}

type evalWaitMacro struct {
	waiter *waiter
}

func (m *evalWaitMacro) MacroID() string { return m.Help() }

func (m *evalWaitMacro) Label() string { return "wait" }

func (m *evalWaitMacro) MacroSheathString() *string { return PtrString("Wait") }

func (m *evalWaitMacro) Help() string { return "Wait" }

func (m *evalWaitMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	wr := m.waiter.Wait()
	switch {
	case wr.Error != nil:
		return nil, nil, span.Errorf(err, "spinned function error")
	case wr.Panic != nil:
		panic(wr.Panic)
	default:
		return wr.Returned, nil, nil
	}
}
