package eval

import (
	"strings"

	. "github.com/kocircuit/kocircuit/lang/circuit/flow"
	. "github.com/kocircuit/kocircuit/lang/circuit/lex"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

type evalEnvelope struct {
	Program Program
	Arg     Arg
	Span    *Span
}

func (env evalEnvelope) boundary() Boundary { return env.Program.System.Boundary }

func (env evalEnvelope) newFlow(frame *Span, r Shape, eff Effect) evalFlow {
	return evalFlow{env: env, Frame: frame, Shape: r, Effect: eff}
}

func (env evalEnvelope) returnFlow(frame *Span, r Return, eff Effect) evalFlow {
	return evalFlow{env: env, Frame: frame, Return: r, Effect: eff}
}

func (env evalEnvelope) Enter(frame *Span) (Flow, error) {
	arg, effect, err := env.boundary().Enter(frame, env.Arg)
	if err != nil {
		return nil, err
	}
	return env.newFlow(frame, arg, effect), nil
}

func (env evalEnvelope) Make(frame *Span, v interface{}) (Flow, error) {
	var w Figure
	switch u := v.(type) {
	case LexString:
		w = String{u.String}
	case bool: // comes from graftTerm
		w = Bool{u}
	case LexInteger:
		w = Integer{u.Int64}
	case LexFloat:
		w = Float{u.Float64}
	default:
		panic("unrecognized figure")
	}
	arg, effect, err := env.boundary().Figure(frame, w)
	if err != nil {
		return nil, err
	}
	return env.newFlow(frame, arg, effect), nil
}

func (env evalEnvelope) MakePkgFunc(frame *Span, pkg string, name string) (Flow, error) {
	if fn := env.Program.Repo.Lookup(pkg, name); fn != nil { // user can overwrite idiomatic functions
		return env.makeMacroFigure(
			frame,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if fn := env.Program.Idiom.Lookup(pkg, name); fn != nil { // otherwise hereditary idiom used
		return env.makeMacroFigure(
			frame,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if macro := env.Program.Faculty[Ideal{Pkg: pkg, Name: name}]; macro != nil {
		return env.makeMacroFigure(frame, macro)
	}
	return nil, frame.Errorf(nil, "no function or macro %q.%s", pkg, name)
}

const IdiomRootPkg = "idiom"

func (env evalEnvelope) MakeOp(frame *Span, ref []string) (Flow, error) {
	// first check if idiomatic repo implements the operator
	// thus, idiomatic circuits overwrite hard-coded macros
	if len(ref) > 0 {
		idiomFu := ref[len(ref)-1]
		idiomPkg := strings.Join(
			append([]string{IdiomRootPkg}, ref[:len(ref)-1]...),
			".",
		)
		if idiomFn := env.Program.Idiom.Lookup(idiomPkg, idiomFu); idiomFn != nil {
			return env.makeMacroFigure(
				frame,
				env.Program.System.Combiner.Interpret(env.Program, idiomFn),
			)
		}
	}
	ideal := Ideal{Name: strings.Join(ref, ".")}
	if macro := env.Program.Faculty[ideal]; macro != nil {
		return env.makeMacroFigure(frame, macro)
	}
	return nil, frame.Errorf(nil, "macro %v not known", ideal)
}

func (env evalEnvelope) makeMacroFigure(frame *Span, macro Macro) (Flow, error) {
	arg, effect, err := env.boundary().Figure(frame, macro)
	if err != nil {
		return nil, err
	}
	return env.newFlow(frame, arg, effect), nil
}

type evalFlow struct {
	env    evalEnvelope
	Shape  Shape
	Return Arg
	Effect Effect
	Frame  *Span // frame during when this flow was created
}

func (f evalFlow) Select(frame *Span, path []string) (Flow, error) {
	if len(path) == 0 {
		return f, nil
	}
	returns, effect, err := f.Shape.Select(frame, Path(path))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(frame, returns, effect), nil
}

func (f evalFlow) Augment(frame *Span, gather []GatherFlow) (Flow, error) {
	returns, effect, err := f.Shape.Augment(frame, gatherFlowArg(gather))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(frame, returns, effect), nil
}

func gatherFlowArg(gather []GatherFlow) Knot {
	var s Knot
	for _, g := range gather {
		fieldFlow := g.Flow.(evalFlow)
		s = append(s, Field{
			Name:   g.Field,
			Shape:  fieldFlow.Shape,
			Effect: fieldFlow.Effect,
			Frame:  fieldFlow.Frame,
		})
	}
	return s
}

func (f evalFlow) Invoke(frame *Span) (Flow, error) {
	returns, effect, err := f.Shape.Invoke(frame)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(frame, returns, effect), nil
}

func (f evalFlow) Leave(frame *Span) (Flow, error) {
	r, effect, err := f.env.boundary().Leave(frame, f.Shape)
	if err != nil {
		return nil, err
	}
	return f.env.returnFlow(frame, r, effect), nil
}
