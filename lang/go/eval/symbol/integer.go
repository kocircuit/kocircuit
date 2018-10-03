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

package symbol

import (
	"github.com/kocircuit/kocircuit/lang/circuit/model"
)

func ExtractMonadicNonEmptyIntegerSeries(span *model.Span, from *StructSymbol) (_ Symbols, signed bool, _ error) {
	series := from.LinkField("value", true).LiftToSeries(span)
	if series.IsEmpty() {
		return nil, false, span.Errorf(nil, "series is empty")
	}
	switch {
	case IsSignedIntegerType(series.Type_.Elem):
		return series.Elem, true, nil
	case IsUnsignedIntegerType(series.Type_.Elem):
		return series.Elem, false, nil
	default:
		return nil, false, span.Errorf(nil, "series element %v not integral", series.Type_.Elem)
	}
}

func IsSignedIntegerType(t Type) bool {
	if basic, ok := t.(BasicType); !ok {
		return false
	} else {
		switch basic {
		case BasicInt8, BasicInt16, BasicInt32, BasicInt64:
			return true
		default:
			return false
		}
	}
}

func IsUnsignedIntegerType(t Type) bool {
	if basic, ok := t.(BasicType); !ok {
		return false
	} else {
		switch basic {
		case BasicUint8, BasicUint16, BasicUint32, BasicUint64:
			return true
		default:
			return false
		}
	}
}

func SignedMaximal(sym Symbol) int64 {
	return sym.(BasicSymbol).GoValue().Convert(goInt64).Interface().(int64)
}

func UnsignedMaximal(sym Symbol) uint64 {
	return sym.(BasicSymbol).GoValue().Convert(goUint64).Interface().(uint64)
}
