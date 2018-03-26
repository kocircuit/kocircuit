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
		switch u := pbSymbol.Symbol.(type) {
		case *pb.Symbol_Basic:
			return ctx.AssembleBasic(u.Basic)
		case *pb.Symbol_Series:
			return ctx.AssembleSeries(u.Series)
		case *pb.Symbol_Struct:
			return ctx.AssembleStruct(u.Struct)
		case *pb.Symbol_Map:
			return ctx.AssembleMap(u.Map)
		case *pb.Symbol_Blob:
			return ctx.AssembleBlob(u.Blob)
		default:
			panic(ctx.Errorf(nil, "unknown symbol"))
		}
	}
}

func (ctx *typingCtx) AssembleBasic(pbBasic *pb.SymbolBasic) Symbol {
	if pbBasic == nil {
		return EmptySymbol{}
	} else {
		switch u := pbBasic.Basic.(type) {
		case *pb.SymbolBasic_Bool:
			return BasicSymbol{Value: bool(u.Bool)}
		case *pb.SymbolBasic_String_:
			return BasicSymbol{Value: string(u.String_)}
		case *pb.SymbolBasic_Int8:
			return BasicSymbol{Value: int8(u.Int8)}
		case *pb.SymbolBasic_Int16:
			return BasicSymbol{Value: int16(u.Int16)}
		case *pb.SymbolBasic_Int32:
			return BasicSymbol{Value: int32(u.Int32)}
		case *pb.SymbolBasic_Int64:
			return BasicSymbol{Value: int64(u.Int64)}
		case *pb.SymbolBasic_Uint8:
			return BasicSymbol{Value: uint8(u.Uint8)}
		case *pb.SymbolBasic_Uint16:
			return BasicSymbol{Value: uint16(u.Uint16)}
		case *pb.SymbolBasic_Uint32:
			return BasicSymbol{Value: uint32(u.Uint32)}
		case *pb.SymbolBasic_Uint64:
			return BasicSymbol{Value: uint64(u.Uint64)}
		case *pb.SymbolBasic_Float32:
			return BasicSymbol{Value: float32(u.Float32)}
		case *pb.SymbolBasic_Float64:
			return BasicSymbol{Value: float64(u.Float64)}
		default:
			panic(ctx.Errorf(nil, "unknown basic symbol"))
		}
	}
}

func (ctx *typingCtx) AssembleSeries(pbSeries *pb.SymbolSeries) Symbol {
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
	panic("XXX")
}

func (ctx *typingCtx) AssembleMap(pbMap *pb.SymbolMap) Symbol {
	panic("XXX")
}

func (ctx *typingCtx) AssembleBlob(pbBlob *pb.SymbolBlob) Symbol {
	panic("XXX")
}
