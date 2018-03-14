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
Nothing() {
	return: ()
}

Return(pass?) {
	return: pass
}

Pick(either, or) {
	return: When(have: either, then: Return, else: Return[or])
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
	return: Take(
		from: Range(
			over: case
			with: switchCase(elem) {
				return: (emit: Yield(if: elem.if, then: elem.yield))
			}
		).image
		otherwise: otherwise
	)
}
`
