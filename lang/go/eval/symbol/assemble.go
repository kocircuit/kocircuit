package symbol

import (
	"fmt"
	"strconv"

	. "github.com/kocircuit/kocircuit/lang/circuit/eval"
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
)

type VarietyAssembler interface {
	AssembleMacro(span *Span, pkgPath, funcName string) (Macro, error)
}

func AssembleWithError(span *Span, asm VarietyAssembler, pbSymbol *pb.Symbol) (res Symbol, err error) {
	defer func() {
		if r := recover(); r != nil {
			res, err = nil, r.(error)
		}
	}()
	return Assemble(span, asm, pbSymbol), nil
}

// Assemble panics on invalid protocol structure.
func Assemble(span *Span, asm VarietyAssembler, pbSymbol *pb.Symbol) Symbol {
	ctx := &assemblingCtx{Span: span, Asm: asm}
	return ctx.Assemble(pbSymbol)
}

func (ctx *assemblingCtx) Assemble(pbSymbol *pb.Symbol) Symbol {
	if pbSymbol == nil {
		return EmptySymbol{}
	} else {
		switch pbSymbol.Symbol.(type) {
		case *pb.Symbol_Basic:
			return ctx.AssembleBasic(pbSymbol.GetBasic())
		case *pb.Symbol_Series:
			return ctx.AssembleSeries(pbSymbol.GetSeries())
		case *pb.Symbol_Struct:
			return ctx.AssembleStruct(pbSymbol.GetStruct())
		case *pb.Symbol_Map:
			return ctx.AssembleMap(pbSymbol.GetMap())
		case *pb.Symbol_Blob:
			return ctx.AssembleBlob(pbSymbol.GetBlob())
		case *pb.Symbol_Variety:
			return ctx.AssembleVariety(pbSymbol.GetVariety())
		default:
			panic(ctx.Errorf(nil, "unknown symbol"))
		}
	}
}

func (ctx *assemblingCtx) AssembleBasic(pbBasic *pb.SymbolBasic) Symbol {
	if pbBasic == nil {
		return EmptySymbol{}
	} else {
		switch pbBasic.Basic.(type) {
		case *pb.SymbolBasic_Bool:
			return BasicSymbol{Value: bool(pbBasic.GetBool())}
		case *pb.SymbolBasic_String_:
			return BasicSymbol{Value: string(pbBasic.GetString_())}
		case *pb.SymbolBasic_Int8:
			return BasicSymbol{Value: int8(pbBasic.GetInt8())}
		case *pb.SymbolBasic_Int16:
			return BasicSymbol{Value: int16(pbBasic.GetInt16())}
		case *pb.SymbolBasic_Int32:
			return BasicSymbol{Value: int32(pbBasic.GetInt32())}
		case *pb.SymbolBasic_Int64:
			return BasicSymbol{Value: int64(pbBasic.GetInt64())}
		case *pb.SymbolBasic_Uint8:
			return BasicSymbol{Value: uint8(pbBasic.GetUint8())}
		case *pb.SymbolBasic_Uint16:
			return BasicSymbol{Value: uint16(pbBasic.GetUint16())}
		case *pb.SymbolBasic_Uint32:
			return BasicSymbol{Value: uint32(pbBasic.GetUint32())}
		case *pb.SymbolBasic_Uint64:
			return BasicSymbol{Value: uint64(pbBasic.GetUint64())}
		case *pb.SymbolBasic_Float32:
			return BasicSymbol{Value: float32(pbBasic.GetFloat32())}
		case *pb.SymbolBasic_Float64:
			return BasicSymbol{Value: float64(pbBasic.GetFloat64())}
		default:
			panic(ctx.Errorf(nil, "unknown basic symbol"))
		}
	}
}

