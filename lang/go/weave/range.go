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
	RegisterGoMacro("Range", new(GoRangeMacro))
}

// GoRangeMacro implements the Range macro, used like so:
// Main() {
// 	return: Range(
// 		over: (1, 2, 3, 4)
// 		with: step(carry, elem) {
// 			return: (
// 				emit: Increment(elem),
// 				carry: Sum(carry, elem)
// 			)
// 		}
//		) // (image: (█), residue: █)
// }
type GoRangeMacro struct{}

func (m GoRangeMacro) MacroID() string { return m.Help() }

func (m GoRangeMacro) Label() string { return "range" }

func (m GoRangeMacro) MacroSheathString() *string { return PtrString("Range") }

func (m GoRangeMacro) Help() string {
	return GoInterfaceTypeAddress(m).String()
}

func (GoRangeMacro) Invoke(span *Span, arg Arg) (returns Return, effect Effect, err error) {
	soln, err := SolveRange(span, arg.(GoStructure))
	if err != nil {
		return nil, nil, err
	}
	switch {
	case soln.OverIsUnknown():
		return GoUnknownMacro{}.Invoke(span, arg)
	case soln.OverIsEmpty():
		empty := NewGoEmpty(span)
		return empty,
			&GoMacroEffect{
				Arg:      arg.(GoType),
				SlotForm: &GoZeroForm{empty},
				Cached:   soln.Cached(),
			}, nil
	}
	return soln.Returns(), SlotFormMacroEffect(arg.(GoStructure), soln), nil
}

type RangeSolnArg struct {
	Arg           GoStructure       `ko:"name=arg"`
	Over          TypeExtractor     `ko:"name=over"`
	OverIsUnknown bool              `ko:"name=over_is_unknown"`
	OverIsEmpty   bool              `ko:"name=over_is_empty"`
	OverLift      Shaper            `ko:"name=over_lift"` // lift over to slice if necessary
	OverLifted    GoSequence        `ko:"name=over_lifted"`
	With          VarietalExtractor `ko:"name=with"`
	Start         TypeExtractor     `ko:"name=start"`
}

func (arg *RangeSolnArg) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		arg.Over.Extractor.CircuitEffect(),
		arg.OverLift.CircuitEffect(),
		arg.With.Extractor.CircuitEffect(),
		arg.Start.Extractor.CircuitEffect(),
	)
}

type RangeSoln struct {
	Origin       *Span            `ko:"name=origin"`
	Arg          *RangeSolnArg    `ko:"name=arg"`
	Step         *Evocation       `ko:"name=step"`
	StepResult   *RangeStepResult `ko:"name=step_result"`
	StartToCarry Shaper           `ko:"name=start_to_carry"`
	JoinResult   *Join            `ko:"name=join_result"`
	Cached_      *AssignCache     `ko:"name=cached"`
}

func (soln *RangeSoln) String() string { return Sprint(soln) }

func (soln *RangeSoln) Cached() *AssignCache {
	return AssignCacheUnion(
		soln.Cached_,
		soln.Step.Cached(),
		soln.JoinResult.Cached(),
	)
}

func (soln *RangeSoln) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		soln.Arg.CircuitEffect(),
		soln.Step.CircuitEffect(),
		soln.StepResult.CircuitEffect(),
		soln.StartToCarry.CircuitEffect(),
		soln.JoinResult.CircuitEffect(),
	).AggregateDuctFunc(soln.FuncExpr())
}

func (soln *RangeSoln) OverIsUnknown() bool { return soln.Arg.OverIsUnknown }

func (soln *RangeSoln) OverIsEmpty() bool { return soln.Arg.OverIsEmpty }

func (soln *RangeSoln) ProgramEffect() *GoProgramEffect {
	return AggregateProgramEffects(
		soln.Step.ProgramEffect(),
		soln.JoinResult.ProgramEffect(),
	)
}

type RangeStepResult struct {
	Carry       GoType `ko:"name=carry"`
	SelectCarry Shaper `ko:"name=select_carry"`
	Emit        GoType `ko:"name=emit"`
	SelectEmit  Shaper `ko:"name=select_emit"`
}

func (stepResult *RangeStepResult) CircuitEffect() *GoCircuitEffect {
	return AggregateCircuitEffects(
		stepResult.SelectCarry.CircuitEffect(),
		stepResult.SelectEmit.CircuitEffect(),
	)
}

func (soln *RangeSoln) Returns() GoType { return soln.JoinResult.Returns() }

func (soln *RangeSoln) Image() GoType {
	return NewGoSlice(soln.StepResult.Emit)
}

