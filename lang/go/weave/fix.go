package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Fix", new(GoFixMacro))
}

type GoFixMacro struct{}

func (m GoFixMacro) MacroID() string { return m.Help() }

func (GoFixMacro) Label() string { return "fix" }

func (m GoFixMacro) MacroSheathString() *string { return PtrString("Fix") }

func (m GoFixMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoFixMacro) Invoke(span *Span, arg Arg) (Return, Effect, error) {
	if monadic := StructureMonadicField(arg.(GoStructure)); monadic == nil {
		return nil, nil, span.Errorf(nil, "fix expects a monadic argument, got %s", Sprint(arg))
	} else {
		_, logic, err := GoSelectSimplify(span, Path{monadic.KoName()}, arg.(GoStructure))
		if err != nil {
			return nil, nil, span.Errorf(err, "fix expects a func argument")
		}
		varietal, ok := logic.(GoVarietal)
		if !ok {
			return nil, nil, span.Errorf(nil, "fix func must be a single-branch variety")
		}
		if _, returns, effect, err := GoFix(span, fmt.Sprintf("fix_%s", varietal.TypeID()), varietal); err != nil {
			return nil, nil, err
		} else {
			return returns, effect, nil
		}
	}
}

func GoFix(span *Span, chamber string, varietal GoVarietal) (*GoValve, Return, Effect, error) {
	if _, err := VerifyUniqueFieldAugmentations(span, varietal); err != nil {
		return nil, nil, nil, span.Errorf(err, "fixing augmentation")
	}
	xform, err := VarietalMacroTransforms(RefineChamber(span, chamber), varietal)
	if err != nil {
		return nil, nil, nil, err
	}
	if len(xform.MacroTransform) != 1 {
		return nil, nil, nil, span.Errorf(nil, "fix expects a single-branch variety, got %s", Sprint(varietal))
	}
	valve := xform.MacroTransform[0].ExpandValve
	if valve == nil {
		return nil, nil, nil, span.Errorf(nil, "fix expects a function variety, got %s", Sprint(varietal))
	}
	returns := NewGoVariety(span, &GoCallMacro{Valve: valve}, nil)
	projected, _ := VarietalProject(span, returns)
	return valve,
		returns,
		&GoMacroEffect{
			Arg:           NewGoEmpty(span),
			SlotForm:      &FixSoln{Variety: returns},
			Cached:        xform.Cached,
			CircuitEffect: xform.CircuitEffect.AggregateDuctType(projected),
			ProgramEffect: xform.ProgramEffect,
		}, nil
}

type FixSoln struct {
	Variety *GoVariety `ko:"name=variety"` // valve variety
}

func (fix *FixSoln) FormExpr(...*GoSlotExpr) GoExpr {
	return fix.Variety.VarietyExpr()
}
