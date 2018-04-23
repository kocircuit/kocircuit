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
	//
	RegisterEvalPkgMacro("integer", "Less", new(EvalIntegerLessMacro))
	RegisterEvalPkgMacro("integer", "Prod", new(EvalIntegerProdMacro))
	RegisterEvalPkgMacro("integer", "Ratio", new(EvalIntegerRatioMacro))
	RegisterEvalPkgMacro("integer", "Sum", new(EvalIntegerSumMacro))
	RegisterEvalPkgMacro("integer", "Moduli", new(EvalIntegerModuliMacro))
	RegisterEvalPkgMacro("integer", "Negative", new(EvalIntegerNegativeMacro))
	//
	RegisterEvalMacro("Less", new(EvalIntegerLessMacro))
	RegisterEvalMacro("Prod", new(EvalIntegerProdMacro))
	RegisterEvalMacro("Ratio", new(EvalIntegerRatioMacro))
	RegisterEvalMacro("Sum", new(EvalIntegerSumMacro))
	RegisterEvalMacro("Moduli", new(EvalIntegerModuliMacro))
	RegisterEvalMacro("Negative", new(EvalIntegerNegativeMacro))
}

type EvalIntegerEqualMacro struct{}

func (m EvalIntegerEqualMacro) MacroID() string { return m.Help() }

func (m EvalIntegerEqualMacro) Label() string { return "equal" }

func (m EvalIntegerEqualMacro) MacroSheathString() *string { return PtrString("Equal") }

func (m EvalIntegerEqualMacro) Help() string { return "Equal" }

func (m EvalIntegerEqualMacro) Doc() string {
	return `Equal returns true if the sequence of integer arguments passed to it is empty,
or all integers are equal.`
}

func (m EvalIntegerEqualMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	if symbols, _, err := ExtractMonadicNonEmptyIntegerSeries(span, arg.(*StructSymbol)); err != nil {
		return nil, nil, span.Errorf(err, "integer equal")
	} else {
		for i := 1; i < len(symbols); i++ {
			if !symbols[i-1].Equal(span, symbols[i]) {
				return BasicFalse, nil, nil
			}
		}
		return BasicTrue, nil, nil
	}
}

type EvalIntegerLessMacro struct{}

func (m EvalIntegerLessMacro) MacroID() string { return m.Help() }

func (m EvalIntegerLessMacro) Label() string { return "less" }

func (m EvalIntegerLessMacro) MacroSheathString() *string { return PtrString("Less") }

func (m EvalIntegerLessMacro) Help() string { return "Less" }

func (m EvalIntegerLessMacro) Doc() string {
	return `Less returns true if the non-empty sequence of integer arguments passed to it is strictly increasing.`
}

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

func (m EvalIntegerProdMacro) MacroSheathString() *string { return PtrString("Prod") }

func (m EvalIntegerProdMacro) Help() string { return "Prod" }

func (m EvalIntegerProdMacro) Doc() string {
	return `Prod returns the product of the integers passed to it.`
}

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

func (m EvalIntegerRatioMacro) MacroSheathString() *string { return PtrString("Ratio") }

func (m EvalIntegerRatioMacro) Help() string { return "Ratio" }

func (m EvalIntegerRatioMacro) Doc() string {
	return `Ratio returns the ratio of the integers passed to it.`
}

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

func (m EvalIntegerSumMacro) MacroSheathString() *string { return PtrString("Sum") }

func (m EvalIntegerSumMacro) Help() string { return "Sum" }

func (m EvalIntegerSumMacro) Doc() string {
	return `Sum returns the sum of the integers passed to it.`
}

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

func (m EvalIntegerModuliMacro) MacroSheathString() *string { return PtrString("Moduli") }

func (m EvalIntegerModuliMacro) Help() string { return "Moduli" }

func (m EvalIntegerModuliMacro) Doc() string {
	return `Moduli(n1, n2) returns n1 modulus n2.`
}

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

func (m EvalIntegerNegativeMacro) MacroSheathString() *string { return PtrString("Negative") }

func (m EvalIntegerNegativeMacro) Help() string { return "Negative" }

func (m EvalIntegerNegativeMacro) Doc() string {
	return `Negative(n) returns -n.`
}

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
