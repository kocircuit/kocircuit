package weave

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

type GoWeaving struct {
	Span    *Span       `ko:"name=span"`
	Func    *Func       `ko:"name=func"`
	Arg     GoStructure `ko:"name=arg"`
	Returns GoType      `ko:"name=returns"`
}

func (weaving *GoWeaving) GoCircuitName() string {
	return weaving.Valve().Address.Name
}

func (weaving *GoWeaving) Valve() *GoValve {
	return MakeValve(weaving.Span, weaving.Func, weaving.Arg, weaving.Returns)
}

type GoWeaveCtx struct {
	Label_       string       `ko:"name=label"`
	Idiom        Repo         `ko:"name=idiom"`
	Repo         Repo         `ko:"name=repo"`
	Faculty      Faculty      `ko:"name=faculty"`
	Parent       *GoWeaveCtx  `ko:"name=parent"`
	Weaving      *GoWeaving   `ko:"name=weaving"`
	AssignCache_ *AssignCache `ko:"name=assignCache"`
}

// var GoBackup = append(GoFaculty().PkgNames(), GoIdiomRepo.PkgNames()...)

func NewGoWeaveCtx(label string, repo Repo, faculty Faculty, idiom Repo) *GoWeaveCtx {
	return &GoWeaveCtx{
		Label_:       label,
		Idiom:        idiom,
		Repo:         repo,
		Faculty:      faculty,
		AssignCache_: nil, // no prior cache
	}
}

func (ctx *GoWeaveCtx) SheathID() *ID {
	return nil
}

func (ctx *GoWeaveCtx) SheathLabel() *string {
	return nil
}

func (ctx *GoWeaveCtx) SheathString() *string {
	return nil
}

func (ctx *GoWeaveCtx) Clear() *GoWeaveCtx {
	return &GoWeaveCtx{
		Label_:       ctx.Label_,
		Idiom:        ctx.Idiom,
		Repo:         ctx.Repo,
		Faculty:      ctx.Faculty,
		AssignCache_: ctx.AssignCache_,
		Parent:       nil, // disconnect weaving context
		Weaving:      nil,
	}
}

func (ctx *GoWeaveCtx) AssignCache() *AssignCache { return ctx.AssignCache_ }

func (ctx *GoWeaveCtx) UseCache(cache *AssignCache) *GoWeaveCtx {
	return &GoWeaveCtx{
		Label_:       ctx.Label_,
		Idiom:        ctx.Idiom,
		Repo:         ctx.Repo,
		Faculty:      ctx.Faculty,
		AssignCache_: cache,
		Parent:       ctx.Parent,
		Weaving:      ctx.Weaving,
	}
}

func (ctx *GoWeaveCtx) PushWeaving(weaving *GoWeaving) *GoWeaveCtx {
	return &GoWeaveCtx{
		Label_:       weaving.GoCircuitName(),
		Idiom:        ctx.Idiom,
		Repo:         ctx.Repo,
		Faculty:      ctx.Faculty,
		AssignCache_: ctx.AssignCache_,
		Parent:       ctx,
		Weaving:      weaving,
	}
}

func (ctx *GoWeaveCtx) LookupWeavingFor(f *Func) *GoWeaveCtx {
	for ctx != nil {
		if ctx.Weaving != nil && ctx.Weaving.Func == f {
			return ctx
		}
		ctx = ctx.Parent
	}
	return nil
}

func (ctx *GoWeaveCtx) Label() string {
	return ctx.Label_
}

func (ctx *GoWeaveCtx) Valve() *GoValve {
	return ctx.Weaving.Valve()
}

func (ctx *GoWeaveCtx) WeaveInstrument(span *Span, f *Func, arg GoStructure) (*GoInstrument, error) {
	if returns, effect, err := ctx.Weave(span, f, arg); err != nil {
		return nil, err
	} else {
		if len(effect.ProgramEffect.Recurrence) > 0 {
			panic("o")
		}
		return &GoInstrument{
			Valve:         effect.Valve,
			Returns:       returns,
			Circuit:       effect.ProgramEffect.Circuit,
			Directive:     effect.ProgramEffect.Directive,
			ProgramEffect: effect.ProgramEffect,
		}, nil
	}
}

