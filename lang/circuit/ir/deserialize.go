package ir

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/compile"
	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func DeserializeRepo(pbRepo *pb.Repo) (repo Repo, err error) {
	defer func() {
		if r := recover(); r != nil {
			repo, err = nil, r.(error)
		}
	}()
	return deserializeRepo(pbRepo), nil
}

func deserializeRepo(pbRepo *pb.Repo) Repo {
	repo := Repo{}
	for _, pbPkg := range pbRepo.Package {
		pkgPath, pkg := deserializePackage(pbPkg)
		repo[pkgPath] = pkg
	}
	return repo
}

func deserializePackage(pbPkg *pb.Package) (pkgPath string, pkg Package) {
	pkgPath, pkg = pbPkg.GetPath(), Package{}
	for _, pbFu := range pbPkg.Func {
		fu := deserializeFunc(pbFu)
		pkg[pbFu.GetName()] = fu
	}
	return pkgPath, pkg
}

func deserializeFunc(pbFu *pb.Func) (fu *Func) {
	fu = &Func{
		Doc:     pbFu.GetDoc(),
		ID:      DeserializeID(pbFu.Id),
		Name:    pbFu.GetName(),
		Pkg:     pbFu.GetPkgPath(),
		Step:    nil, // filled later
		Spread:  nil, // filled later
		Enter:   nil, // filled later
		Field:   nil, // filled later
		Monadic: pbFu.GetMonadic(),
		Leave:   nil, // filled later
		Syntax:  DeserializeSyntax(pbFu.Source),
	}
	unresolvedSteps, steps := deserializeUnresolvedSteps(pbFu.Step)
	fu.Step = steps
	fu.Enter = lookupStepLabelled(steps, pbFu.GetEnter())
	fu.Leave = lookupStepLabelled(steps, pbFu.GetLeave())
	fu.Field = map[string]*Step{}
	for _, arg := range pbFu.Arg {
		fu.Field[arg.GetName()] = lookupStepLabelled(steps, arg.GetStep())
	}
	// resolve Step.Gather
	for _, unresolvedStep := range unresolvedSteps {
		pbGathers := unresolvedStep.Proto.Gather
		unresolvedStep.Step.Gather = make([]*Gather, len(pbGathers))
		for i, pbGather := range pbGathers {
			unresolvedStep.Step.Gather[i] = &Gather{
				Field: pbGather.GetArg(),
				Step:  lookupStepLabelled(steps, pbGather.GetStep()),
			}
		}
	}
	// fill Func.Spread
	BacklinkFunc(fu)
	return fu
}

func lookupStepLabelled(steps []*Step, label string) *Step {
	for _, step := range steps {
		if step.Label == label {
			return step
		}
	}
	panic("step not found")
}

type unresolvedStep struct {
	Step  *Step    `ko:"name=step"`
	Proto *pb.Step `ko:"name=proto"`
}

func deserializeUnresolvedSteps(pbSteps []*pb.Step) ([]*unresolvedStep, []*Step) {
	unresolved := make([]*unresolvedStep, len(pbSteps))
	steps := make([]*Step, len(pbSteps))
	for i, pbStep := range pbSteps {
		steps[i] = &Step{
			ID:     DeserializeID(pbStep.Id),
			Label:  pbStep.GetLabel(),
			Logic:  DeserializeLogic(pbStep.Logic),
			Syntax: DeserializeSyntax(pbStep.Source),
			Func:   nil, // backlinks to func added in deserializeFunc
			Gather: nil, // backlinks to upstream steps added in deserializeFunc
		}
		unresolved[i] = &unresolvedStep{
			Proto: pbStep,
			Step:  steps[i],
		}
	}
	return unresolved, steps
}

func DeserializeID(pbID *pb.ID) ID {
	return IDFromProtoData(pbID.GetData())
}

func DeserializeSyntax(pbSource *pb.Source) Syntax {
	return &StartEndRegion{
		Start: Position{
			Filename: pbSource.GetFile(),
			Offset:   int(pbSource.Start.GetOffset()),
			Line:     int(pbSource.Start.GetLine()),
			Column:   int(pbSource.Start.GetColumn()),
		},
		End: Position{
			Filename: pbSource.GetFile(),
			Offset:   int(pbSource.End.GetOffset()),
			Line:     int(pbSource.End.GetLine()),
			Column:   int(pbSource.End.GetColumn()),
		},
	}
}

func DeserializeLogic(pbLogic *pb.Logic) Logic {
	switch u := pbLogic.Logic.(type) {
	case *pb.Logic_Enter:
		return Enter{}
	case *pb.Logic_Leave:
		return Leave{}
	case *pb.Logic_Select:
		return Select{Path: u.Select.Path}
	case *pb.Logic_Augment:
		return Augment{}
	case *pb.Logic_Invoke:
		return Invoke{}
	case *pb.Logic_Operator:
		return Operator{Path: u.Operator.Path}
	case *pb.Logic_PkgFunc:
		return PkgFunc{
			Pkg:  u.PkgFunc.GetPkg(),
			Func: u.PkgFunc.GetFunc(),
		}
	case *pb.Logic_Number:
		switch w := u.Number.Number.(type) {
		case *pb.LogicNumber_Bool:
			return Number{Value: w.Bool}
		case *pb.LogicNumber_String_:
			return Number{Value: LexString{String: w.String_}}
		case *pb.LogicNumber_Int64:
			return Number{Value: LexInteger{Int64: w.Int64}}
		case *pb.LogicNumber_Float64:
			return Number{Value: LexFloat{Float64: w.Float64}}
		}
	}
	panic("o")
}
