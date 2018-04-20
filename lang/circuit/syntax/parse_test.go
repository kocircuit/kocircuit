package syntax

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
)

func TestParseFileString(t *testing.T) {
	for i, test := range testParse {
		parsedFile, err := ParseFileString("test.ko", test.Text)
		if err != nil {
			t.Errorf("test %d: parsing %v", i, err)
			continue
		}
		if err := VerifyIsSubset(test.Parsed, parsedFile); err != nil {
			// t.Errorf("test %d: %v", i, err)
			t.Errorf(
				"test %d: expecting %s, got %s, because %v",
				i, Sprint(test.Parsed), Sprint(parsedFile), err,
			)
			continue
		}
	}
}

var testParse = []struct {
	Text   string
	Parsed interface{}
}{
	{
		Text: `
import "path/to/nowhere" as foo
// Comment for sum.
Sum(x,y) {return: Add(x:x, y:y)}
`,
		Parsed: File{
			Path: "test.ko",
			Import: []Import{
				{Path: "path/to/nowhere", As: Ref{Path: []string{"foo"}}, Comment: ""},
			},
			Design: []Design{
				{
					Comment: " Comment for sum.\n",
					Name:    Ref{Path: []string{"Sum"}},
					Factor: []Factor{
						{Comment: "", Name: Ref{Path: []string{"x"}}},
						{Comment: "", Name: Ref{Path: []string{"y"}}},
					},
					Returns: Assembly{
						Type: "{}",
						Term: []Term{
							{
								Label: Ref{Path: []string{"return"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"Add"}},
									Type: "()",
									Term: []Term{
										{Comment: "", Label: Ref{Path: []string{"x"}}, Hitch: Ref{Path: []string{"x"}}},
										{Comment: "", Label: Ref{Path: []string{"y"}}, Hitch: Ref{Path: []string{"y"}}},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Text: `
use "path/to/nowhere" as foo
// Comment for sum.
Sum(x,y,etc?) {return: Add(x,etc,y)}
`,
		Parsed: File{
			Path: "test.ko",
			Import: []Import{
				{Path: "path/to/nowhere", As: Ref{Path: []string{"foo"}}, Comment: ""},
			},
			Design: []Design{
				{
					Comment: " Comment for sum.\n",
					Name:    Ref{Path: []string{"Sum"}},
					Factor: []Factor{
						{Comment: "", Name: Ref{Path: []string{"x"}}},
						{Comment: "", Name: Ref{Path: []string{"y"}}},
						{Comment: "", Name: Ref{Path: []string{"etc"}}, Monadic: true},
					},
					Returns: Assembly{
						Type: "{}",
						Term: []Term{
							{
								Label: Ref{Path: []string{"return"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"Add"}},
									Type: "()",
									Term: []Term{
										{Comment: "", Label: Ref{Path: []string{NoLabel}}, Hitch: Ref{Path: []string{"x"}}},
										{Comment: "", Label: Ref{Path: []string{NoLabel}}, Hitch: Ref{Path: []string{"etc"}}},
										{Comment: "", Label: Ref{Path: []string{NoLabel}}, Hitch: Ref{Path: []string{"y"}}},
									},
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Text: `
SumUp(etc) {
	return: Reduce(
		over: etc
		with: sumUpReducer(tally, elem) {
			return: Sum(tally, elem)
		}
	)
}`,
		Parsed: File{
			Path: "test.ko",
			Design: []Design{
				{
					Comment: "",
					Name:    Ref{Path: []string{"SumUp"}},
					Factor:  []Factor{{Comment: "", Name: Ref{Path: []string{"etc"}}}},
					Returns: Assembly{
						Type: "{}",
						Term: []Term{
							{
								Comment: "",
								Label:   Ref{Path: []string{"return"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"Reduce"}},
									Type: "()",
									Term: []Term{
										{Comment: "", Label: Ref{Path: []string{"over"}}, Hitch: Ref{Path: []string{"etc"}}},
										{Comment: "", Label: Ref{Path: []string{"with"}}, Hitch: Ref{Path: []string{"sumUpReducer"}}},
									},
								},
							},
						},
					},
				},
				{
					Name:   Ref{Path: []string{"sumUpReducer"}},
					Factor: []Factor{{Name: Ref{Path: []string{"tally"}}}, {Comment: "", Name: Ref{Path: []string{"elem"}}}},
					Returns: Assembly{
						Type: "{}",
						Term: []Term{{
							Comment: "",
							Label:   Ref{Path: []string{"return"}},
							Hitch: Assembly{
								Sign: Ref{Path: []string{"Sum"}},
								Type: "()",
								Term: []Term{
									{Label: Ref{Path: []string{NoLabel}}, Hitch: Ref{Path: []string{"tally"}}},
									{Comment: "", Label: Ref{Path: []string{NoLabel}}, Hitch: Ref{Path: []string{"elem"}}},
								},
							}},
						},
					},
				},
			},
		},
	},
	{
		Text: `
Main(x, y) {
	return: F(x) (y) [a:1, b:2]
}`,
		Parsed: File{
			Path: "test.ko",
			Design: []Design{
				{
					Comment: "",
					Name:    Ref{Path: []string{"Main"}},
					Factor: []Factor{
						{Comment: "", Name: Ref{Path: []string{"x"}}},
						{Comment: "", Name: Ref{Path: []string{"y"}}},
					},
					Returns: Assembly{
						Type: "{}",
						Term: []Term{
							{
								Comment: "",
								Label:   Ref{Path: []string{"return"}},
								Hitch:   Ref{Path: []string{"0_inline_0_return_0_series_2"}},
							},
							{
								Comment: "",
								Label:   Ref{Path: []string{"0_inline_0_return_0_series_0"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"F"}},
									Type: "()",
									Term: []Term{{
										Comment: "",
										Label:   Ref{Path: []string{NoLabel}},
										Hitch:   Ref{Path: []string{"x"}},
									}},
								},
							},
							{
								Comment: "",
								Label:   Ref{Path: []string{"0_inline_0_return_0_series_1"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"0_inline_0_return_0_series_0"}},
									Type: "()",
									Term: []Term{{
										Comment: "",
										Label:   Ref{Path: []string{NoLabel}},
										Hitch:   Ref{Path: []string{"y"}},
									}},
								},
							},
							{
								Comment: "",
								Label:   Ref{Path: []string{"0_inline_0_return_0_series_2"}},
								Hitch: Assembly{
									Sign: Ref{Path: []string{"0_inline_0_return_0_series_1"}},
									Type: "[]",
									Term: []Term{
										{Comment: "", Label: Ref{Path: []string{"a"}}, Hitch: Literal{Value: LexInteger{Int64: 1}}},
										{Comment: "", Label: Ref{Path: []string{"b"}}, Hitch: Literal{Value: LexInteger{Int64: 2}}},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}
