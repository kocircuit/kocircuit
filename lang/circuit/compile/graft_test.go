package compile

import (
	"fmt"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestGraftFile(t *testing.T) {
	for i, test := range testGraftFile {
		if !test.Enabled {
			continue
		}
		fmt.Printf("TEST %d\n", i)
		pkgName := "ko/test"
		repo, err := CompileString(pkgName, "test.ko", test.File)
		if err != nil {
			t.Errorf("test %d: %v", i, err)
			continue
		}
		fmt.Println(repo[pkgName].BodyString())
		if err := VerifyIsSubset(test.Pkg, repo[pkgName]); err != nil {
			t.Errorf(
				"test %d: expecting %v, got %v, because %v",
				i, Sprint(test.Pkg), Sprint(repo[pkgName]), err,
			)
		}
	}
}

var testGraftFile = []struct {
	Enabled bool
	File    string
	Pkg     Package
}{
	{ // test 0
		Enabled: true,
		File: `F(etc) { 
			_: H(etc)
			__: H(etc: etc)
			return: 3.14 
		}`,
		Pkg: Package{
			"F": &Func{
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Enter: &Step{ID: StepID("0_enter"), Label: "0_enter", Logic: Enter{}},
				Field: testField{"etc"}.Make(),
				Step: testBody{
					{"_", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"__", Invoke{}, []testGather{{MainFlowLabel, "4"}}},
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"2", Augment{}, []testGather{
						{MainFlowLabel, "1"},
						{NoLabel, "0_enter_etc"},
					}},
					{"4", Augment{}, []testGather{
						{MainFlowLabel, "3"},
						{"etc", "0_enter_etc"},
					}},
					{"return", Number{LexFloat{Float64: 3.14}}, nil},
					{"1", Operator{Path: []string{"H"}}, nil},
					{"3", Operator{Path: []string{"H"}}, nil},
					{"0_enter_etc", SelectArg{Name: "etc"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{ // test 1
		Enabled: true,
		File: `
		import "path/to/pkg"
		F(x, y) { return: pkg.G(u: x, u: y) }`,
		Pkg: Package{
			"F": &Func{
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Field: testField{"x", "y"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"2", Augment{}, []testGather{
						{MainFlowLabel, "1"},
						{"u", "0_enter_x"},
						{"u", "0_enter_y"},
					}},
					{"1", PkgFunc{"path/to/pkg", "G"}, nil},
					{"0_enter_x", SelectArg{Name: "x"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_y", SelectArg{Name: "y"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{ // test 2
		Enabled: true,
		File: `
		import "path/to/pkg" as other
		H() { return: true }
		F(x, y) {
			p: H
			q: H(x: x, y: y)
			return: other.G[u: p, u: q]
		}`,
		Pkg: Package{
			"H": &Func{
				ID:   FuncID("ko/test", "H"),
				Name: "H",
				Pkg:  "ko/test",
				Step: testBody{
					{"0_enter", Enter{}, nil},
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Number{true}, nil},
				}.Make(),
			},
			"F": &Func{
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Field: testField{"x", "y"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Augment{}, []testGather{{MainFlowLabel, "3"}, {"u", "p"}, {"u", "q"}}},
					{"3", PkgFunc{"path/to/pkg", "G"}, nil},
					{"p", PkgFunc{"ko/test", "H"}, nil},
					{"q", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"2", Augment{}, []testGather{{MainFlowLabel, "1"}, {"x", "0_enter_x"}, {"y", "0_enter_y"}}},
					{"1", PkgFunc{"ko/test", "H"}, nil},
					{"0_enter_x", SelectArg{Name: "x"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_y", SelectArg{Name: "y"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{ // test 3
		Enabled: true,
		File: `
		import "path/to/pkg"
		F(g, h) {
			x: 1
			y: "a"
			return: pkg.G[u: x, u: y, u: g, w: h]
		}`,
		Pkg: Package{
			"F": &Func{
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Field: testField{"g", "h"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Augment{}, []testGather{
						{MainFlowLabel, "1"}, {"u", "x"}, {"u", "y"},
						{"u", "0_enter_g"}, {"w", "0_enter_h"},
					}},
					{"1", PkgFunc{"path/to/pkg", "G"}, nil},
					{"x", Number{LexInteger{Int64: 1}}, nil},
					{"y", Number{LexString{"a"}}, nil},
					{"0_enter_g", SelectArg{Name: "g"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_h", SelectArg{Name: "h"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{
		Enabled: true,
		File: `
		F(g, h) {
			return: f(g, h) { return: F(g: g, h: h) }
		}`,
		Pkg: Package{
			"F": &Func{
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Field: testField{"g", "h"}.Make(),
				Step: testBody{
					{"0_enter_g", SelectArg{Name: "g"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_h", SelectArg{Name: "h"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"0_enter", Enter{}, nil},
					{"return", PkgFunc{"ko/test", "f"}, nil},
				}.Make(),
			},
			"f": &Func{
				ID:    FuncID("ko/test", "f"),
				Name:  "f",
				Pkg:   "ko/test",
				Field: testField{"g", "h"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"2", Augment{}, []testGather{
						{MainFlowLabel, "1"}, {"g", "0_enter_g"}, {"h", "0_enter_h"},
					}},
					{"1", PkgFunc{"ko/test", "F"}, nil},
					{"0_enter_g", SelectArg{Name: "g"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_h", SelectArg{Name: "h"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{
		Enabled: true,
		File: `// inline functions + series composition
		F(u, v) {
			return: f(g, h) { 
				return: F(u: g, v: h) 
			} (g: u, h: v)
		}`,
		Pkg: Package{
			"F": &Func{
				Doc:   " inline functions + series composition\n",
				ID:    FuncID("ko/test", "F"),
				Name:  "F",
				Pkg:   "ko/test",
				Field: testField{"u", "v"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "0_inline_0_return_0_series_1"}}},
					{"0_inline_0_return_0_series_1", Invoke{}, []testGather{{MainFlowLabel, "1"}}},
					{"1", Augment{}, []testGather{
						{MainFlowLabel, "0_inline_0_return_0_series_0"},
						{"g", "0_enter_u"},
						{"h", "0_enter_v"},
					}},
					{"0_inline_0_return_0_series_0", PkgFunc{"ko/test", "f"}, nil},
					{"0_enter_u", SelectArg{Name: "u"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_v", SelectArg{Name: "v"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
			"f": &Func{
				ID:    FuncID("ko/test", "f"),
				Name:  "f",
				Pkg:   "ko/test",
				Field: testField{"g", "h"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "return"}}},
					{"return", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"2", Augment{}, []testGather{
						{MainFlowLabel, "1"}, {"u", "0_enter_g"}, {"v", "0_enter_h"},
					}},
					{"1", PkgFunc{"ko/test", "F"}, nil},
					{"0_enter_g", SelectArg{Name: "g"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter_h", SelectArg{Name: "h"}, []testGather{{MainFlowLabel, "0_enter"}}},
					{"0_enter", Enter{}, nil},
				}.Make(),
			},
		},
	},
	{
		Enabled: true,
		File: `
		S(v) {
			return: (v, v).etc // tests multi-edges in kahnsort
		}
		`,
		Pkg: Package{
			"S": &Func{
				ID:    FuncID("ko/test", "S"),
				Name:  "S",
				Pkg:   "ko/test",
				Field: testField{"v"}.Make(),
				Step: testBody{
					{"0_leave", Leave{}, []testGather{{MainFlowLabel, "0_inline_0_return_0_series_1"}}},
					{"0_inline_0_return_0_series_1", Select{Path: []string{"etc"}}, []testGather{{MainFlowLabel, "0_inline_0_return_0_series_0"}}},
					{"0_inline_0_return_0_series_0", Invoke{}, []testGather{{MainFlowLabel, "2"}}},
					{"2", Augment{}, []testGather{
						{MainFlowLabel, "1"},
						{NoLabel, "0_enter_v"},
						{NoLabel, "0_enter_v"},
					}},
					{"1", Operator{}, nil},
					{"0_enter_v", SelectArg{Name: "v"}, []testGather{{MainFlowLabel, "0_enter"}}},
				}.Make(),
			},
		},
	},
}

type testStep struct {
	Label  string
	Logic  Logic
	Gather []testGather
}

type testGather struct {
	Field string
	Label string
}

func (t testStep) Make() *Step {
	return &Step{ID: StepID(t.Label), Label: t.Label, Logic: t.Logic, Gather: testMakeGather(t.Gather)}
}

type testBody []*testStep

func (tt testBody) Make() []*Step {
	ss := make([]*Step, len(tt))
	for i, t := range tt {
		ss[i] = t.Make()
	}
	return ss
}

func testMakeGather(gg []testGather) []*Gather {
	g := make([]*Gather, len(gg))
	for i := range gg {
		g[i] = &Gather{
			Field: gg[i].Field,
			Step:  &Step{ID: StepID(gg[i].Label), Label: gg[i].Label},
		}
	}
	return g
}

type testField []string

func (ff testField) Make() map[string]*Step {
	s := map[string]*Step{}
	for _, f := range ff {
		var label string
		if f == NoLabel {
			label = "0_enter_nolabel"
		} else {
			label = fmt.Sprintf("0_enter_%s", f)
		}
		s[f] = testStep{
			label,
			SelectArg{Name: f, Monadic: f == NoLabel},
			[]testGather{{MainFlowLabel, "0_enter"}},
		}.Make()
	}
	return s
}
