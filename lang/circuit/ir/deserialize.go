//
// Copyright © 2018 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ir

import (
	"github.com/kocircuit/kocircuit/lang/circuit/compile"
	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	"github.com/kocircuit/kocircuit/lang/circuit/lex"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func DeserializeRepo(pbRepo *pb.Repo) (repo model.Repo, err error) {
	defer func() {
		if r := recover(); r != nil {
			repo, err = nil, r.(error)
		}
	}()
	return deserializeRepo(pbRepo), nil
}

func deserializeRepo(pbRepo *pb.Repo) model.Repo {
	repo := model.Repo{}
	for _, pbPkg := range pbRepo.Package {
		pkgPath, pkg := deserializePackage(pbPkg)
		repo[pkgPath] = pkg
	}
	return repo
}

func deserializePackage(pbPkg *pb.Package) (pkgPath string, pkg model.Package) {
	pkgPath, pkg = pbPkg.GetPath(), model.Package{}
	for _, pbFu := range pbPkg.Func {
		fu := deserializeFunc(pbFu)
		pkg[pbFu.GetName()] = fu
	}
	return pkgPath, pkg
}

func deserializeFunc(pbFu *pb.Func) (fu *model.Func) {
	fu = &model.Func{
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
	fu.Field = map[string]*model.Step{}
	for _, arg := range pbFu.Arg {
		fu.Field[arg.GetName()] = lookupStepLabelled(steps, arg.GetStep())
	}
	// resolve Step.Gather
	for _, unresolvedStep := range unresolvedSteps {
		pbGathers := unresolvedStep.Proto.Gather
		unresolvedStep.Step.Gather = make([]*model.Gather, len(pbGathers))
		for i, pbGather := range pbGathers {
			unresolvedStep.Step.Gather[i] = &model.Gather{
				Field: pbGather.GetArg(),
				Step:  lookupStepLabelled(steps, pbGather.GetStep()),
			}
		}
	}
	// fill Func.Spread
	compile.BacklinkFunc(fu)
	return fu
}

func lookupStepLabelled(steps []*model.Step, label string) *model.Step {
	for _, step := range steps {
		if step.Label == label {
			return step
		}
	}
	panic("step not found")
}

type unresolvedStep struct {
	Step  *model.Step `ko:"name=step"`
	Proto *pb.Step    `ko:"name=proto"`
}

func deserializeUnresolvedSteps(pbSteps []*pb.Step) ([]*unresolvedStep, []*model.Step) {
	unresolved := make([]*unresolvedStep, len(pbSteps))
	steps := make([]*model.Step, len(pbSteps))
	for i, pbStep := range pbSteps {
		steps[i] = &model.Step{
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

func DeserializeID(pbID *pb.ID) model.ID {
	return model.IDFromProtoData(pbID.GetData())
}

func DeserializeSyntax(pbSource *pb.Source) syntax.Syntax {
	return &lex.StartEndRegion{
		Start: lex.Position{
			Filename: pbSource.GetFile(),
			Offset:   int(pbSource.Start.GetOffset()),
			Line:     int(pbSource.Start.GetLine()),
			Column:   int(pbSource.Start.GetColumn()),
		},
		End: lex.Position{
			Filename: pbSource.GetFile(),
			Offset:   int(pbSource.End.GetOffset()),
			Line:     int(pbSource.End.GetLine()),
			Column:   int(pbSource.End.GetColumn()),
		},
	}
}

func DeserializeLogic(pbLogic *pb.Logic) model.Logic {
	switch u := pbLogic.Logic.(type) {
	case *pb.Logic_Enter:
		return model.Enter{}
	case *pb.Logic_Leave:
		return model.Leave{}
	case *pb.Logic_Link:
		return model.Link{
			Name:    u.Link.GetName(),
			Monadic: u.Link.GetMonadic(),
		}
	case *pb.Logic_Select:
		return model.Select{Path: u.Select.Path}
	case *pb.Logic_Augment:
		return model.Augment{}
	case *pb.Logic_Invoke:
		return model.Invoke{}
	case *pb.Logic_Operator:
		return model.Operator{Path: u.Operator.Path}
	case *pb.Logic_PkgFunc:
		return model.PkgFunc{
			Pkg:  u.PkgFunc.GetPkg(),
			Func: u.PkgFunc.GetFunc(),
		}
	case *pb.Logic_Number:
		switch w := u.Number.Number.(type) {
		case *pb.LogicNumber_Bool:
			return model.Number{Value: w.Bool}
		case *pb.LogicNumber_String_:
			return model.Number{Value: lex.LexString{String: w.String_}}
		case *pb.LogicNumber_Int64:
			return model.Number{Value: lex.LexInteger{Int64: w.Int64}}
		case *pb.LogicNumber_Float64:
			return model.Number{Value: lex.LexFloat{Float64: w.Float64}}
		}
	}
	panic("o")
}