func SolveRangeArg(span *Span, arg GoStructure) (soln *RangeSolnArg, err error) {
	soln = &RangeSolnArg{Arg: arg}
	if soln.Over.Extractor, soln.Over.Type, err = GoSelectSimplify(span, Path{"over"}, arg); err != nil {
		return nil, span.Errorf(err, "range expects an “over” argument")
	}
	simplifiedOver, _ := Simplify(span, soln.Over.Type)
	switch simplifiedOver.(type) {
	case Unknown:
		soln.OverIsUnknown = true
		return soln, nil
	case *GoEmpty:
		soln.OverIsEmpty = true
		return soln, nil
	}
	if soln.OverLifted, soln.OverLift, err = LiftToSequence(span, soln.Over.Type); err != nil {
		return nil, span.Errorf(err, "range lifting “over” to series")
	}
	if soln.With.Extractor, soln.With.Varietal, err = GoSelectVarietal(span, Path{"with"}, arg); err != nil {
		return nil, span.Errorf(err, "range expects a “with” argument")
	}
	if soln.Start.Extractor, soln.Start.Type, err = GoSelect(span, Path{"start"}, arg); err != nil {
		return nil, span.Errorf(err, "range expects a “start” argument")
	}
	return soln, nil
}

func SolveRange(span *Span, arg GoStructure) (soln *RangeSoln, err error) {
	soln = &RangeSoln{Origin: span}
	if soln.Arg, err = SolveRangeArg(span, arg); err != nil {
		return nil, err
	}
	if soln.Arg.OverIsUnknown || soln.Arg.OverIsEmpty {
		return soln, nil
	}
	carry := soln.Arg.Start.Type
	span = SpanClearWeavingCtx(span) // clear weaving recursion history
	assn := NewAssignCtx(span)
	for {
		if soln.Step, err = GoEvoke(
			span,
			soln.Arg.With.Varietal,
			rangeStepWith(carry, soln.Arg.OverLifted.SequenceElem()),
		); err != nil {
			return nil, err
		}
		if soln.StepResult, err = rangeExtractStepResult(span, soln.Step.Returns()); err != nil {
			return nil, err
		}
		if unified, err := Generalize(span, carry, soln.StepResult.Carry); err != nil {
			return nil, span.Errorf(err, "solving carry generalization")
		} else {
			if _, err = assn.Assign(unified, carry); err != nil {
				carry = unified
			} else { // otherwise, carry is at fixed point
				if soln.StartToCarry, err = assn.Assign(soln.Arg.Start.Type, carry); err != nil {
					panic("o")
				}
				if soln.JoinResult, err = GoJoin(span, rangeJoinWith(soln.Image(), carry)); err != nil {
					return nil, span.Errorf(err, "joining range result")
				}
				soln.Cached_ = assn.Flush()
				return soln, nil
			}
		}
	}
	panic("o")
}

func rangeStepWith(carry, elem GoType) []*GoAugmentField {
	return WithGoField(
		&GoField{Type: carry, Name: GoNameFor("carry"), Tag: KoTags("carry", false)},
		&GoField{Type: elem, Name: GoNameFor("elem"), Tag: KoTags("elem", false)},
	)
}

func rangeJoinWith(image, carry GoType) []*GoAugmentField {
	return WithGoField(
		&GoField{Type: image, Name: GoNameFor("image"), Tag: KoTags("image", false)},
		&GoField{Type: carry, Name: GoNameFor("residue"), Tag: KoTags("residue", false)},
	)
}

func rangeExtractStepResult(span *Span, stepReturned GoType) (result *RangeStepResult, err error) {
	result = &RangeStepResult{}
	if result.SelectCarry, result.Carry, err = GoSelect(span, Path{"carry"}, stepReturned); err != nil {
		return nil, span.Errorf(err, "range step must return a carry")
	}
	if result.SelectEmit, result.Emit, err = GoSelect(span, Path{"emit"}, stepReturned); err != nil {
		return nil, span.Errorf(err, "range step must return an emit")
	}
	return
}

// func range_xyz789(step_ctx *runtime.Contex, arg RangeArg) RangeResult {
// 	carry := StartToCarry(ExtractStart(arg))
// 	var image []StepResultEmit
// 	for i, elem := range ExtractOver(LiftOver(arg)) {
// 		stepResult := EVOKE(ExtractWith(arg), carry, elem)
// 		carry = StepResultSelectCarry(stepResult)
// 		image = append(image, StepResultSelectEmit(stepResult))
// 	}
// 	return JOIN(image, carry)
// }

func (soln *RangeSoln) FuncName() *GoNameExpr {
	return &GoNameExpr{
		Origin: soln.Origin,
		Name:   fmt.Sprintf("range_%s", soln.Origin.SpanID().String()),
	}
}

