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

func (env evalEnvelope) newFlow(span *Span, r Shape, eff Effect) evalFlow {
	return evalFlow{env: env, Frame: span, Shape: r, Effect: eff}
}

func (env evalEnvelope) returnFlow(span *Span, r Return, eff Effect) evalFlow {
	return evalFlow{env: env, Frame: span, Return: r, Effect: eff}
}

func (env evalEnvelope) Enter(span *Span) (Flow, error) {
	arg, effect, err := env.boundary().Enter(span, env.Arg)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

func (env evalEnvelope) Make(span *Span, v interface{}) (Flow, error) {
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
	arg, effect, err := env.boundary().Figure(span, w)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

func (env evalEnvelope) MakePkgFunc(span *Span, pkg string, name string) (Flow, error) {
	if fn := env.Program.Repo.Lookup(pkg, name); fn != nil { // user can overwrite idiomatic functions
		return env.makeMacroFigure(
			span,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if fn := env.Program.Idiom.Lookup(pkg, name); fn != nil { // otherwise hereditary idiom used
		return env.makeMacroFigure(
			span,
			env.Program.System.Combiner.Interpret(env.Program, fn),
		)
	}
	if macro := env.Program.Faculty[Ideal{Pkg: pkg, Name: name}]; macro != nil {
		return env.makeMacroFigure(span, macro)
	}
	return nil, span.Errorf(nil, "no function or macro %q.%s", pkg, name)
}

const IdiomRootPkg = "idiom"

func (env evalEnvelope) MakeOp(span *Span, ref []string) (Flow, error) {
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
				span,
				env.Program.System.Combiner.Interpret(env.Program, idiomFn),
			)
		}
	}
	ideal := Ideal{Name: strings.Join(ref, ".")}
	if macro := env.Program.Faculty[ideal]; macro != nil {
		return env.makeMacroFigure(span, macro)
	}
	return nil, span.Errorf(nil, "macro %v not known", ideal)
}

func (env evalEnvelope) makeMacroFigure(span *Span, macro Macro) (Flow, error) {
	arg, effect, err := env.boundary().Figure(span, macro)
	if err != nil {
		return nil, err
	}
	return env.newFlow(span, arg, effect), nil
}

type evalFlow struct {
	env    evalEnvelope
	Shape  Shape
	Return Arg
	Effect Effect
	Frame  *Span // span during when this flow was created
}

func (f evalFlow) Link(span *Span, name string, monadic bool) (Flow, error) {
	returns, effect, err := f.Shape.Link(span, name, monadic)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Select(span *Span, path []string) (Flow, error) {
	if len(path) == 0 {
		return f, nil
	}
	returns, effect, err := f.Shape.Select(span, Path(path))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Augment(span *Span, gather []GatherFlow) (Flow, error) {
	returns, effect, err := f.Shape.Augment(span, gatherFlowArg(gather))
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func gatherFlowArg(gather []GatherFlow) Fields {
	var s Fields
	for _, g := range gather {
		fieldFlow := g.Flow.(evalFlow)
		s = append(s, Field{
			Name:  g.Field,
			Shape: fieldFlow.Shape,
		})
	}
	return s
}

func (f evalFlow) Invoke(span *Span) (Flow, error) {
	returns, effect, err := f.Shape.Invoke(span)
	if err != nil {
		return nil, err
	}
	return f.env.newFlow(span, returns, effect), nil
}

func (f evalFlow) Leave(span *Span) (Flow, error) {
	r, effect, err := f.env.boundary().Leave(span, f.Shape)
	if err != nil {
		return nil, err
	}
	return f.env.returnFlow(span, r, effect), nil
}
