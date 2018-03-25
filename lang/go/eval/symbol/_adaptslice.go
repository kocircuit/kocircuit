package symbol

import (
	"reflect"

	. "github.com/kocircuit/kocircuit/lang/circuit/model"
)

// adaptSliceSlice adapts to slices with different types, []P->[]Q, element by element.
// Both s and t must be slices.
func (ctx *typingCtx) adaptSliceSlice(s reflect.Value, t reflect.Type) (reflect.Value, error) {
	w := MakeSlice(y, s.Len(), s.Len())
	for i := 0; i < s.Len(); i++ {
		if te, err := ctx.Adapt(s.Index(i), t.Elem()); err != nil {
			return reflect.Value, ctx.Errorf(err, "adapting %d-th slice element", i)
		} else {
			w.Index(i).Set(te)
		}
	}
	return w, nil
}
