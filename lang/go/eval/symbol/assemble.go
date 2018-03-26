package symbol

import (
	. "github.com/kocircuit/kocircuit/lang/circuit/model"
	pb "github.com/kocircuit/kocircuit/lang/go/eval/symbol/proto"
)

// Assemble panics on invalid protocol structure.
func Assemble(span *Span, pbSymbol *pb.Symbol) Symbol {
	ctx := &typingCtx{Span: span}
	return ctx.Assemble(pbSymbol)
}

func (ctx *typingCtx) Assemble(pbSymbol *pb.Symbol) Symbol {
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
		default:
			panic(ctx.Errorf(nil, "unknown symbol"))
		}
	}
}

func (ctx *typingCtx) AssembleBasic(pbBasic *pb.SymbolBasic) Symbol {
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

func (ctx *typingCtx) AssembleSeries(pbSeries *pb.SymbolSeries) Symbol {
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

func (ctx *typingCtx) AssembleStruct(pbStruct *pb.SymbolStruct) Symbol {
	if pbStruct == nil {
		return EmptySymbol{}
	}
	asmFields := make(FieldSymbols, 0, len(pbStruct.Field))
	for _, field := range pbStruct.Field {
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
	return MakeStructSymbol(asmFields)
}

func (ctx *typingCtx) AssembleMap(pbMap *pb.SymbolMap) Symbol {
	if pbMap == nil {
		return EmptySymbol{}
	}
	asmKeyValues := make(KeyValueSymbols, 0, len(pbMap.KeyValue))
	for _, keyValue := range pbMap.KeyValue {
		ctx2 := ctx.Refine(keyValue.GetKey())
		if asmValue := ctx2.Assemble(keyValue.GetValue()); !IsEmptySymbol(asmValue) {
			asmKeyValues = append(asmKeyValues,
				&KeyValueSymbol{
					Key:   keyValue.GetKey(),
					Value: asmValue,
				},
			)
		}
	}
	if ms, err := MakeMapSymbol(ctx.Span, asmKeyValues); err != nil {
		panic(err)
	} else {
		return ms
	}
}

func (ctx *typingCtx) AssembleBlob(pbBlob *pb.SymbolBlob) Symbol {
	panic("XXX")
}
