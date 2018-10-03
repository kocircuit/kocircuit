package test

import (
	"reflect"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
)

func TestReassembly(t *testing.T) {
	// make eval
	repo, err := CompileString("test", "test.ko", testReassemblyProg)
	if err != nil {
		t.Fatalf("compile (%v)", err)
	}
	eval := NewEvaluator(EvalFaculty(), repo)
	span := NewSpan()
	// eval pre-disassembly
	preReturned, _, _, err := eval.Eval(span, repo["test"]["Pre"], MakeStructSymbol(nil))
	if err != nil {
		t.Fatalf("eval pre-disassembly (%v)", err)
	}
	// disassemble
	pbDisassembled, err := preReturned.DisassembleToPB(span)
	if err != nil {
		t.Fatalf("disassembly (%v)", err)
	}
	// re-assemble
	reassembled, err := AssembleWithError(span, eval, pbDisassembled)
	if err != nil {
		t.Fatalf("re-assembly (%v)", err)
	}
	// eval post-reassembly
	postArg := MakeStructSymbol(
		FieldSymbols{{Name: "reassembled", Value: reassembled}},
	)
	postReturned, _, _, err := eval.Eval(span, repo["test"]["Post"], postArg)
	if err != nil {
		t.Fatalf("eval post-reassembly (%v)", err)
	}
	// verify post-eval returned true
	result, err := Integrate(span, postReturned, reflect.TypeOf(true))
	if err != nil {
		t.Fatalf("result integration (%v)", err)
	}
	if !result.Bool() {
		t.Errorf("reassembly test unexpected result")
	}
}

const testReassemblyProg = `
F(foo, bar) { return: (foo: foo, bar: bar) }
Data() {
	return: (
		integer: 7
		string: "abc"
		bool: true
		sequence: (1.1, 2.2, 3.3)
		structure: (
			variety: F[foo: "foo"]
		)
	)
}
Pre() {
	return: Data()
}
Post(reassembled) {
	_: reassembled.structure.variety(bar: "bar")
	return: Equal(reassembled, Data())
}`