func (ctx *assemblingCtx) AssembleSeries(pbSeries *pb.SymbolSeries) Symbol {
	if pbSeries == nil {
		return EmptySymbol{}
	}
	asmElems := make(Symbols, 0, len(pbSeries.Element))
	ctx2 := ctx.Refine("()")
	for _, elem := range pbSeries.Element {
		if asmElem := ctx2.Assemble(elem); !IsEmptySymbol(asmElem) {
			asmElems = append(asmElems, asmElem)
		}
	}
	if series, err := MakeSeriesSymbol(ctx.Span, asmElems); err != nil {
		panic(err)
	} else {
		return series
	}
}

func (ctx *assemblingCtx) AssembleStruct(pbStruct *pb.SymbolStruct) Symbol {
	if pbStruct == nil {
		return EmptySymbol{}
	}
	return MakeStructSymbol(ctx.AssembleFields(pbStruct.GetField()))
}

func (ctx *assemblingCtx) AssembleFields(pbFields []*pb.SymbolField) FieldSymbols {
	asmFields := make(FieldSymbols, 0, len(pbFields))
	for _, field := range pbFields {
		ctx2 := ctx.Refine(field.GetName())
		if asmFieldValue := ctx2.Assemble(field.GetValue()); !IsEmptySymbol(asmFieldValue) {
			asmFields = append(asmFields,
				&FieldSymbol{
					Name:    field.GetName(),
					Monadic: field.GetMonadic(),
					Value:   asmFieldValue,
				},
			)
		}
	}
	return asmFields
}

func (ctx *assemblingCtx) AssembleMap(pbMap *pb.SymbolMap) Symbol {
	if pbMap == nil {
		return EmptySymbol{}
	}
	asmMap := map[string]Symbol{}
	for _, keyValue := range pbMap.KeyValue {
		ctx2 := ctx.Refine(keyValue.GetKey())
		if asmValue := ctx2.Assemble(keyValue.GetValue()); !IsEmptySymbol(asmValue) {
			asmMap[keyValue.GetKey()] = asmValue
		}
	}
	if ms, err := MakeMapSymbol(ctx.Span, asmMap); err != nil {
		panic(err)
	} else {
		return ms
	}
}

func (ctx *assemblingCtx) AssembleBlob(pbBlob *pb.SymbolBlob) Symbol {
	return MakeBlobSymbol(pbBlob.GetBytes())
}

func (ctx *assemblingCtx) AssembleVariety(pbVariety *pb.SymbolVariety) Symbol {
	if m, err := ctx.Asm.AssembleMacro(
		ctx.Span,
		pbVariety.GetPkgPath(),
		pbVariety.GetFuncName(),
	); err != nil {
		panic(err)
	} else {
		return MakeVarietySymbol(m, ctx.AssembleFields(pbVariety.GetArg()))
	}
}

// context

type assemblingCtx struct {
	Parent *assemblingCtx   `ko:"name=parent"`
	Span   *Span            `ko:"name=span"`
	Walk   string           `ko:"name=walk"`
	Asm    VarietyAssembler `ko:"name=asm"`
}

func (ctx *assemblingCtx) Refine(walk string) *assemblingCtx {
	return &assemblingCtx{
		Parent: ctx,
		Span:   ctx.Span,
		Walk:   walk,
		Asm:    ctx.Asm,
	}
}

func (ctx *assemblingCtx) RefineIndex(i int) *assemblingCtx {
	return ctx.Refine(strconv.Itoa(i))
}

func (ctx *assemblingCtx) Path() Path {
	if ctx == nil {
		return nil
	} else if ctx.Parent == nil {
		return Path{ctx.Walk}
	} else {
		return append(ctx.Parent.Path(), ctx.Walk)
	}
}

func (ctx *assemblingCtx) Errorf(cause error, format string, arg ...interface{}) error {
	return ctx.Span.ErrorfSkip(
		2, cause,
		fmt.Sprintf("%v: %s", ctx.Path(), fmt.Sprintf(format, arg...)),
	)
}
