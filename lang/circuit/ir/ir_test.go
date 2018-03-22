package ir

import (
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
)

func TestSerializeDeserialize(t *testing.T) {
	repo, err := CompileString("test/pkg", "test.ko", testFile)
	if err != nil {
		t.Fatalf("compile (%v)", err)
	}
	filteredRepo, err := DeserializeRepo(SerializeRepo(repo))
	if err != nil {
		t.Fatalf("deserialize (%v)", err)
	}
	if repo.BodyString() != filteredRepo.BodyString() {
		t.Errorf("repo body strings differ")
	}
	removeNotSerializedFeatures(repo)
	removeNotSerializedFeatures(filteredRepo)
	if err := VerifyIsSubset(filteredRepo, repo); err != nil {
		t.Errorf("subset (%v)", err)
	}
	if err := VerifyIsSubset(repo, filteredRepo); err != nil {
		t.Errorf("superset (%v)", err)
	}
}

const testFile = `
G(x?) {
	label: (1, 2, 3)
	return: F(a: x, b: 1, c: "abc", d: 3.14e-12, e: true, f: label)
}
F(a, b?) {
	return: G(b)
}
`

func removeNotSerializedFeatures(repo Repo) {
	for _, pkg := range repo {
		for _, fu := range pkg {
			fu.Syntax = nil
			fu.Spread = nil // spread is based on step pointer values (cannot be used with subset verification logic)
			for _, step := range fu.Step {
				step.Syntax = nil
				step.Func = nil
			}
		}
	}
}
