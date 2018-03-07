package weave

import (
	"fmt"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	. "github.com/kocircuit/kocircuit/lang/go/kit/tree"
	. "github.com/kocircuit/kocircuit/lang/go/kit/util"
	. "github.com/kocircuit/kocircuit/lang/go/model"
)

func init() {
	RegisterGoMacro("Loop", new(GoLoopMacro))
}

// Loop(start: start_value, step: step_func, stop: optional_stop_condition)
type GoLoopMacro struct{}

func (m GoLoopMacro) MacroID() string { return m.Help() }

func (m GoLoopMacro) Label() string { return "loop" }

func (m GoLoopMacro) MacroSheathString() *string { return PtrString("Loop") }

func (m GoLoopMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

// Emitted effect is GoMacroEffect.
func (GoLoopMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveLoop(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, span.Errorf(err, "solving loop")
	}
	return soln.Returns(), SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type LoopSoln struct {
	Origin *Span        `ko:"name=origin"`
	Arg    *LoopSolnArg `ko:"name=arg"`
	Iter   GoType       `ko:"name=iter"`
	Start  struct {
		Start       GoType `ko:"name=start"`
		StartToIter Shaper `ko:"name=startToIter"`
	} `ko:"name=start"`
	Step struct {
		Evo              *Evocation `ko:"name=evocation"`
		StepResultToIter Shaper     `ko:"name=stepResultToIter"`
	} `ko:"name=step"`
	Bootstrap struct { // bootstrap step (first step invocation, when start not given)
		Inv                       *Invocation `ko:"name=invocation"`
		BootstrapStepResultToIter Shaper      `ko:"name=bootstrapStepResultToIter"`
	} `ko:"name=bootstrap"`
	Stop struct {
		Evo              *Evocation `ko:"name=evocation"`
		StopResultToBool Shaper     `ko:"name=stopResultToBool"`
	} `ko:"name=stop"`
	Cached_ *AssignCache `ko:"name=cached"`
}

func (soln *LoopSoln) String() string { return Sprint(soln) }

func (soln *LoopSoln) Cached() *AssignCache {
	return AssignCacheUnion(
		soln.Cached_,
		soln.Step.Evo.Cached(),
		soln.Bootstrap.Inv.Cached(),
		soln.Stop.Evo.Cached(),
	)
}

func (soln *LoopSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		soln.Arg.CircuitEffect(),
		CircuitEffectIfNotNil(soln.Start.StartToIter),
		CircuitEffectIfNotNil(soln.Bootstrap.Inv),
		CircuitEffectIfNotNil(soln.Bootstrap.BootstrapStepResultToIter),
		soln.Step.Evo.CircuitEffect(),
		soln.Step.StepResultToIter.CircuitEffect(),
		soln.Stop.Evo.CircuitEffect(),
		CircuitEffectIfNotNil(soln.Stop.StopResultToBool),
	).AggregateDuctFunc(soln.FuncExpr())
}

func (soln *LoopSoln) ProgramEffect() *GoProgramEffect {
	return AggregateProgramEffects(
		ProgramEffectIfNotNil(soln.Bootstrap.Inv),
		ProgramEffectIfNotNil(soln.Step.Evo),
		ProgramEffectIfNotNil(soln.Stop.Evo),
	)
}

type LoopSolnArg struct {
	Arg   GoStructure       `ko:"name=arg"`
	Start TypeExtractor     `ko:"name=start"`
	Step  VarietalExtractor `ko:"name=step"`
	Stop  TypeExtractor     `ko:"name=stop"`
}

func (arg *LoopSolnArg) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		arg.Start.Extractor.CircuitEffect(),
		arg.Step.Extractor.CircuitEffect(),
		arg.Stop.Extractor.CircuitEffect(),
	)
}

func SolveLoopArg(span *Span, arg GoStructure) (soln *LoopSolnArg, err error) {
	soln = &LoopSolnArg{Arg: arg}
	if soln.Start.Extractor, soln.Start.Type, err = GoSelect(span, Path{"start"}, arg); err != nil {
		return nil, span.Errorf(err, "loop expects a “start” argument")
	}
	if soln.Step.Extractor, soln.Step.Varietal, err = GoSelectVarietal(span, Path{"step"}, arg); err != nil {
		return nil, span.Errorf(err, "loop expects a “step” argument")
	}
	if soln.Stop.Extractor, soln.Stop.Type, err = GoSelectSimplify(span, Path{"stop"}, arg); err != nil {
		return nil, span.Errorf(err, "loop expects a “stop” argument")
	}
	return soln, nil
}

