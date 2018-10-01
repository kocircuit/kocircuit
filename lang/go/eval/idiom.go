package eval

import (
	"github.com/kocircuit/kocircuit/lang/circuit/compile"
	"github.com/kocircuit/kocircuit/lang/circuit/eval"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

var EvalIdiomRepo = compileIdiomRepo()

func compileIdiomRepo() model.Repo {
	repo := compile.MustCompileString(eval.IdiomRootPkg, "idiom.ko", EvalIdiomSrc) // package "idiom"
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

Inc(number?) {
	return: Sum(number, 1)
}

// RunTests runs a set of test varieties.
// packages is a series of (name, tests) pairs.
// tests is a series of (name, func) pairs.
RunTests(packages?) {
	ranged: Range(
		over: packages                // range of all packages
		with: runTestPackageIterator  // iterator function
		start: (passed: 0, failed: 0) // initial carry
	)
	results: ranged.residue
	after: Show(results: results)
	return: ranged
}

runTestPackageIterator(carry, elem) {
	before: Show(package: elem.name)
	ranged: Range(
		over: elem.tests              // range of all tests
		with: runTestIterator         // iterator function
		start: (passed: 0, failed: 0) // initial carry
		_before: before
	)
	pkgResults: ranged.residue
	after: Show(package: elem.name, passed: pkgResults.passed, failed: pkgResults.failed)
	result: (
		passed: Sum(carry.passed, pkgResults.passed)
		failed: Sum(carry.failed, pkgResults.failed)
	)
	return: (
		carry: result
	)
}

runTestIterator(carry, elem) {
	testName: elem.name
	testResult: Recover(
		invoke: runTestRunner[testFunc: elem.func, testName: testName]
		panic: runTestFailRecover[testName: testName]
	)
	result: (
		passed: Sum(carry.passed, testResult.passed)
		failed: Sum(carry.failed, testResult.failed)
	)
	return: (
		carry: result
	)
}

runTestRunner(testFunc, testName) {
	s1: testFunc()
	s2: Show(pass: testName, _after: ignore(s1))
	return: (
		passed: 1
		failed: 0
	)
}

runTestFailRecover(panicValue?, testName) {
	s1: Show(fail: testName)
	return: (
		passed: 0
		failed: 1
	)
}

ignore(value?) {
	return: Empty()
}
`
