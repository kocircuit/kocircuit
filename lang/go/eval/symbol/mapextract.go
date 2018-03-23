package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

func ExtractMapOpaque(span *Span, s Symbol, mapType reflect.Type) (*OpaqueSymbol, error) {
	if mapValue, err := extractMapValue(span, s, mapType); err != nil {
		return nil, err
	} else {
		return &OpaqueSymbol{Value: mapValue}, nil
	}
}

func extractMapValue(span *Span, s Symbol, mapType reflect.Type) (reflect.Value, error) {
	series := s.LiftToSeries(span)
	kt, vt := mapType.Key(), mapType.Elem()
	mapValue := reflect.MakeMapWithSize(mapType, len(series.Elem))
	for _, e := range series.Elem {
		row, ok := e.(*StructSymbol) // (key: K, value: V)
		if !ok {
			return reflect.Value{}, span.Errorf(nil, "expecting key/value struct")
		}
		kv, err := Integrate(span, row.Walk("key"), kt)
		if err != nil {
			return reflect.Value{}, span.Errorf(err, "expecting key of type %v", kt)
		}
		vv, err := Integrate(span, row.Walk("value"), vt)
		if err != nil {
			return reflect.Value{}, span.Errorf(err, "expecting value of type %v", vt)
		}
		mapValue.SetMapIndex(kv, vv)
	}
	return mapValue, nil
}
