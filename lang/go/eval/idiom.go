package eval

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

var EvalIdiomRepo = compileIdiomRepo()

func compileIdiomRepo() Repo {
	repo := MustCompileString(IdiomRootPkg, "idiom.ko", EvalIdiomSrc) // package "idiom"
	return repo
}

const EvalIdiomSrc = `
// Return returns its default argument.
Return(pass?) {
	return: pass
}

// Nothing returns an empty value.
Nothing() {
	return: Return()
}

// Empty returns an empty value.
Empty() {
	return: Return()
}

// Expect panics if its default argument is empty.
Expect(value?) {
	return: When(
		have: value
		then: Return
		else: Panic[expect: "expecting a value"]
	)
}

Have(value?) {
	return: When(
		have: value
		then: Return[pass: true]
		else: Return[pass: false]
	)
}

Pick(either, or) {
	return: When(
		have: either
		then: Return
		else: Return[or]
	)
}

Invoke(variety?) {
	return: variety()
}

Branch(if, then, else) {
	return: When(
		have: Yield(if: if, then: then, else: else)
		then: Invoke
		else: Nothing
	)
}

Switch(case, otherwise) {
	branches: switchRangeCases(case)
	took: Take(branches)
	return: When(have: took.first, then: Return, else: otherwise)
}

switchRangeCases(case?) {
	return: Range(
		over: case
		with: switchRangeCase(carry, elem) {
			invoked: elem() // invoke case
			return: (emit: invoked, carry: invoked)
		}
		stop: Have
	).image
}
`
