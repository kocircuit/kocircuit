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
	"github.com/golang/protobuf/proto"

	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	"github.com/kocircuit/kocircuit/lang/circuit/lex"
	"github.com/kocircuit/kocircuit/lang/circuit/model"
	"github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func SerializeRepo(repo model.Repo) *pb.Repo {
	pbRepo := &pb.Repo{
		Package: make([]*pb.Package, 0, len(repo)),
	}
	for pkgPath, pkg := range repo {
		pbRepo.Package = append(pbRepo.Package, serializePackage(pkgPath, pkg))
	}
	return pbRepo
}

func serializePackage(pkgPath string, pkg model.Package) *pb.Package {
	pbPackage := &pb.Package{
		Path: proto.String(pkgPath),
		Func: make([]*pb.Func, 0, len(pkg)),
	}
	for _, fu := range pkg {
		pbPackage.Func = append(pbPackage.Func, serializeFunc(fu))
	}
	return pbPackage
}

func serializeFunc(fu *model.Func) *pb.Func {
	pbFunc := &pb.Func{
		Doc:     proto.String(fu.Doc),
		Id:      SerializeID(fu.ID),
		Name:    proto.String(fu.Name),
		PkgPath: proto.String(fu.Pkg),
		Step:    serializeSteps(fu.Step),
		Enter:   proto.String(fu.Enter.Label),
		Leave:   proto.String(fu.Leave.Label),
		Arg:     serializeArg(fu.Field),
		Monadic: proto.String(fu.Monadic),
		Source:  SerializeSyntax(fu.Syntax),
	}
	return pbFunc
}

func serializeSteps(steps []*model.Step) []*pb.Step {
	pbSteps := make([]*pb.Step, len(steps))
	for i, step := range steps {
		pbSteps[i] = &pb.Step{
			Id:      SerializeID(step.ID),
			Label:   proto.String(step.Label),
			Gather:  serializeGather(step.Gather),
			Logic:   SerializeLogic(step.Logic),
			Source:  SerializeSyntax(step.Syntax),
			PkgPath: proto.String(step.Func.Pkg),
			Func:    proto.String(step.Func.Name),
		}
	}
	return pbSteps
}

func serializeGather(gather []*model.Gather) []*pb.Gather {
	pbGather := make([]*pb.Gather, len(gather))
	for i, g := range gather {
		pbGather[i] = &pb.Gather{Arg: proto.String(g.Field), Step: proto.String(g.Step.Label)}
	}
	return pbGather
}

func serializeArg(arg map[string]*model.Step) []*pb.Arg {
	pbArg := make([]*pb.Arg, 0, len(arg))
	for name, step := range arg {
		pbArg = append(pbArg, &pb.Arg{Name: proto.String(name), Step: proto.String(step.Label)})
	}
	return pbArg
}

func SerializeID(id model.ID) *pb.ID {
	return &pb.ID{Data: proto.Uint64(id.ProtoData())}
}

func SerializeSyntax(syntax syntax.Syntax) *pb.Source {
	return &pb.Source{
		File:  proto.String(syntax.FilePath()),
		Start: SerializePosition(syntax.StartPosition()),
		End:   SerializePosition(syntax.EndPosition()),
	}
}

func SerializePosition(pos lex.Position) *pb.Position {
	return &pb.Position{
		Offset: proto.Int64(int64(pos.Offset)),
		Line:   proto.Int64(int64(pos.Line)),
		Column: proto.Int64(int64(pos.Column)),
	}
}

func SerializeLogic(logic model.Logic) *pb.Logic {
	switch u := logic.(type) {
	case model.Operator:
		return &pb.Logic{
			Logic: &pb.Logic_Operator{
				Operator: &pb.LogicOperator{
					Path: u.Path,
				},
			},
		}
	case model.PkgFunc:
		return &pb.Logic{
			Logic: &pb.Logic_PkgFunc{
				PkgFunc: &pb.LogicPkgFunc{
					Pkg:  proto.String(u.Pkg),
					Func: proto.String(u.Func),
				},
			},
		}
	case model.Enter:
		return &pb.Logic{
			Logic: &pb.Logic_Enter{
				Enter: &pb.LogicEnter{},
			},
		}
	case model.Leave:
		return &pb.Logic{
			Logic: &pb.Logic_Leave{
				Leave: &pb.LogicLeave{},
			},
		}
	case model.Link:
		return &pb.Logic{
			Logic: &pb.Logic_Link{
				Link: &pb.LogicLink{
					Name:    proto.String(u.Name),
					Monadic: proto.Bool(u.Monadic),
				},
			},
		}
	case model.Select:
		return &pb.Logic{
			Logic: &pb.Logic_Select{
				Select: &pb.LogicSelect{
					Path: u.Path,
				},
			},
		}
	case model.Augment:
		return &pb.Logic{
			Logic: &pb.Logic_Augment{
				Augment: &pb.LogicAugment{},
			},
		}
	case model.Invoke:
		return &pb.Logic{
			Logic: &pb.Logic_Invoke{
				Invoke: &pb.LogicInvoke{},
			},
		}
	case model.Number:
		return SerializeNumberLogic(u)
	}
	panic("o")
}

func SerializeNumberLogic(n model.Number) *pb.Logic {
	switch u := n.Value.(type) {
	case bool:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Bool{Bool: u},
				},
			},
		}
	case lex.LexInteger:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Int64{Int64: u.Int64},
				},
			},
		}
	case lex.LexString:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_String_{String_: u.String},
				},
			},
		}
	case lex.LexFloat:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Float64{Float64: u.Float64},
				},
			},
		}
	}
	panic("o")
}
