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

Branch(if, then, else) {
	return: Yield(
		if: if
		then: Yield(if: Have(then), then: then, else: [])
		else: Yield(if: Have(else), then: else, else: [])
	)()
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
