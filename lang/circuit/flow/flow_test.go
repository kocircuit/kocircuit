package flow

import (
	"fmt"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestSeqFlowPlayer(t *testing.T) {
	for i, test := range testFlowPlayer {
		repo, err := CompileString("test", "test.ko", test.File)
		if err != nil {
			t.Errorf("test %d: parse/grafting (%v)", i, err)
			continue
		}
		flow := &testFlow{Repo: repo}
		fmt.Println(repo.BodyString())
		got, _, err := PlaySeqFlow(NewSpan(), repo.Lookup("test", "Main"), flow)
		if err != nil {
			t.Errorf("test %d: play flow (%v)", i, err)
			continue
		}
		if err := VerifyIsSubset(test.Flow, got); err != nil {
			t.Errorf("test %d: expecting %v, got %v (%v)", i, Sprint(test.Flow), Sprint(got), err)
			continue
		}
	}
}

func TestParFlowPlayer(t *testing.T) {
	for i, test := range testFlowPlayer {
		repo, err := CompileString("test", "test.ko", test.File)
		if err != nil {
			t.Errorf("test %d: parse/grafting (%v)", i, err)
			continue
		}
		flow := &testFlow{Repo: repo}
		fmt.Println(repo.BodyString())
		got, _, err := PlayParFlow(NewSpan(), repo.Lookup("test", "Main"), flow)
		if err != nil {
			t.Errorf("test %d: play flow (%v)", i, err)
			continue
		}
		if err := VerifyIsSubset(test.Flow, got); err != nil {
			t.Errorf("test %d: expecting %v, got %v (%v)", i, Sprint(test.Flow), Sprint(got), err)
			continue
		}
	}
}

var testFlowPlayer = []struct {
	File string
	Flow *testFlow
}{
	{
		File: `
		Main() {
			call: Returns8
			return: call( unnecessary_argument: 3 )
		}
		Returns8() { return: 8 }
		`,
		Flow: func() *testFlow {
			main := &testFlow{}
			_, _ = main.Enter(newTestFrame("0_enter", Enter{}))
			three, _ := main.Make(newTestFrame("1", Number{int64(3)}), int64(3))
			dfunc, _ := main.MakePkgFunc(newTestFrame("call", PkgFunc{"test", "Returns8"}), "test", "Returns8")
			dfunc.Augment(newTestFrame("2", Augment{}), []GatherFlow{{Field: "unnecessary_argument", Flow: three}})
			result, _ := dfunc.(*testFlow).InvokeReturn(
				newTestFrame("return", Invoke{}),
				func(frame *Span, envelope *testFlow) (Flow, error) {
					_, _ = envelope.Enter(newTestFrame("0_enter", Enter{}))
					eight, _ := envelope.Make(newTestFrame("1", Number{LexInteger{Int64: 8}}), LexInteger{Int64: 8})
					return eight.Leave(newTestFrame("0_leave", Leave{}))
				},
			)
			done, _ := result.Leave(newTestFrame("0_leave", Leave{}))
			done.(*testFlow).Frame = nil // erase span information
			return done.(*testFlow)
		}(),
	},
}

func newTestFrame(stepLabel string, stepLogic Logic) *Span {
	return NewSpan()
}

type testFlow struct {
	Repo      Repo
	Frame     *Span // frame that created this flow
	Number    interface{}
	Selected  []string // accumulated selection
	FuncPkg   string
	FuncName  string
	OpRef     []string
	Augmented []GatherFlow
}

func (f *testFlow) Copy(frame *Span) *testFlow {
	return &testFlow{
		Repo: f.Repo, Frame: frame,
		Number: f.Number, FuncPkg: f.FuncPkg, FuncName: f.FuncName, OpRef: f.OpRef, Augmented: f.Augmented,
	}
}

func (f *testFlow) newChild(frame *Span) *testFlow {
	return &testFlow{Repo: f.Repo, Frame: frame}
}

func (f *testFlow) Enter(frame *Span) (Flow, error) {
	return f.Copy(frame), nil
}

func (f *testFlow) Make(frame *Span, v interface{}) (Flow, error) {
	if f.Number != nil {
		return nil, fmt.Errorf("overwriting literal %v with %v", f.Number, v)
	}
	g := f.newChild(frame)
	g.Number = v
	return g, nil
}

func (f *testFlow) MakePkgFunc(frame *Span, pkg string, fu string) (Flow, error) {
	if f.FuncPkg != "" || f.FuncName != "" {
		return nil, fmt.Errorf(
			"overwriting function gate %q.%s with %q.%s",
			f.FuncPkg, f.FuncName, pkg, fu,
		)
	}
	g := f.newChild(frame)
	g.FuncPkg, g.FuncName = pkg, fu
	return g, nil
}

func (f *testFlow) MakeOp(frame *Span, ref []string) (Flow, error) {
	if f.OpRef != nil {
		return nil, fmt.Errorf("overwriting operator gate %v with %v", f.OpRef, ref)
	}
	g := f.newChild(frame)
	g.OpRef = ref
	return g, nil
}

func (f *testFlow) Link(frame *Span, name string, monadic bool) (Flow, error) {
	panic("o")
}

func (f *testFlow) Select(frame *Span, path []string) (Flow, error) {
	g := f.Copy(frame)
	g.Selected = append(g.Selected, path...)
	return g, nil
}

func (f *testFlow) Augment(frame *Span, gather []GatherFlow) (Flow, error) {
	g := f.Copy(frame)
	g.Augmented = append(g.Augmented, gather...)
	return g, nil
}

func (f *testFlow) Invoke(frame *Span) (Flow, error) {
	g := f.Copy(frame) // envelope flow for child function
	r, _, err := PlaySeqFlow(frame, f.Repo[f.FuncPkg][f.FuncName], g)
	if err != nil {
		return nil, err
	}
	return r.(*testFlow).Copy(frame), nil
}

func (f *testFlow) InvokeReturn(frame *Span, bypassPlay func(*Span, *testFlow) (Flow, error)) (Flow, error) {
	g := f.Copy(frame) // envelope flow for child function
	r, err := bypassPlay(frame, g)
	if err != nil {
		return nil, err
	}
	return r.(*testFlow).Copy(frame), nil
}

func (f *testFlow) Leave(frame *Span) (Flow, error) {
	return f.Copy(frame), nil
}