// On entry, arg is a projection line.Real (GoNeverNilPtr{GoStruct{█}}).
func (ctx *GoWeaveCtx) Weave(span *Span, f *Func, arg GoStructure) (
	returns GoType,
	effect *GoCombineEffect,
	err error,
) {
	arg = RenameMonadicForFunc(span, arg, f)
	if priorCtx := ctx.LookupWeavingFor(f); priorCtx != nil { // function is already weaving
		valve := priorCtx.Weaving.Valve()
		return valve.Returns, &GoCombineEffect{
			Valve: valve, // valve with unknown return
			ProgramEffect: &GoProgramEffect{
				Recurrence: []*GoArgRecurrence{
					&GoArgRecurrence{Func: f, Arg: arg, Over: priorCtx},
				},
			},
		}, nil
	} else { // function is not currently weaving
		returns, effect, err = ctx.WeaveFixedPoint(span, f, arg)
		return
	}
}

type GoClamp struct {
	Arg     GoStructure `ko:"name=arg"`
	Returns GoType      `ko:"name=returns"`
}

func (ctx *GoWeaveCtx) WeaveFixedPoint(span *Span, f *Func, arg GoStructure) (
	returns GoType,
	effect *GoCombineEffect,
	err error,
) {
	stat := &GoWeavingStat{RecursionCount: 1}
	current := &GoClamp{Arg: arg, Returns: NewGoUnknown(span)}
	burning := false
	for {
		var returned GoType
		if returned, effect, err = ctx.evalIterate(current.Returns, span, f, current.Arg); err != nil {
			return nil, nil, err
		}
		stat.IterationCount++
		// generalizations
		next := &GoClamp{}
		argFixed, returnFixed := true, true
		// generalize argument occurrences
		recurrence, remainder := effect.SubtractRecurrences(f)
		// ensure non-self-recursive functions don't iterate to fixed point,
		// this would introduce exponential dependence on count of non-self-recursive circuits.
		if len(recurrence) == 0 {
			return returned, remainder.AggregateWeavingStat(stat), nil
		}
		next.Arg = current.Arg
		for _, re := range recurrence {
			if next.Arg, err = GeneralizeStructure(span, next.Arg, re.Arg); err != nil {
				return nil, nil, span.Errorf(err, "cannot generalize arguments %s", Sprint(re))
			}
		}
		if _, _, err := Assign(span, next.Arg, current.Arg); err != nil {
			argFixed = false
		}
		// generalize return values
		if next.Returns, err = Generalize(span, current.Returns, returned); err != nil {
			return nil, nil, span.Errorf(err, "cannot generalize returns %s", Sprint(returned))
		} else {
			if _, _, err := Assign(span, next.Returns, current.Returns); err != nil {
				returnFixed = false
			}
		}
		// exit logic
		if burning {
			if argFixed && returnFixed {
				return returned, remainder.AggregateWeavingStat(stat), nil
			} else {
				burning = false
			}
		} else {
			if argFixed && returnFixed {
				burning = true
			}
		}
		current = next
	}
	panic("o")
}

func (ctx *GoWeaveCtx) evalIterate(priorReturns GoType, span *Span, f *Func, iterate GoStructure) (
	returns GoType,
	effect *GoCombineEffect,
	err error,
) {
	weaving := &GoWeaving{Span: span, Func: f, Arg: iterate, Returns: priorReturns}
	undeterminedArg := weaving.Valve().Real().(GoStructure)
	iterationSpan := RefineWeaveCtx(span, ctx.PushWeaving(weaving))
	if returns, effect, err = ctx.Eval(iterationSpan, f, undeterminedArg); err != nil {
		return nil, nil, err
	}
	return returns, effect, nil
}

// span holds assignment cache passed to evaluator.
// Returned effect holds the newly cached assignments.
func (ctx *GoWeaveCtx) Eval(span *Span, f *Func, arg GoStructure) (GoType, *GoCombineEffect, error) {
	prog := Program{
		Idiom: ctx.Idiom,
		Repo:  ctx.Repo,
		System: System{
			Faculty:  ctx.Faculty, // macros for registrations from RegisterMacro
			Boundary: GoBoundary{},
			Combiner: GoCombiner{},
		},
	}
	shape, effect, err := prog.EvalSeq(span, f, arg)
	if err != nil {
		return nil, nil, err
	}
	return shape.(GoType), effect.(*GoCombineEffect), nil
}

type GoWeavingStat struct {
	IterationCount int `ko:"name=iterationCount"`
	RecursionCount int `ko:"name=recursionCount"`
}

func SumGoWeavingStat(x, y *GoWeavingStat) *GoWeavingStat {
	if x == nil {
		x = &GoWeavingStat{}
	}
	if y == nil {
		y = &GoWeavingStat{}
	}
	return &GoWeavingStat{
		IterationCount: x.IterationCount + y.IterationCount,
		RecursionCount: x.RecursionCount + y.RecursionCount,
	}
}