// Return condition on iter: Generalize(iter, Step(iter)) assigns to iter.
func SolveLoop(span *Span, argStruct GoStructure) (soln *LoopSoln, err error) {
	arg, err := SolveLoopArg(span, argStruct)
	if err != nil {
		return nil, err
	}
	soln = &LoopSoln{Origin: span, Arg: arg, Iter: arg.Start.Type}
	span = SpanClearWeavingCtx(span) // clear weaving recursion history
	assn := NewAssignCtx(span)
	for { // solution power-iteration loop
		if soln.Step.Evo, err = GoEvoke(span, arg.Step.Varietal, soln.IterField()); err != nil {
			return nil, span.Errorf(err, "loop step")
		}
		if unified, err := Generalize(span, soln.Iter, soln.Step.Evo.Returns()); err != nil {
			return nil, span.Errorf(err, "loop iterator generalization")
		} else {
			if _, err = assn.Assign(unified, soln.Iter); err != nil { // if unification continues to generalize/expand, continue iterating
				soln.Iter = unified
			} else { // otherwise, iter has stabilized at fixed point
				if _, startIsEmpty := arg.Start.Type.(*GoEmpty); !startIsEmpty {
					soln.Start.Start = arg.Start.Type
					if soln.Start.StartToIter, err = assn.Assign(arg.Start.Type, soln.Iter); err != nil {
						return nil, span.Errorf(err, "loop assigning start to iterator")
					}
				} else {
					if soln.Bootstrap.Inv, err = GoInvoke(span, arg.Step.Varietal); err != nil {
						return nil, span.Errorf(err, "loop invoking step without a start argument")
					}
					if soln.Bootstrap.BootstrapStepResultToIter, err = assn.Assign(soln.Bootstrap.Inv.Returns, soln.Iter); err != nil {
						return nil, span.Errorf(err, "loop shaping first step result to iterator")
					}
				}
				if soln.Step.StepResultToIter, err = assn.Assign(soln.Step.Evo.Returns(), soln.Iter); err != nil {
					return nil, span.Errorf(err, "loop step")
				}
				if err = solveLoopStop(assn, span, soln); err != nil {
					return nil, err
				}
				soln.Cached_ = assn.Flush()
				return soln, nil
			}
		}
	} // for (solve) loop
	panic("u")
}

const LoopIterArgName = "iter"

func (soln *LoopSoln) IterField() []*GoAugmentField {
	return WithGoField(
		&GoField{
			Type: soln.Iter,
			Name: GoNameFor(LoopIterArgName),
			Tag:  KoTags(LoopIterArgName, true),
		},
	)
}

func solveLoopStop(assn *AssignCtx, span *Span, soln *LoopSoln) (err error) {
	switch stop := soln.Arg.Stop.Type.(type) {
	case *GoEmpty:
		return nil
	case GoVarietal:
		if soln.Stop.Evo, err = GoEvoke(span, stop, soln.IterField()); err != nil {
			return span.Errorf(err, "loop stop")
		}
		if soln.Stop.StopResultToBool, err = assn.Assign(
			soln.Stop.Evo.Invocation.Returns,
			GoBool,
		); err != nil {
			return span.Errorf(err, "loop stop returns")
		}
		return nil
	default:
		return span.Errorf(nil, "loop stop argument must be a varietal, got %s", Sprint(stop))
	}
	panic("o")
}

func (soln *LoopSoln) Returns() GoType { return soln.Iter }

func (soln *LoopSoln) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: soln.Origin,
		Name:   fmt.Sprintf("loop_%s", soln.Origin.SpanID().String()),
	}
}

// FormExpr returns an invocation of the loop_ function with (step_ctx, LoopArgStructure)
func (soln *LoopSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: soln.FuncExpr(),
		Arg: []GoExpr{
			&GoVerbatimExpr{"step_ctx"},
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}

// FuncExpr returns the loop_ function definition:
//
// func loop_xyz(step_ctx *runtime.Context, step_arg LoopArg) (iter Iter) {
// //--- if start given
// 	iter = StartToIter(ExtractStart(step_arg))
// //--- else
// 	bootstrapStepResult := INVOKE(SelectStep(step_arg))
// 	_ = bootstrapStepResult
// 	iter = BootstrapStepResultToIter(bootstrapStepResult)
// //--- endif
// 	for {
// 		stepVty := AUGMENT(SelectStep(step_arg), iter)
// 		stepResult := INVOKE(stepVty)
// 		iter = StepResultToIter(stepResult)
// 		//--- if stop given
// 		stopResult := EVOKE(SelectStop(step_arg), iter)
// 		_ = stopResult
// 		if StopResultToBool(stopResult) {
// 			break
// 		}
// 		//--- endif
// 	}
// 	return
// }
func (soln *LoopSoln) FuncExpr() GoFuncExpr {
	return &GoDuctFuncExpr{
		Comment:    fmt.Sprintf("(%s)", soln.Origin.SourceLine()),
		FuncName:   soln.FuncName(),
		ArgName:    soln.argExpr(),
		ArgType:    soln.Arg.Arg,
		ReturnName: soln.iterExpr(),
		ReturnType: soln.Returns(),
		Line: MergeLine(
			soln.bootstrapLogic(),
			[]GoExpr{
				soln.forLogic(),
				&GoReturnExpr{EmptyExpr},
			},
		),
	}
}

func (soln *LoopSoln) iterExpr() GoExpr {
	return &GoVerbatimExpr{"iter"}
}

func (soln *LoopSoln) argExpr() GoExpr {
	return &GoVerbatimExpr{"step_arg"}
}

func (soln *LoopSoln) startExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.Start.Extractor, Expr: soln.argExpr()}
}

