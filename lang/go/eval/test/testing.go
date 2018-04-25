package test

import (
	"fmt"
	"reflect"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	_ "github.com/kocircuit/kocircuit/lang/go/eval/macros"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	"github.com/kocircuit/kocircuit/lang/go/runtime"
)

type EvalTests struct {
	T    *testing.T  `ko:"name=t"`
	Test []*EvalTest `ko:"name=test"`
}

func (tests *EvalTests) Play(ctx *runtime.Context) struct{} {
	for i, test := range tests.Test {
		if err := test.Play(ctx); err != nil {
			tests.T.Errorf("test %d (%v)", i, err)
		}
	}
	return struct{}{}
}

type EvalTest struct {
	Name    string      `ko:"name=name"` // test name
	Enabled bool        `ko:"name=enabled"`
	File    string      `ko:"name=file"`   // test ko source
	Arg     interface{} `ko:"name=arg"`    // arg for Main
	Result  interface{} `ko:"name=result"` // expecting result
}

func (test *EvalTest) Play(ctx *runtime.Context) error {
	if !test.Enabled {
		return nil
	}
	repo, err := CompileString("test", "test.ko", test.File)
	if err != nil {
		return fmt.Errorf("compile (%v)", err)
	}
	// fmt.Println(repo["test"].BodyString())
	eval := NewEvaluator(EvalFaculty(), repo)
	span := NewSpan()
	span = RefineChamber(span, "testEval")
	span = RefineOutline(span, "Main")
	// deconstruct test arguments
	arg := Deconstruct(span, reflect.ValueOf(test.Arg))
	returned, _, _, err := eval.Eval(span, repo["test"]["Main"], arg)
	if err != nil {
		return fmt.Errorf("eval (%v)", err)
	}
	// integrate test results
	result, err := Integrate(span, returned, reflect.TypeOf(test.Result))
	if err != nil {
		return fmt.Errorf("result integration (%v)", err)
	}
	// test result values
	if !IsSubset(test.Result, result.Interface()) {
		return fmt.Errorf("expecting %s, got %s", Sprint(test.Result), Sprint(result.Interface()))
	}
	return nil
}
