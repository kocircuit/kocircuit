package weave

import (
	"fmt"
	"os"
	"path"
	"testing"

	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/subset"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type WeaveTest struct {
	Enabled bool      `ko:"name=enabled"`
	Name    string    `ko:"name=name"`
	File    string    `ko:"name=file"`   // source
	Arg     *GoStruct `ko:"name=arg"`    // arg for main
	Result  GoType    `ko:"name=result"` // subset of expected result
}

type WeaveTests []*WeaveTest

func (test WeaveTests) Run(t *testing.T) {
	for _, test := range test {
		if !test.Enabled {
			continue
		}
		repo, err := CompileString("test", "test.ko", test.File)
		if err != nil {
			t.Errorf("test %q: compile (%v)", test.Name, err)
			continue
		}
		fmt.Println(repo["test"].BodyString())
		weaveCtx := NewGoWeaveCtx("TEST", repo, GoFaculty(), GoIdiomRepo)
		span := RefineWeaveCtx(NewSpan(), weaveCtx)
		span = RefineChamber(span, "testWeave")
		span = RefineOutline(span, "Main")
		inst, err := weaveCtx.WeaveInstrument(span, repo["test"]["Main"], test.Arg)
		if err != nil {
			fmt.Println(repo.BodyString())
			t.Errorf("test %q: weave (%v)", test.Name, err.Error())
			continue
		}
		fmt.Printf(
			"recursions=%d iterations=%d\n",
			inst.ProgramEffect.WeavingStat.RecursionCount,
			inst.ProgramEffect.WeavingStat.IterationCount,
		)
		if err := VerifyIsSubset(test.Result, inst.Returns); err != nil {
			t.Errorf("test %q: expecting %s, got %s (%v)", test.Name, Sprint(test.Result), Sprint(inst.Returns), err)
		}
		fmt.Printf("result=%s\n", Sprint(inst.Returns))
		ctx := &RenderCtx{}
		files := ctx.Shred(inst.Circuit, inst.Directive)
		tmpDir := os.TempDir()
		fmt.Printf("cp %s\n", path.Join(tmpDir, "/ko__/test/G/"))
		SourceRepo(files).Materialize(tmpDir)
	}
}

func TestGoField(comment, name string, typ GoType) *GoField {
	return &GoField{
		Comment: comment,
		Name:    GoNameFor(name),
		Type:    typ,
		Tag:     KoTags(name, false),
	}
}