func (soln *LoopSoln) stepExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.Step.Extractor, Expr: soln.argExpr()}
}

func (soln *LoopSoln) stopExpr() GoExpr {
	return &GoShapeExpr{Shaper: soln.Arg.Stop.Extractor, Expr: soln.argExpr()}
}

// 	// --- if start given
// 	iter = shape_start_iter(arg.Field_start)
// 	// --- else
// 	bootstrapStepResult := INVOKE(step_ctx, arg.Field_step)
//		_ = bootstrapStepResult
// 	iter = bootstrapStepResultToIter(bootstrapStepResult)
// 	// --- endif
func (soln *LoopSoln) bootstrapLogic() []GoExpr {
	bootstrapStepResult := &GoVerbatimExpr{"bootstrapStepResult"}
	if soln.Start.Start != nil {
		return []GoExpr{
			&GoAssignExpr{
				Left: soln.iterExpr(),
				Right: &GoShapeExpr{
					Shaper: soln.Start.StartToIter,
					Expr:   soln.startExpr(),
				},
			},
		}
	} else {
		return []GoExpr{
			&GoColonAssignExpr{
				Left: bootstrapStepResult,
				Right: &GoSlotFormExpr{
					SlotExpr: []*GoSlotExpr{{Slot: RootSlot{}, Expr: soln.stepExpr()}},
					Form:     soln.Bootstrap.Inv,
				},
			},
			&GoAssignExpr{
				Left:  UnderlineExpr,
				Right: bootstrapStepResult,
			},
			&GoAssignExpr{
				Left: soln.iterExpr(),
				Right: &GoShapeExpr{
					Shaper: soln.Bootstrap.BootstrapStepResultToIter,
					Expr:   bootstrapStepResult,
				},
			},
		}
	}
}

// 	// --- if stop given
// 	stopVty := AUGMENT(arg.Field_stop, iter)
// 	stopResult := INVOKE(step_ctx, stopVty)
// 	if stopResultToBool(stopResult) {
// 		break
// 	}
// 	// --- endif
func (soln *LoopSoln) stopLogic() []GoExpr {
	if soln.Stop.Evo == nil {
		return nil
	}
	stopResult := &GoVerbatimExpr{"stopResult"}
	return []GoExpr{
		&GoColonAssignExpr{
			Left: stopResult,
			Right: &GoSlotFormExpr{
				SlotExpr: []*GoSlotExpr{
					{Slot: RootSlot{}, Expr: soln.stopExpr()},
					{Slot: NameSlot{LoopIterArgName}, Expr: soln.iterExpr()},
				},
				Form: soln.Stop.Evo,
			},
		},
		&GoAssignExpr{
			Left:  UnderlineExpr,
			Right: stopResult,
		},
		&GoIfThenExpr{
			If: &GoShapeExpr{
				Shaper: soln.Stop.StopResultToBool,
				Expr:   stopResult,
			},
			Then: []GoExpr{BreakExpr},
		},
	}
}

// 	for {
// 		stepVty := AUGMENT(arg.Field_step, iter)
// 		stepResult := INVOKE(step_ctx, stepVty)
// 		iter = stepResultToIter(stepResult)
//		// --- if stop given
// 		...
// 		// --- endif
// 	}
func (soln *LoopSoln) forLogic() GoExpr {
	stepResult := &GoVerbatimExpr{"stepResult"}
	return &GoForExpr{
		Line: MergeLine(
			[]GoExpr{
				&GoColonAssignExpr{
					Left: stepResult,
					Right: &GoSlotFormExpr{
						SlotExpr: []*GoSlotExpr{
							{Slot: RootSlot{}, Expr: soln.stepExpr()},
							{Slot: NameSlot{LoopIterArgName}, Expr: soln.iterExpr()},
						},
						Form: soln.Step.Evo,
					},
				},
				&GoAssignExpr{
					Left:  UnderlineExpr,
					Right: stepResult,
				},
				&GoAssignExpr{
					Left: soln.iterExpr(),
					Right: &GoShapeExpr{
						Shaper: soln.Step.StepResultToIter,
						Expr:   stepResult,
					},
				},
			},
			soln.stopLogic(),
		),
	}
}
