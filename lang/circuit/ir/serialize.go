package ir

import (
	"github.com/golang/protobuf/proto"

	pb "github.com/kocircuit/kocircuit/lang/circuit/ir/proto"
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/circuit/syntax"
)

func SerializeRepo(repo Repo) *pb.Repo {
	ctx := &serializingCtx{}
	return ctx.serializeRepo(repo)
}

type serializingCtx struct{}

func (ctx *serializingCtx) serializeRepo(repo Repo) *pb.Repo {
	pbRepo := &pb.Repo{
		Package: make([]*pb.Package, 0, len(repo)),
	}
	for pkgPath, pkg := range repo {
		pbRepo.Package = append(pbRepo.Package, ctx.serializePackage(pkgPath, pkg))
	}
	return pbRepo
}

func (ctx *serializingCtx) serializePackage(pkgPath string, pkg Package) *pb.Package {
	pbPackage := &pb.Package{
		Path: proto.String(pkgPath),
		Func: make([]*pb.Func, 0, len(pkg)),
	}
	for _, fu := range pkg {
		pbPackage.Func = append(pbPackage.Func, ctx.serializeFunc(fu))
	}
	return pbPackage
}

func (ctx *serializingCtx) serializeFunc(fu *Func) *pb.Func {
	pbFunc := &pb.Func{
		Id:      SerializeID(fu.ID),
		Name:    proto.String(fu.Name),
		PkgPath: proto.String(fu.Pkg),
		Step:    ctx.serializeSteps(fu.Step),
		Enter:   proto.String(fu.Enter.Label),
		Leave:   proto.String(fu.Leave.Label),
		Arg:     ctx.serializeArg(fu.Field),
		Monadic: proto.String(fu.Monadic),
		Source:  SerializeSyntax(fu.Syntax),
	}
	return pbFunc
}

func (ctx *serializingCtx) serializeSteps(steps []*Step) []*pb.Step {
	pbSteps := make([]*pb.Step, len(steps))
	for i, step := range steps {
		pbSteps[i] = &pb.Step{
			Id:      SerializeID(step.ID),
			Label:   proto.String(step.Label),
			Gather:  ctx.serializeGather(step.Gather),
			Logic:   SerializeLogic(step.Logic),
			Source:  SerializeSyntax(step.Syntax),
			PkgPath: proto.String(step.Func.Pkg),
			Func:    proto.String(step.Func.Name),
		}
	}
	return pbSteps
}

func (ctx *serializingCtx) serializeGather(gather []*Gather) []*pb.Gather {
	pbGather := make([]*pb.Gather, len(gather))
	for i, g := range gather {
		pbGather[i] = &pb.Gather{Arg: proto.String(g.Field), Step: proto.String(g.Step.Label)}
	}
	return pbGather
}

func (ctx *serializingCtx) serializeArg(arg map[string]*Step) []*pb.Arg {
	pbArg := make([]*pb.Arg, 0, len(arg))
	for name, step := range arg {
		pbArg = append(pbArg, &pb.Arg{Name: proto.String(name), Step: proto.String(step.Label)})
	}
	return pbArg
}

func SerializeID(id ID) *pb.ID {
	return &pb.ID{Data: proto.Uint64(id.ProtoData())}
}

func SerializeSyntax(syntax Syntax) *pb.Source {
	return &pb.Source{
		File:  proto.String(syntax.FilePath()),
		Start: SerializePosition(syntax.StartPosition()),
		End:   SerializePosition(syntax.EndPosition()),
	}
}

func SerializePosition(pos Position) *pb.Position {
	return &pb.Position{
		Offset: proto.Int64(int64(pos.Offset)),
		Line:   proto.Int64(int64(pos.Line)),
		Column: proto.Int64(int64(pos.Column)),
	}
}

func SerializeLogic(logic Logic) *pb.Logic {
	switch u := logic.(type) {
	case Operator:
		return &pb.Logic{
			Logic: &pb.Logic_Operator{
				Operator: &pb.LogicOperator{
					Path: u.Path,
				},
			},
		}
	case PkgFunc:
		return &pb.Logic{
			Logic: &pb.Logic_PkgFunc{
				PkgFunc: &pb.LogicPkgFunc{
					Pkg:  proto.String(u.Pkg),
					Func: proto.String(u.Func),
				},
			},
		}
	case Enter:
		return &pb.Logic{
			Logic: &pb.Logic_Enter{
				Enter: &pb.LogicEnter{},
			},
		}
	case Leave:
		return &pb.Logic{
			Logic: &pb.Logic_Leave{
				Leave: &pb.LogicLeave{},
			},
		}
	case Select:
		return &pb.Logic{
			Logic: &pb.Logic_Select{
				Select: &pb.LogicSelect{
					Path: u.Path,
				},
			},
		}
	case Augment:
		return &pb.Logic{
			Logic: &pb.Logic_Augment{
				Augment: &pb.LogicAugment{},
			},
		}
	case Invoke:
		return &pb.Logic{
			Logic: &pb.Logic_Invoke{
				Invoke: &pb.LogicInvoke{},
			},
		}
	case Number:
		return SerializeNumberLogic(u)
	}
	panic("o")
}

func SerializeNumberLogic(n Number) *pb.Logic {
	switch u := n.Value.(type) {
	case bool:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Bool{Bool: u},
				},
			},
		}
	case int64:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Int64{Int64: u},
				},
			},
		}
	case string:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_String_{String_: u},
				},
			},
		}
	case float64:
		return &pb.Logic{
			Logic: &pb.Logic_Number{
				Number: &pb.LogicNumber{
					Number: &pb.LogicNumber_Float64{Float64: u},
				},
			},
		}
	}
	panic("o")
}
