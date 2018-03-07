package macros

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/eval"
	. "github.com/kocircuit/kocircuit/lang/go/eval/symbol"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
)

func init() {
	RegisterEvalPkgMacro("integer", "Equal", new(EvalIntegerEqualMacro))
	RegisterEvalPkgMacro("integer", "Less", new(EvalIntegerLessMacro))
	RegisterEvalPkgMacro("integer", "Prod", new(EvalIntegerProdMacro))
	RegisterEvalPkgMacro("integer", "Ratio", new(EvalIntegerRatioMacro))
	RegisterEvalPkgMacro("integer", "Sum", new(EvalIntegerSumMacro))
	RegisterEvalPkgMacro("integer", "Moduli", new(EvalIntegerModuliMacro))
	RegisterEvalPkgMacro("integer", "Negative", new(EvalIntegerNegativeMacro))
}

type EvalIntegerEqualMacro struct{}

func (m EvalIntegerEqualMacro) MacroID() string { return m.Help() }

func (m EvalIntegerEqualMacro) Label() string { return "equal" }

func (m EvalIntegerEqualMacro) MacroSheathString() *string { return PtrString("integer.Equal") }

func (m EvalIntegerEqualMacro) Help() string { return "integer.Equal" }

func (m EvalIntegerEqualMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, _, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer equal")
	} else {
		for i := 1; i < len(symbols); i++ {
			if !symbols[i-1].Equal(symbols[i]) {
				return BasicFalse, nil, nil
			}
		}
		return BasicTrue, nil, nil
	}
}

type EvalIntegerLessMacro struct{}

func (m EvalIntegerLessMacro) MacroID() string { return m.Help() }

func (m EvalIntegerLessMacro) Label() string { return "less" }

func (m EvalIntegerLessMacro) MacroSheathString() *string { return PtrString("integer.Less") }

func (m EvalIntegerLessMacro) Help() string { return "integer.Less" }

func (m EvalIntegerLessMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer less")
	} else if signed { // signed
		for i := 1; i < len(symbols); i++ {
			if SignedMaximal(symbols[i-1]) >= SignedMaximal(symbols[i]) {
				return BasicFalse, nil, nil
			}
		}
		return BasicTrue, nil, nil
	} else { // unsigned
		for i := 1; i < len(symbols); i++ {
			if UnsignedMaximal(symbols[i-1]) >= UnsignedMaximal(symbols[i]) {
				return BasicFalse, nil, nil
			}
		}
		return BasicTrue, nil, nil
	}
}

type EvalIntegerProdMacro struct{}

func (m EvalIntegerProdMacro) MacroID() string { return m.Help() }

func (m EvalIntegerProdMacro) Label() string { return "product" }

func (m EvalIntegerProdMacro) MacroSheathString() *string { return PtrString("integer.Prod") }

func (m EvalIntegerProdMacro) Help() string { return "integer.Prod" }

func (m EvalIntegerProdMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer product")
	} else if signed {
		prod := SignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			prod *= SignedMaximal(symbols[i])
		}
		return makeAndConvert(span, prod, symbols[0].Type().(BasicType))
	} else { // unsigned
		prod := UnsignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			prod *= UnsignedMaximal(symbols[i])
		}
		return makeAndConvert(span, prod, symbols[0].Type().(BasicType))
	}
}

type EvalIntegerRatioMacro struct{}

func (m EvalIntegerRatioMacro) MacroID() string { return m.Help() }

func (m EvalIntegerRatioMacro) Label() string { return "ratio" }

func (m EvalIntegerRatioMacro) MacroSheathString() *string { return PtrString("integer.Ratio") }

func (m EvalIntegerRatioMacro) Help() string { return "integer.Ratio" }

func (m EvalIntegerRatioMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer ratio")
	} else if signed {
		ratio := SignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			ratio /= SignedMaximal(symbols[i])
		}
		return makeAndConvert(span, ratio, symbols[0].Type().(BasicType))
	} else { // unsigned
		ratio := UnsignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			ratio /= UnsignedMaximal(symbols[i])
		}
		return makeAndConvert(span, ratio, symbols[0].Type().(BasicType))
	}
}

type EvalIntegerSumMacro struct{}

func (m EvalIntegerSumMacro) MacroID() string { return m.Help() }

func (m EvalIntegerSumMacro) Label() string { return "sum" }

func (m EvalIntegerSumMacro) MacroSheathString() *string { return PtrString("integer.Sum") }

func (m EvalIntegerSumMacro) Help() string { return "integer.Sum" }

func (m EvalIntegerSumMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer sum")
	} else if signed {
		sum := SignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			sum += SignedMaximal(symbols[i])
		}
		return makeAndConvert(span, sum, symbols[0].Type().(BasicType))
	} else { // unsigned
		sum := UnsignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			sum += UnsignedMaximal(symbols[i])
		}
		return makeAndConvert(span, sum, symbols[0].Type().(BasicType))
	}
}

type EvalIntegerModuliMacro struct{}

func (m EvalIntegerModuliMacro) MacroID() string { return m.Help() }

func (m EvalIntegerModuliMacro) Label() string { return "moduli" }

func (m EvalIntegerModuliMacro) MacroSheathString() *string { return PtrString("integer.Moduli") }

func (m EvalIntegerModuliMacro) Help() string { return "integer.Moduli" }

func (m EvalIntegerModuliMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer moduli")
	} else if signed {
		moduli := SignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			moduli %= SignedMaximal(symbols[i])
		}
		return makeAndConvert(span, moduli, symbols[0].Type().(BasicType))
	} else { // unsigned
		moduli := UnsignedMaximal(symbols[0])
		for i := 1; i < len(symbols); i++ {
			moduli %= UnsignedMaximal(symbols[i])
		}
		return makeAndConvert(span, moduli, symbols[0].Type().(BasicType))
	}
}

type EvalIntegerNegativeMacro struct{}

func (m EvalIntegerNegativeMacro) MacroID() string { return m.Help() }

func (m EvalIntegerNegativeMacro) Label() string { return "negative" }

func (m EvalIntegerNegativeMacro) MacroSheathString() *string { return PtrString("integer.Negative") }

func (m EvalIntegerNegativeMacro) Help() string { return "integer.Negative" }

func (m EvalIntegerNegativeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, signed, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer negative")
	} else {
		if len(symbols) != 1 {
			return nil, nil, span.Errorf(err, "integer negative expects a single integer")
		} else {
			if signed { // signed
				negated := -SignedMaximal(symbols[0])
				return makeAndConvert(span, negated, symbols[0].Type().(BasicType))
			} else { // unsigned
				negated := -UnsignedMaximal(symbols[0])
				return makeAndConvert(span, negated, symbols[0].Type().(BasicType))
			}
		}
	}
}

func makeAndConvert(span *Span, v interface{}, to BasicType) (returns Return, effect Effect, err error) {
	if r, err := MakeBasicSymbol(span, v).ConvertTo(span, to); err != nil {
		panic(err)
	} else {
		return r, nil, nil
	}
}
