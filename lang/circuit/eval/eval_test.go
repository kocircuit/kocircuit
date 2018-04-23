package eval

import (
	"fmt"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

var testShapeFlow = []struct {
	File   string
	Result Shape
}{
	{
		File: `
		Return(this) { return: this }
		Call(this) { x: this(), return: x }
		Main() { return: Call( this: Return[ this: 8 ] ) }`,
		Result: Integer{Value_: 8},
	},
	{
		File: `
		Main() {
			x: ( field: "abc" )
			return: x.field
		}`,
		Result: String{Value_: "abc"},
	},
	{
		File: `
		Main() {
			x: ( field: Unknown[ name: "abc" ] )
			return: x.field()
		}`,
		Result: testUnknown{Payload: "abc"},
	},
	{
		File: `
		Return(x) { return: x }
		ReturnReturn(fmt) {
			return: Return(x: fmt)
		}
		Main() {
			return: ReturnReturn(fmt: "a", fmt: "b")
		}`,
		Result: Knot{
			{Name: "fmt", Shape: String{Value_: "a"}},
			{Name: "fmt", Shape: String{Value_: "b"}},
		},
	},
	{
		File: `
		Main() {
			return: ("a", "b")
		}`,
		Result: Knot{
			{Name: NoLabel, Shape: String{Value_: "a"}},
			{Name: NoLabel, Shape: String{Value_: "b"}},
		},
	},
	{
		File: `
		Return(etc?) { return: etc }
		Main() {
			x: Return["a", "b"]
			return: x()
		}`,
		Result: Knot{
			{Name: "etc", Shape: String{Value_: "a"}},
			{Name: "etc", Shape: String{Value_: "b"}},
		},
	},
	{
		File: `// series composition
		ReturnReturn() { return: Return }
		Return(x) { return: x }
		Main() {
			return: ReturnReturn() (x: "a")
		}`,
		Result: String{Value_: "a"},
	},
	{
		File: `// series composition + inline function definition
		Return(x) { return: x }
		Main() {
			return: ReturnReturn() {
				return: Return
			}() (x: "a")
		}`,
		Result: String{Value_: "a"},
	},
	{
		File: `// series composition + recursive inline function definition
		Main() {
			return: ReturnReturn() {
				return: Return(x) {
					return: x
				}
			}() (x: "a")
		}`,
		Result: String{Value_: "a"},
	},
}

func TestShapeFlow(t *testing.T) {
	for i, test := range testShapeFlow {
		repo, err := CompileString("test", "test.ko", test.File)
		if err != nil {
			t.Errorf("test %d: parse/grafting (%v)", i, err)
			continue
		}
		got, _, err := Program{
			Repo: repo,
			System: System{
				Faculty:  Faculty{Ideal{}: testMakeMacro{}, Ideal{Name: "Unknown"}: testUnknownMacro{}},
				Boundary: IdentityBoundary{},
				Combiner: IdentityCombiner{},
			},
		}.EvalSeq(NewSpan(), repo.Lookup("test", "Main"), Knot{})
		if err != nil {
			t.Errorf("test %d: evaluating (%v)", i, err)
			continue
		}
		if !IsSubset(test.Result, got) {
			t.Errorf("test %d: expecting %v, got %v", i, Sprint(test.Result), Sprint(got))
			continue
		}
	}
}

type testMakeMacro struct{}

func (m testMakeMacro) MacroID() string { return m.Help() }

func (testMakeMacro) MacroSheathString() *string { return PtrString("testMake") }

func (testMakeMacro) Label() string { return "testMake" }

func (testMakeMacro) Help() string { return "testMake" }

func (testMakeMacro) Doc() string { return "testMake" }

func (testMakeMacro) Invoke(_ *Span, arg Arg) (returns Return, effect Effect, err error) {
	return arg, nil, nil
}

type testUnknownMacro struct{}

func (m testUnknownMacro) MacroID() string { return m.Help() }

func (testUnknownMacro) MacroSheathString() *string { return PtrString("testUnknown") }

func (testUnknownMacro) Label() string { return "testUnknown" }

func (testUnknownMacro) Help() string { return "testUnknown" }

func (testUnknownMacro) Doc() string { return "testUnknown" }

func (testUnknownMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	pay, err := arg.(Knot).StringField("name")
	if err != nil {
		return nil, nil, span.Errorf(err, "unknown name argument")
	}
	return testUnknown{Payload: pay}, nil, nil
}

type testUnknown struct {
	Payload   interface{}
	Selection Path // selection path into the go payload type
}

func (u testUnknown) String() string { return fmt.Sprintf("unknown(%v)", u.Payload) }

func (u testUnknown) Select(span *Span, path Path) (Shape, Effect, error) {
	return testUnknown{
		Payload:   u.Payload,
		Selection: append(append(Path{}, u.Selection...), path...), // copy-and-append
	}, nil, nil
}

func (u testUnknown) Augment(span *Span, _ Knot) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "augmenting test unknown")
}

func (u testUnknown) Invoke(span *Span) (Shape, Effect, error) {
	return nil, nil, span.Errorf(nil, "invoking test unknown")
}