// FormExpr returns an invocation of the range_ function with (step_ctx, RangeArgStructure)
func (soln *RangeSoln) FormExpr(arg ...*GoSlotExpr) GoExpr {
	return &GoCallFuncExpr{
		Func: soln.FuncExpr(),
		Arg: []GoExpr{
			&GoVerbatimExpr{"step_ctx"},
			FindSlotExpr(arg, RootSlot{}),
		},
	}
}

func (soln *RangeSoln) FuncExpr() GoFuncExpr {
	return &GoDuctFuncExpr{
		Comment:    fmt.Sprintf("(%s)", soln.Origin.SourceLine()),
		FuncName:   soln.FuncName(),
		ArgName:    soln.argExpr(),
		ArgType:    soln.Arg.Arg,
		ReturnType: soln.Returns(),
		Line: MergeLine(
			soln.initExpr(),
			soln.forExpr(),
			soln.returnExpr(),
		),
	}
}

func (soln *RangeSoln) argExpr() GoExpr { return &GoVerbatimExpr{"arg"} }

func (soln *RangeSoln) initExpr() []GoExpr {
	return []GoExpr{
		// carry := StartToCarry(SelectStart(arg))
		&GoColonAssignExpr{
			Left: soln.carryExpr(),
			Right: &GoShapeExpr{
				Shaper: soln.StartToCarry,
				Expr: &GoShapeExpr{
					Shaper: soln.Arg.Start.Extractor,
					Expr:   soln.argExpr(),
				},
			},
		},
		// var image []StepResultEmit
		&GoVarDeclExpr{Name: soln.imageExpr(), Type: soln.Image()},
	}
}

func (soln *RangeSoln) carryExpr() GoExpr { return &GoVerbatimExpr{"carry"} }

func (soln *RangeSoln) imageExpr() GoExpr { return &GoVerbatimExpr{"image"} }

func (soln *RangeSoln) forExpr() []GoExpr {
	indexExpr, elemExpr := &GoVerbatimExpr{"i"}, &GoVerbatimExpr{"elem"}
	stepResultExpr := &GoVerbatimExpr{"stepResult"}
	return []GoExpr{
		&GoForExpr{
			// for i, elem := range ExtractOver(LiftOver(arg))
			Range: &GoColonAssignExpr{
				Left:  &GoListExpr{Elem: []GoExpr{indexExpr, elemExpr}},
				Right: &GoRangeExpr{Range: soln.overExpr()},
			},
			Line: []GoExpr{
				// _ = i
				&GoAssignExpr{
					Left:  UnderlineExpr,
					Right: indexExpr,
				},
				// stepResult := Evoke(ExtractWith(arg), carry, elem)
				&GoColonAssignExpr{
					Left: stepResultExpr,
					Right: &GoSlotFormExpr{
						SlotExpr: []*GoSlotExpr{
							{Slot: RootSlot{}, Expr: soln.withExpr()},
							{Slot: NameSlot{"carry"}, Expr: soln.carryExpr()},
							{Slot: NameSlot{"elem"}, Expr: elemExpr},
						},
						Form: soln.Step,
					},
				},
				// carry = StepResultSelectCarry(stepResult)
				&GoAssignExpr{
					Left: soln.carryExpr(),
					Right: &GoShapeExpr{
						Shaper: soln.StepResult.SelectCarry,
						Expr:   stepResultExpr,
					},
				},
				// image = append(image, StepResultSelectEmit(stepResult))
				&GoAssignExpr{
					Left: soln.imageExpr(),
					Right: &GoAppendExpr{
						Base: soln.imageExpr(),
						Elem: []GoExpr{
							soln.imageExpr(),
							&GoShapeExpr{
								Shaper: soln.StepResult.SelectEmit,
								Expr:   stepResultExpr,
							},
						},
					},
				},
			}, // for line
		}, // for
	}
}

func (soln *RangeSoln) overExpr() GoExpr { // select and lift
	return &GoShapeExpr{
		Shaper: soln.Arg.OverLift,
		Expr: &GoShapeExpr{
			Shaper: soln.Arg.Over.Extractor,
			Expr:   soln.argExpr(),
		},
	}
}

func (soln *RangeSoln) withExpr() GoExpr { // select and lift
	return &GoShapeExpr{
		Shaper: soln.Arg.With.Extractor,
		Expr:   soln.argExpr(),
	}
}

func (soln *RangeSoln) returnExpr() []GoExpr {
	return []GoExpr{
		&GoReturnExpr{
			Expr: &GoSlotFormExpr{
				SlotExpr: []*GoSlotExpr{
					{Slot: NameSlot{"image"}, Expr: soln.imageExpr()},
					{Slot: NameSlot{"residue"}, Expr: soln.carryExpr()},
				},
				Form: soln.JoinResult,
			},
		},
	}
}
